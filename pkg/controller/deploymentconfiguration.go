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

	"stash.appscode.dev/apimachinery/apis"
	stash_rbac "stash.appscode.dev/stash/pkg/rbac"
	"stash.appscode.dev/stash/pkg/util"

	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/discovery"
	"kmodules.xyz/client-go/tools/queue"
	ocapps "kmodules.xyz/openshift/apis/apps/v1"
	ocapps_util "kmodules.xyz/openshift/client/clientset/versioned/typed/apps/v1/util"
	"kmodules.xyz/webhook-runtime/admission"
	hooks "kmodules.xyz/webhook-runtime/admission/v1beta1"
	webhook "kmodules.xyz/webhook-runtime/admission/v1beta1/workload"
	wapi "kmodules.xyz/webhook-runtime/apis/workload/v1"
	wcs "kmodules.xyz/webhook-runtime/client/workload/v1"
)

func (c *StashController) NewDeploymentConfigWebhook() hooks.AdmissionHook {
	return webhook.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.stash.appscode.com",
			Version:  "v1alpha1",
			Resource: "deploymentconfigmutators",
		},
		"deploymentconfigmutator",
		"DeploymentConfigMutator",
		nil,
		&admission.ResourceHandlerFuncs{
			CreateFunc: func(obj runtime.Object) (runtime.Object, error) {
				w := obj.(*wapi.Workload)
				// apply stash backup/restore logic on this workload
				_, err := c.applyStashLogic(w, apis.CallerWebhook)
				return w, err
			},
			UpdateFunc: func(oldObj, newObj runtime.Object) (runtime.Object, error) {
				w := newObj.(*wapi.Workload)
				// apply stash backup/restore logic on this workload
				_, err := c.applyStashLogic(w, apis.CallerWebhook)
				return w, err
			},
		},
	)
}

func (c *StashController) initDeploymentConfigWatcher() {
	if !discovery.IsPreferredAPIResource(c.kubeClient.Discovery(), ocapps.GroupVersion.String(), apis.KindDeploymentConfig) {
		klog.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", ocapps.GroupVersion.String(), apis.KindDeploymentConfig)
		return
	}
	c.dcInformer = c.ocInformerFactory.Apps().V1().DeploymentConfigs().Informer()
	c.dcQueue = queue.New(apis.KindDeploymentConfig, c.MaxNumRequeues, c.NumThreads, c.runDeploymentConfigProcessor)
	c.dcInformer.AddEventHandler(queue.DefaultEventHandler(c.dcQueue.GetQueue(), core.NamespaceAll))
	c.dcLister = c.ocInformerFactory.Apps().V1().DeploymentConfigs().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *StashController) runDeploymentConfigProcessor(key string) error {
	obj, exists, err := c.dcInformer.GetIndexer().GetByKey(key)
	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a DeploymentConfig, so that we will see a delete for one deployment
		klog.Warningf("DeploymentConfig %s does not exist anymore\n", key)

		ns, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		// workload does not exist anymore. so delete respective ConfigMapLocks if exist
		err = util.DeleteAllConfigMapLocks(c.kubeClient, ns, name, apis.KindDeploymentConfig)
		if err != nil && !kerr.IsNotFound(err) {
			return err
		}
	} else {
		klog.Infof("Sync/Add/Update for DeploymentConfig %s", key)

		dc := obj.(*ocapps.DeploymentConfig).DeepCopy()
		dc.GetObjectKind().SetGroupVersionKind(ocapps.GroupVersion.WithKind(apis.KindDeploymentConfig))

		// convert DeploymentConfig into a common object (Workload type) so that
		// we don't need to re-write stash logic for DeploymentConfig separately
		w, err := wcs.ConvertToWorkload(dc.DeepCopy())
		if err != nil {
			klog.Errorf("failed to convert DeploymentConfig %s/%s to workload type. Reason: %v", dc.Namespace, dc.Name, err)
			return err
		}

		// apply stash backup/restore logic on this workload
		modified, err := c.applyStashLogic(w, apis.CallerController)
		if err != nil {
			klog.Errorf("failed to apply stash logic on DeploymentConfig %s/%s. Reason: %v", dc.Namespace, dc.Name, err)
			return err
		}

		if modified {
			// workload has been modified. Patch the workload so that respective pods start with the updated spec
			_, _, err := ocapps_util.PatchDeploymentConfigObject(context.TODO(), c.ocClient, dc, w.Object.(*ocapps.DeploymentConfig), metav1.PatchOptions{})
			if err != nil {
				klog.Errorf("failed to update DeploymentConfig %s/%s. Reason: %v", dc.Namespace, dc.Name, err)
				return err
			}

			// TODO: Should we force restart all pods while restore?
			// otherwise one pod will restore while others are writing/reading?

			// wait until newly patched deploymentconfigs pods are ready
			err = util.WaitUntilDeploymentConfigReady(c.ocClient, dc.ObjectMeta)
			if err != nil {
				return err
			}
		}

		// if the workload does not have any stash sidecar/init-container then
		// delete respective ConfigMapLock and RBAC stuffs if exist
		err = c.ensureUnnecessaryConfigMapLockDeleted(w)
		if err != nil {
			return err
		}
		err = stash_rbac.EnsureUnnecessaryWorkloadRBACDeleted(c.kubeClient, w)
		if err != nil {
			return err
		}
	}
	return nil
}
