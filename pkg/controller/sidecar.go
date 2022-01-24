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
	"context"
	"fmt"
	"time"

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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	core_util "kmodules.xyz/client-go/core/v1"
	wapi "kmodules.xyz/webhook-runtime/apis/workload/v1"
)

func (c *StashController) ensureBackupSidecar(w *wapi.Workload, inv invoker.BackupInvoker, targetInfo invoker.BackupTargetInfo, caller string) error {
	invMeta := inv.GetObjectMeta()

	sa := stringz.Val(w.Spec.Template.Spec.ServiceAccountName, "default")
	owner, err := ownerWorkload(w)
	if err != nil {
		return err
	}

	// Don't create RBAC stuff when the caller is webhook to make the webhooks side effect free.
	if caller != apis.CallerWebhook {
		err = stash_rbac.EnsureSidecarRoleBinding(c.kubeClient, owner, invMeta.Namespace, sa, inv.GetLabels())
		if err != nil {
			return err
		}
	}

	// if the Stash is using a private registry, then ensure the image pull secrets
	if c.ImagePullSecrets != nil {
		var imagePullSecrets []core.LocalObjectReference
		imagePullSecrets, err = c.ensureImagePullSecrets(invMeta, inv.GetOwnerRef())
		if err != nil {
			return err
		}
		w.Spec.Template.Spec.ImagePullSecrets = imagePullSecrets
	}

	repository, err := inv.GetRepository()
	if err != nil {
		klog.Errorf("unable to get repository %s/%s: Reason: %v", inv.GetRepoRef().Namespace, inv.GetRepoRef().Name, err)
		return err
	}

	if w.Spec.Template.Annotations == nil {
		w.Spec.Template.Annotations = map[string]string{}
	}
	// mark pods with BackupConfiguration spec hash. used to force restart pods for rc/rs
	w.Spec.Template.Annotations[api_v1beta1.AppliedBackupInvokerSpecHash] = inv.GetHash()

	if targetInfo.Target == nil {
		return fmt.Errorf("target is nil")
	}

	image := docker.Docker{
		Registry: c.DockerRegistry,
		Image:    c.StashImage,
		Tag:      c.StashImageTag,
	}
	w.Spec.Template.Spec.Containers = core_util.UpsertContainer(
		w.Spec.Template.Spec.Containers,
		util.NewBackupSidecarContainer(inv, targetInfo, &repository.Spec.Backend, image),
	)

	w.Spec.Template.Spec.Volumes = util.UpsertTmpVolume(w.Spec.Template.Spec.Volumes, targetInfo.TempDir)

	if w.Annotations == nil {
		w.Annotations = make(map[string]string)
	}
	jsonObj, err := inv.GetObjectJSON()
	if err != nil {
		return err
	}
	w.Annotations[api_v1beta1.KeyLastAppliedBackupInvoker] = jsonObj
	w.Annotations[api_v1beta1.KeyLastAppliedBackupInvokerKind] = inv.GetTypeMeta().Kind

	return nil
}

func (c *StashController) ensureBackupSidecarDeleted(w *wapi.Workload) {
	// remove resource hash annotation
	if w.Spec.Template.Annotations != nil {
		delete(w.Spec.Template.Annotations, api_v1beta1.AppliedBackupInvokerSpecHash)
	}
	// remove sidecar container
	w.Spec.Template.Spec.Containers = core_util.EnsureContainerDeleted(w.Spec.Template.Spec.Containers, apis.StashContainer)

	// backup sidecar has been removed but workload still may have restore init-container
	// so removed respective volumes that were added to the workload only if the workload does not have restore init-container
	if !util.HasStashContainer(w) {
		// remove the helpers volumes that were added for sidecar
		w.Spec.Template.Spec.Volumes = util.EnsureVolumeDeleted(w.Spec.Template.Spec.Volumes, apis.ScratchDirVolumeName)
	}

	// remove respective annotations
	if w.Annotations != nil {
		delete(w.Annotations, api_v1beta1.KeyLastAppliedBackupInvoker)
		delete(w.Annotations, api_v1beta1.KeyLastAppliedBackupInvokerKind)
	}
}

