/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"

	"stash.appscode.dev/apimachinery/apis"
	api_v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	"stash.appscode.dev/apimachinery/pkg/conditions"
	"stash.appscode.dev/apimachinery/pkg/docker"
	"stash.appscode.dev/apimachinery/pkg/invoker"
	"stash.appscode.dev/stash/pkg/eventer"
	stash_rbac "stash.appscode.dev/stash/pkg/rbac"
	"stash.appscode.dev/stash/pkg/util"

	stringz "gomodules.xyz/x/strings"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	core_util "kmodules.xyz/client-go/core/v1"
	wapi "kmodules.xyz/webhook-runtime/apis/workload/v1"
)

func (c *StashController) ensureRestoreInitContainer(w *wapi.Workload, inv invoker.RestoreInvoker, targetInfo invoker.RestoreTargetInfo, caller string) error {
	invMeta := inv.GetObjectMeta()
	// if RBAC is enabled then ensure ServiceAccount and respective ClusterRole and RoleBinding
	sa := stringz.Val(w.Spec.Template.Spec.ServiceAccountName, "default")

	// Don't create RBAC stuff when the caller is webhook to make the webhooks side effect free.
	if caller != apis.CallerWebhook {
		owner, err := ownerWorkload(w)
		if err != nil {
			return err
		}
		err = stash_rbac.EnsureRestoreInitContainerRBAC(c.kubeClient, owner, invMeta.Namespace, sa, inv.GetLabels())
		if err != nil {
			return err
		}
	}

	// if the Stash is using a private registry, then ensure the image pull secrets
	if c.ImagePullSecrets != nil {
		var imagePullSecrets []core.LocalObjectReference
		imagePullSecrets, err := c.ensureImagePullSecrets(invMeta, inv.GetOwnerRef())
		if err != nil {
			return err
		}
		w.Spec.Template.Spec.ImagePullSecrets = imagePullSecrets
	}

	repository, err := inv.GetRepository()
	if err != nil {
		klog.Errorf("unable to get Repository %s/%s: Reason: %v", inv.GetRepoRef().Namespace, inv.GetRepoRef().Name, err)
		return err
	}

	if w.Spec.Template.Annotations == nil {
		w.Spec.Template.Annotations = map[string]string{}
	}

	// mark pods with restore Invoker spec hash. used to force restart pods for rc/rs
	w.Spec.Template.Annotations[api_v1beta1.AppliedRestoreInvokerSpecHash] = inv.GetHash()

	image := docker.Docker{
		Registry: c.DockerRegistry,
		Image:    c.StashImage,
		Tag:      c.StashImageTag,
	}

	// insert restore init container
	initContainers := []core.Container{util.NewRestoreInitContainer(inv, targetInfo, repository, image)}
	for i := range w.Spec.Template.Spec.InitContainers {
		initContainers = core_util.UpsertContainer(initContainers, w.Spec.Template.Spec.InitContainers[i])
	}
	w.Spec.Template.Spec.InitContainers = initContainers

	// add an emptyDir volume for holding temporary files
	w.Spec.Template.Spec.Volumes = util.UpsertTmpVolume(w.Spec.Template.Spec.Volumes, targetInfo.TempDir)

	// add RestoreSession definition as annotation of the workload
	if w.Annotations == nil {
		w.Annotations = make(map[string]string)
	}
	jsonObj, err := inv.GetObjectJSON()
	if err != nil {
		return err
	}
	w.Annotations[api_v1beta1.KeyLastAppliedRestoreInvoker] = jsonObj
	w.Annotations[api_v1beta1.KeyLastAppliedRestoreInvokerKind] = inv.GetTypeMeta().Kind

	return nil
}

func (c *StashController) ensureRestoreInitContainerDeleted(w *wapi.Workload) {
	// remove resource hash annotation
	if w.Spec.Template.Annotations != nil {
		delete(w.Spec.Template.Annotations, api_v1beta1.AppliedRestoreInvokerSpecHash)
	}
	// remove init-container
	w.Spec.Template.Spec.InitContainers = core_util.EnsureContainerDeleted(w.Spec.Template.Spec.InitContainers, apis.StashInitContainer)

	// restore init-container has been removed but workload still may have backup sidecar
	// so removed respective volumes that were added to the workload only if the workload does not have backup sidecar
	if !util.HasStashContainer(w) {
		// remove the helpers volumes added for init-container
		w.Spec.Template.Spec.Volumes = util.EnsureVolumeDeleted(w.Spec.Template.Spec.Volumes, apis.ScratchDirVolumeName)
	}

	// remove respective annotations
	if w.Annotations != nil {
		delete(w.Annotations, api_v1beta1.KeyLastAppliedRestoreInvoker)
		delete(w.Annotations, api_v1beta1.KeyLastAppliedRestoreInvokerKind)
	}
}

func (c *StashController) handleInitContainerInjectionFailure(w *wapi.Workload, inv invoker.RestoreInvoker, ref api_v1beta1.TargetRef, err error) error {
	klog.Warningf("Failed to inject stash init-container into %s %s/%s. Reason: %v", w.Kind, w.Namespace, w.Name, err)

	// Set "StashInitContainerInjected" condition to "False"
	cerr := conditions.SetInitContainerInjectedConditionToFalse(inv, ref, err)

	// write event to respective resource
	_, err2 := eventer.CreateEvent(
		c.kubeClient,
		eventer.EventSourceWorkloadController,
		w.Object,
		core.EventTypeWarning,
		eventer.EventReasonInitContainerInjectionFailed,
		fmt.Sprintf("Failed to inject stash init-container into %s %s/%s. Reason: %v", w.Kind, w.Namespace, w.Name, err),
	)
	return errors.NewAggregate([]error{err2, cerr})
}

func (c *StashController) handleInitContainerInjectionSuccess(w *wapi.Workload, inv invoker.RestoreInvoker, ref api_v1beta1.TargetRef) error {
	klog.Infof("Successfully injected stash init-container into %s %s/%s.", w.Kind, w.Namespace, w.Name)

	// Set "StashInitContainerInjected" condition to "True"
	cerr := conditions.SetInitContainerInjectedConditionToTrue(inv, ref)

	// write event to respective resource
	_, err2 := eventer.CreateEvent(
		c.kubeClient,
		eventer.EventSourceWorkloadController,
		w.Object,
		core.EventTypeNormal,
		eventer.EventReasonInitContainerInjectionSucceeded,
		fmt.Sprintf("Successfully injected stash init-container into %s %s/%s.", w.Kind, w.Namespace, w.Name),
	)
	return errors.NewAggregate([]error{err2, cerr})
}

func (c *StashController) handleInitContainerDeletionSuccess(w *wapi.Workload) error {
	klog.Infof("Successfully removed stash init-container from %s %s/%s.", w.Kind, w.Namespace, w.Name)

	// write event to respective resource
	_, err2 := eventer.CreateEvent(
		c.kubeClient,
		eventer.EventSourceWorkloadController,
		w.Object,
		core.EventTypeNormal,
		eventer.EventReasonInitContainerDeletionSucceeded,
		fmt.Sprintf("Successfully stash init-container from %s %s/%s.", w.Kind, w.Namespace, w.Name),
	)
	return err2
}