// ensureWorkloadLatestState check if the workload's pod has latest update of workload specification
// if a pod does not have latest update, it deletes that pod so that new pod start with updated spec
func (c *StashController) ensureWorkloadLatestState(w *wapi.Workload) (bool, error) {
	stateChanged := false

	err := wait.PollImmediate(3*time.Second, 5*time.Minute, func() (done bool, err error) {
		r, err := metav1.LabelSelectorAsSelector(w.Spec.Selector)
		if err != nil {
			return false, err
		}
		// list all pods of this workload
		pods, err := c.kubeClient.CoreV1().Pods(w.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: r.String()})
		if err != nil {
			if errors.IsUnauthorized(err) || errors.IsForbidden(err) {
				return false, err
			}
			return false, nil // ignore temporary server errors
		}

		workloadSidecarState := util.HasStashSidecar(w.Spec.Template.Spec.Containers)
		workloadInitContainerState := util.HasStashInitContainer(w.Spec.Template.Spec.InitContainers)
		workloadBackupInvokerResourceHash := util.GetString(w.Spec.Template.Annotations, api_v1beta1.AppliedBackupInvokerSpecHash)
		workloadRestoreInvokerResourceHash := util.GetString(w.Spec.Template.Annotations, api_v1beta1.AppliedRestoreInvokerSpecHash)

		// identify the pods that does not have latest update.
		// we have to restart these pods so that it starts with latest update
		var podsToRestart []core.Pod
		for _, pod := range pods.Items {
			if !isPodOwnedByWorkload(w, pod) {
				continue
			}
			podSidecarState := util.HasStashSidecar(pod.Spec.Containers)
			podInitContainerState := util.HasStashInitContainer(pod.Spec.InitContainers)
			podBackupInvokerResourceHash := util.GetString(pod.Annotations, api_v1beta1.AppliedBackupInvokerSpecHash)
			podRestoreInvokerResourceHash := util.GetString(pod.Annotations, api_v1beta1.AppliedRestoreInvokerSpecHash)

			if workloadSidecarState != podSidecarState ||
				workloadInitContainerState != podInitContainerState ||
				workloadBackupInvokerResourceHash != podBackupInvokerResourceHash ||
				workloadRestoreInvokerResourceHash != podRestoreInvokerResourceHash {

				podsToRestart = append(podsToRestart, pod)
			}
		}

		if len(podsToRestart) == 0 {
			return true, nil // done
		}
		stateChanged = true
		for _, pod := range podsToRestart {
			err := c.kubeClient.CoreV1().Pods(w.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				klog.Errorln(err)
			}
		}
		return false, nil // try again
	})
	if err != nil {
		return false, err
	}

	return stateChanged, nil
}

func isPodOwnedByWorkload(w *wapi.Workload, pod core.Pod) bool {
	for _, ref := range pod.OwnerReferences {
		if ref.Kind == w.Kind && ref.Name == w.Name {
			return true
		}
	}
	return false
}

func (c *StashController) handleSidecarInjectionFailure(w *wapi.Workload, inv invoker.BackupInvoker, tref api_v1beta1.TargetRef, err error) error {
	klog.Warningf("Failed to inject stash sidecar into %s %s/%s. Reason: %v", w.Kind, w.Namespace, w.Name, err)

	// Failed to inject stash sidecar. So, set "StashSidecarInjected" condition to "False".
	cerr := conditions.SetSidecarInjectedConditionToFalse(inv, tref, err)

	// write event to respective resource
	_, err2 := eventer.CreateEvent(
		c.kubeClient,
		eventer.EventSourceWorkloadController,
		w.Object,
		core.EventTypeWarning,
		eventer.EventReasonSidecarInjectionFailed,
		fmt.Sprintf("Failed to inject stash sidecar into %s %s/%s. Reason: %v", w.Kind, w.Namespace, w.Name, err),
	)
	return utilerrors.NewAggregate([]error{err2, cerr})
}

func (c *StashController) handleSidecarInjectionSuccess(w *wapi.Workload, inv invoker.BackupInvoker, tref api_v1beta1.TargetRef) error {
	klog.Infof("Successfully injected stash sidecar into %s %s/%s.", w.Kind, w.Namespace, w.Name)

	// Set "StashSidecarInjected" condition to "True"
	cerr := conditions.SetSidecarInjectedConditionToTrue(inv, tref)

	// write event to respective resource
	_, err2 := eventer.CreateEvent(
		c.kubeClient,
		eventer.EventSourceWorkloadController,
		w.Object,
		core.EventTypeNormal,
		eventer.EventReasonSidecarInjectionSucceeded,
		fmt.Sprintf("Successfully injected stash sidecar into %s %s/%s.", w.Kind, w.Namespace, w.Name),
	)
	return utilerrors.NewAggregate([]error{err2, cerr})
}

func (c *StashController) handleSidecarDeletionSuccess(w *wapi.Workload) error {
	klog.Infof("Successfully removed stash sidecar from %s %s/%s.", w.Kind, w.Namespace, w.Name)

	// write event to respective resource
	_, err2 := eventer.CreateEvent(
		c.kubeClient,
		eventer.EventSourceWorkloadController,
		w.Object,
		core.EventTypeNormal,
		eventer.EventReasonSidecarDeletionSucceeded,
		fmt.Sprintf("Successfully removed stash sidecar from %s %s/%s.", w.Kind, w.Namespace, w.Name),
	)
	return err2
}
