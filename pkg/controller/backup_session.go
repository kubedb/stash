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
	"sort"
	"strconv"
	"strings"
	"time"

	"stash.appscode.dev/apimachinery/apis"
	"stash.appscode.dev/apimachinery/apis/stash"
	api_v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	"stash.appscode.dev/apimachinery/pkg/conditions"
	"stash.appscode.dev/apimachinery/pkg/docker"
	stashHooks "stash.appscode.dev/apimachinery/pkg/hooks"
	"stash.appscode.dev/apimachinery/pkg/invoker"
	"stash.appscode.dev/apimachinery/pkg/metrics"
	api_util "stash.appscode.dev/apimachinery/pkg/util"
	stash_rbac "stash.appscode.dev/stash/pkg/rbac"
	"stash.appscode.dev/stash/pkg/resolve"
	"stash.appscode.dev/stash/pkg/util"

	"gomodules.xyz/pointer"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	batch_util "kmodules.xyz/client-go/batch/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/meta"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/tools/queue"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	"kmodules.xyz/webhook-runtime/admission"
	hooks "kmodules.xyz/webhook-runtime/admission/v1beta1"
	webhook "kmodules.xyz/webhook-runtime/admission/v1beta1/generic"
)

const (
	BackupExecutorSidecar   = "sidecar"
	BackupExecutorCSIDriver = "csi-driver"
	BackupExecutorJob       = "job"
)

func (c *StashController) NewBackupSessionWebhook() hooks.AdmissionHook {
	return webhook.NewGenericWebhook(
		schema.GroupVersionResource{
			Group:    "admission.stash.appscode.com",
			Version:  "v1beta1",
			Resource: "backupsessionvalidators",
		},
		"backupsessionvalidator",
		[]string{stash.GroupName},
		api.SchemeGroupVersion.WithKind(api_v1beta1.ResourceKindBackupSession),
		nil,
		&admission.ResourceHandlerFuncs{
			CreateFunc: func(obj runtime.Object) (runtime.Object, error) {
				return nil, obj.(*api_v1beta1.BackupSession).IsValid()
			},
			UpdateFunc: func(oldObj, newObj runtime.Object) (runtime.Object, error) {
				// should not allow spec update
				if !meta.Equal(oldObj.(*api_v1beta1.BackupSession).Spec, newObj.(*api_v1beta1.BackupSession).Spec) {
					return nil, fmt.Errorf("BackupSession spec is immutable")
				}
				return nil, nil
			},
		},
	)
}

// process only add events
func (c *StashController) initBackupSessionWatcher() {
	c.backupSessionInformer = c.stashInformerFactory.Stash().V1beta1().BackupSessions().Informer()
	c.backupSessionQueue = queue.New(api_v1beta1.ResourceKindBackupSession, c.MaxNumRequeues, c.NumThreads, c.runBackupSessionProcessor)
	if c.auditor != nil {
		c.backupSessionInformer.AddEventHandler(c.auditor.ForGVK(api_v1beta1.SchemeGroupVersion.WithKind(api_v1beta1.ResourceKindBackupSession)))
	}
	c.backupSessionInformer.AddEventHandler(queue.DefaultEventHandler(c.backupSessionQueue.GetQueue(), core.NamespaceAll))
	c.backupSessionLister = c.stashInformerFactory.Stash().V1beta1().BackupSessions().Lister()
}

func (c *StashController) runBackupSessionProcessor(key string) error {
	obj, exists, err := c.backupSessionInformer.GetIndexer().GetByKey(key)
	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if !exists {
		klog.Warningf("BackupSession %s does not exist anymore\n", key)
		return nil
	}

	backupSession := obj.(*api_v1beta1.BackupSession)
	klog.Infof("Sync/Add/Update for BackupSession %s", backupSession.GetName())
	session := invoker.NewBackupSessionHandler(c.stashClient, backupSession)
	return c.applyBackupSessionReconciliationLogic(session)
}

func (c *StashController) applyBackupSessionReconciliationLogic(session *invoker.BackupSessionHandler) error {
	inv, err := session.GetInvoker()
	if err != nil {
		return err
	}

	if isSessionSkipped(session) {
		klog.Infof("Skipping processing BackupSession %s/%s. Reason: %q condition is 'True'.",
			session.GetObjectMeta().Namespace,
			session.GetObjectMeta().Name,
			api_v1beta1.BackupSkipped,
		)
		return nil
	}

	if isBackupRunning(session) && isBackupDeadlineSet(session) {
		if isBackupDeadlineExceeded(session) {
			klog.Infof("Time Limit exceeded for BackupSession  %s/%s.",
				session.GetObjectMeta().Namespace,
				session.GetObjectMeta().Name,
			)
			return conditions.SetBackupDeadlineExceededConditionToTrue(session, inv.GetTimeOut())
		}
	}

	if invoker.BackupCompletedForAllTargets(session.GetTargetStatus(), len(inv.GetTargetInfo())) {
		if !postBackupHooksExecuted(inv.GetTargetInfo(), session.GetTargetStatus()) {
			klog.Infof("Waiting for postBackup hook to be executed for %s %s/%s.",
				inv.GetTypeMeta().Kind,
				inv.GetObjectMeta().Namespace,
				inv.GetObjectMeta().Name,
			)
			return nil
		}

		// Cleanup old BackupSession according to backupHistoryLimit
		if !backupHistoryCleaned(session.GetConditions()) {
			err = c.cleanupBackupHistory(session, inv.GetBackupHistoryLimit())
			if err != nil {
				condErr := conditions.SetBackupHistoryCleanedConditionToFalse(session, err)
				if condErr != nil {
					return condErr
				}
			}
		}

		if !globalPostBackupHookExecuted(inv, session) {
			err = c.executeGlobalPostBackupHook(inv, session)
			if err != nil {
				condErr := conditions.SetGlobalPostBackupHookSucceededConditionToFalse(session, err)
				if condErr != nil {
					return condErr
				}
			}
		}

		if !backupMetricPushed(session.GetConditions()) {
			err = c.sendBackupMetrics(inv, session)
			if err != nil {
				condErr := conditions.SetBackupMetricsPushedConditionToFalse(session, err)
				if condErr != nil {
					return condErr
				}
			}
		}

		klog.Infof("Skipping processing BackupSession %s/%s. Reason: phase is %q.",
			session.GetObjectMeta().Namespace,
			session.GetObjectMeta().Name,
			session.GetStatus().Phase,
		)
		return nil
	}

	skippingReason, err := c.checkIfBackupShouldBeSkipped(inv, session)
	if err != nil {
		return err
	}

	if skippingReason != "" {
		klog.Infof(skippingReason)
		return conditions.SetBackupSkippedConditionToTrue(session, skippingReason)
	}

	if !globalPreBackupHookExecuted(inv, session) {
		err = c.executeGlobalPreBackupHook(inv, session)
		if err != nil {
			return conditions.SetGlobalPreBackupHookSucceededConditionToFalse(session, err)
		}
	}

	// ===================== Run Backup for the Individual Targets ============================
	shouldRequeue := false
	for i, targetInfo := range inv.GetTargetInfo() {
		if targetInfo.Target != nil {
			// Skip processing if the backup has been already initiated before for this target
			if invoker.TargetBackupInitiated(targetInfo.Target.Ref, session.GetTargetStatus()) {
				klog.Infof("Skipping initiating backup by BackupSession %s/%s for target %s %s/%s. Reason: Backup already initiated.",
					session.GetObjectMeta().Namespace,
					session.GetObjectMeta().Name,
					targetInfo.Target.Ref.Kind,
					targetInfo.Target.Ref.Namespace,
					targetInfo.Target.Ref.Name,
				)
				continue
			}

			pendingReason, err := c.checkIfBackupShouldBePending(inv, session, targetInfo.Target.Ref)
			if err != nil {
				return err
			}
			if pendingReason != "" {
				klog.Infof("Skipping initiating backup for BackupSession %s/%s. Reason: %s.",
					session.GetObjectMeta().Namespace,
					session.GetObjectMeta().Name,
					pendingReason,
				)
				err = c.setTargetBackupPending(targetInfo.Target.Ref, session)
				if err != nil {
					return err
				}
				shouldRequeue = true
				continue
			}

			err = c.ensureBackupExecutor(inv, targetInfo, session, i)
			if err != nil {
				msg := fmt.Sprintf("failed to ensure backup executor. Reason: %v", err)
				klog.Warning(msg, fmt.Sprint("target: ", targetInfo.Target.Ref))
				return conditions.SetBackupExecutorEnsuredToFalse(session, targetInfo.Target.Ref, err)
			}

			// Set target backup phase to "Running"
			err = c.initiateTargetBackup(inv, session, i)
			if err != nil {
				return err
			}
		}
	}
	if shouldRequeue {
		return c.requeueBackupSession(session, requeueTimeInterval)
	}

	if inv.GetTimeOut() != "" {
		if err := c.requeueBackupAfterTimeOut(session, inv.GetTimeOut()); err != nil {
			return err
		}
	}

	return nil
}

func isBackupDeadlineSet(session *invoker.BackupSessionHandler) bool {
	deadline := session.GetStatus().SessionDeadline
	return !deadline.IsZero()
}

func (c *StashController) requeueBackupAfterTimeOut(session *invoker.BackupSessionHandler, timeOut string) error {
	klog.Infof("\nTimeOut set for BackupSession %s/%s",
		session.GetObjectMeta().Namespace,
		session.GetObjectMeta().Name,
	)
	if !isBackupDeadlineSet(session) {
		timeOut, err := time.ParseDuration(timeOut)
		if err != nil {
			return err
		}
		if err := c.requeueBackupSession(session, timeOut); err != nil {
			return err
		}
		return c.setBackupDeadline(session, timeOut)
	}
	return nil
}

func isBackupDeadlineExceeded(session *invoker.BackupSessionHandler) bool {
	return metav1.Now().After(session.GetStatus().SessionDeadline.Time)
}

func (c *StashController) requeueBackupSession(session *invoker.BackupSessionHandler, timeOut time.Duration) error {
	klog.Infof("Requeueing BackupSession %s/%s after %s seconds.....",
		session.GetObjectMeta().Namespace,
		session.GetObjectMeta().Name,
		timeOut.String(),
	)

	key, err := cache.MetaNamespaceKeyFunc(session.GetBackupSession())
	if err != nil {
		return err
	}
	c.backupSessionQueue.GetQueue().AddAfter(key, timeOut)
	return nil
}

func (c *StashController) ensureBackupExecutor(inv invoker.BackupInvoker, targetInfo invoker.BackupTargetInfo, session *invoker.BackupSessionHandler, idx int) error {
	switch backupExecutor(inv, targetInfo) {
	case BackupExecutorSidecar:
		err := c.ensureLatestSidecarConfiguration(targetInfo)
		if err != nil {
			return err
		}
	case BackupExecutorCSIDriver:
		// VolumeSnapshotter driver has been used. So, ensure VolumeSnapshotter job
		err := c.ensureVolumeSnapshotterJob(inv, targetInfo, session, idx)
		if err != nil {
			return err
		}
	case BackupExecutorJob:
		err := c.ensureBackupJob(inv, targetInfo, session, idx)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unable to identify backup executor entity")
	}

	return conditions.SetBackupExecutorEnsuredToTrue(session, targetInfo.Target.Ref)
}

func (c *StashController) ensureBackupJob(inv invoker.BackupInvoker, targetInfo invoker.BackupTargetInfo, session *invoker.BackupSessionHandler, index int) error {
	invMeta := inv.GetObjectMeta()
	ownerRef := inv.GetOwnerRef()
	invRef, err := inv.GetObjectRef()
	runtimeSettings := targetInfo.RuntimeSettings

	if err != nil {
		return err
	}
	jobMeta := metav1.ObjectMeta{
		Name:      getBackupJobName(session, strconv.Itoa(index)),
		Namespace: session.GetObjectMeta().Namespace,
		Labels:    inv.GetLabels(),
	}

	psps, err := c.getBackupJobPSPNames(targetInfo.Task)
	if err != nil {
		return err
	}

	rbacOptions, err := c.getBackupRBACOptions(inv, &index)
	if err != nil {
		return err
	}
	rbacOptions.PodSecurityPolicyNames = psps

	if runtimeSettings.Pod != nil {
		if runtimeSettings.Pod.ServiceAccountName != "" {
			rbacOptions.ServiceAccount.Name = runtimeSettings.Pod.ServiceAccountName
		}

		if runtimeSettings.Pod.ServiceAccountAnnotations != nil {
			rbacOptions.ServiceAccount.Annotations = runtimeSettings.Pod.ServiceAccountAnnotations
		}
	}

	err = rbacOptions.EnsureBackupJobRBAC()
	if err != nil {
		return err
	}

	// if the Stash is using a private registry, then ensure the image pull secrets
	var imagePullSecrets []core.LocalObjectReference
	if c.ImagePullSecrets != nil {
		imagePullSecrets, err = c.ensureImagePullSecrets(invMeta, ownerRef)
		if err != nil {
			return err
		}
	}

	// get repository for backupConfig
	repository, err := inv.GetRepository()
	if err != nil {
		return err
	}

	addon, err := api_util.ExtractAddonInfo(c.appCatalogClient, targetInfo.Task, targetInfo.Target.Ref)
	if err != nil {
		return err
	}

	// resolve task template
	repoInputs, err := c.inputsForRepository(repository)
	if err != nil {
		return fmt.Errorf("cannot resolve implicit inputs for Repository %s/%s, reason: %s", repository.Namespace, repository.Name, err)
	}

	bcInputs, err := c.inputsForBackupInvoker(inv, targetInfo)
	if err != nil {
		return fmt.Errorf("cannot resolve implicit inputs for backup invoker  %s %s/%s, reason: %s", inv.GetTypeMeta().Kind, invMeta.Namespace, invMeta.Name, err)
	}

	implicitInputs := meta_util.OverwriteKeys(repoInputs, bcInputs)
	implicitInputs[apis.Namespace] = session.GetObjectMeta().Namespace
	implicitInputs[apis.BackupSession] = session.GetObjectMeta().Name

	// add docker image specific input
	implicitInputs[apis.StashDockerRegistry] = c.DockerRegistry
	implicitInputs[apis.StashDockerImage] = c.StashImage
	implicitInputs[apis.StashImageTag] = c.StashImageTag
	// license related inputs
	implicitInputs[apis.LicenseApiService] = c.LicenseApiService

	taskResolver := resolve.TaskResolver{
		StashClient:     c.stashClient,
		TaskName:        addon.BackupTask.Name,
		Inputs:          meta_util.OverwriteKeys(explicitInputs(addon.BackupTask.Params), implicitInputs), // TODO: reverse priority ???
		RuntimeSettings: targetInfo.RuntimeSettings,
		TempDir:         targetInfo.TempDir,
	}

	// if preBackup or postBackup Hook is specified, add their specific inputs
	if targetInfo.Hooks != nil && targetInfo.Hooks.PreBackup != nil {
		taskResolver.PreTaskHookInput = make(map[string]string)
		taskResolver.PreTaskHookInput[apis.HookType] = apis.PreBackupHook
	}
	if targetInfo.Hooks != nil &&
		targetInfo.Hooks.PostBackup != nil &&
		targetInfo.Hooks.PostBackup.Handler != nil {
		taskResolver.PostTaskHookInput = make(map[string]string)
		taskResolver.PostTaskHookInput[apis.HookType] = apis.PostBackupHook
	}

	podSpec, err := taskResolver.GetPodSpec(invRef.Kind, invRef.Name, targetInfo.Target.Ref)
	if err != nil {
		return fmt.Errorf("can't get PodSpec for backup invoker %s/%s, reason: %s", invMeta.Namespace, invMeta.Name, err)
	}

	ownerBackupSession := metav1.NewControllerRef(session.GetBackupSession(), api_v1beta1.SchemeGroupVersion.WithKind(api_v1beta1.ResourceKindBackupSession))

	// upsert InterimVolume to hold the backup/restored data temporarily
	podSpec, err = util.UpsertInterimVolume(c.kubeClient, podSpec, targetInfo.InterimVolumeTemplate.ToCorePVC(), invMeta.Namespace, ownerBackupSession)
	if err != nil {
		return err
	}

	// create Backup Job
	_, _, err = batch_util.CreateOrPatchJob(
		context.TODO(),
		c.kubeClient,
		jobMeta,
		func(in *batch.Job) *batch.Job {
			// set BackupSession as owner of this Job so that it get cleaned automatically
			// when the BackupSession gets deleted according to backupHistoryLimit
			core_util.EnsureOwnerReference(&in.ObjectMeta, ownerBackupSession)

			in.Spec.Template.Spec = podSpec
			in.Spec.Template.Labels = meta_util.OverwriteKeys(in.Spec.Template.Labels, inv.GetLabels())
			in.Spec.Template.Spec.ImagePullSecrets = core_util.MergeLocalObjectReferences(in.Spec.Template.Spec.ImagePullSecrets, imagePullSecrets)
			in.Spec.Template.Spec.ServiceAccountName = rbacOptions.ServiceAccount.Name
			in.Spec.BackoffLimit = pointer.Int32P(0)
			if runtimeSettings.Pod != nil && runtimeSettings.Pod.PodAnnotations != nil {
				in.Spec.Template.Annotations = runtimeSettings.Pod.PodAnnotations
			}
			return in
		},
		metav1.PatchOptions{},
	)

	return err
}

func (c *StashController) ensureVolumeSnapshotterJob(inv invoker.BackupInvoker, targetInfo invoker.BackupTargetInfo, session *invoker.BackupSessionHandler, index int) error {
	invMeta := inv.GetObjectMeta()
	ownerRef := inv.GetOwnerRef()
	runtimeSettings := targetInfo.RuntimeSettings

	jobMeta := metav1.ObjectMeta{
		Name:      getVolumeSnapshotterJobName(targetInfo.Target.Ref, session.GetObjectMeta().Name),
		Namespace: session.GetObjectMeta().Namespace,
		Labels:    inv.GetLabels(),
	}

	rbacOptions, err := c.getBackupRBACOptions(inv, &index)
	if err != nil {
		return err
	}

	if runtimeSettings.Pod != nil {
		if runtimeSettings.Pod.ServiceAccountName != "" {
			rbacOptions.ServiceAccount.Name = runtimeSettings.Pod.ServiceAccountName
		}

		if runtimeSettings.Pod.ServiceAccountAnnotations != nil {
			rbacOptions.ServiceAccount.Annotations = runtimeSettings.Pod.ServiceAccountAnnotations
		}
	}

	err = rbacOptions.EnsureVolumeSnapshotterJobRBAC()
	if err != nil {
		return err
	}

	// if the Stash is using a private registry, then ensure the image pull secrets
	var imagePullSecrets []core.LocalObjectReference
	if c.ImagePullSecrets != nil {
		imagePullSecrets, err = c.ensureImagePullSecrets(invMeta, ownerRef)
		if err != nil {
			return err
		}
	}

	image := docker.Docker{
		Registry: c.DockerRegistry,
		Image:    c.StashImage,
		Tag:      c.StashImageTag,
	}

	jobTemplate, err := util.NewVolumeSnapshotterJob(session, targetInfo.Target, targetInfo.RuntimeSettings, image)
	if err != nil {
		return err
	}

	ownerBackupSession := metav1.NewControllerRef(session.GetBackupSession(), api_v1beta1.SchemeGroupVersion.WithKind(api_v1beta1.ResourceKindBackupSession))
	// Create VolumeSnapshotter job
	_, _, err = batch_util.CreateOrPatchJob(
		context.TODO(),
		c.kubeClient,
		jobMeta,
		func(in *batch.Job) *batch.Job {
			// set BackupSession as owner of this Job so that it get cleaned automatically
			// when the BackupSession gets deleted according to backupHistoryLimit
			core_util.EnsureOwnerReference(&in.ObjectMeta, ownerBackupSession)

			in.Spec.Template.Spec = jobTemplate.Spec
			in.Spec.Template.Labels = meta_util.OverwriteKeys(in.Spec.Template.Labels, inv.GetLabels())
			in.Spec.Template.Spec.ImagePullSecrets = core_util.MergeLocalObjectReferences(in.Spec.Template.Spec.ImagePullSecrets, imagePullSecrets)
			in.Spec.Template.Spec.ServiceAccountName = rbacOptions.ServiceAccount.Name
			in.Spec.BackoffLimit = pointer.Int32P(0)
			if runtimeSettings.Pod != nil && runtimeSettings.Pod.PodAnnotations != nil {
				in.Spec.Template.Annotations = runtimeSettings.Pod.PodAnnotations
			}
			return in
		},
		metav1.PatchOptions{},
	)

	return err
}

func (c *StashController) setTargetBackupPending(targetRef api_v1beta1.TargetRef, session *invoker.BackupSessionHandler) error {
	return session.UpdateStatus(&api_v1beta1.BackupSessionStatus{
		Targets: []api_v1beta1.BackupTargetStatus{
			{
				Ref: targetRef,
			},
		},
	})
}

func (c *StashController) setBackupDeadline(session *invoker.BackupSessionHandler, timeOut time.Duration) error {
	return session.UpdateStatus(&api_v1beta1.BackupSessionStatus{
		SessionDeadline: metav1.NewTime(session.GetObjectMeta().CreationTimestamp.Add(timeOut)),
	})
}

func (c *StashController) initiateTargetBackup(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler, index int) error {
	targetsInfo := inv.GetTargetInfo()
	target := targetsInfo[index].Target
	// find out the total number of hosts in target that will be backed up in this backup session
	totalHosts, err := c.getTotalHosts(target, inv.GetDriver())
	if err != nil {
		return err
	}
	// For Restic driver, set preBackupAction and postBackupAction
	var preBackupActions, postBackupActions []string
	if inv.GetDriver() == api_v1beta1.ResticSnapshotter {
		// assign preBackupAction to the first target
		if index == 0 {
			preBackupActions = []string{api_v1beta1.InitializeBackendRepository}
		}
		// assign postBackupAction to the last target
		if index == len(targetsInfo)-1 {
			postBackupActions = []string{api_v1beta1.ApplyRetentionPolicy, api_v1beta1.VerifyRepositoryIntegrity, api_v1beta1.SendRepositoryMetrics}
		}
	}
	if session.GetStatus().Phase != api_v1beta1.BackupSessionRunning {
		klog.Infof("Initiating backup for BackupSession %s/%s", session.GetObjectMeta().Namespace, session.GetObjectMeta().Name)
	}
	return session.UpdateStatus(&api_v1beta1.BackupSessionStatus{
		Targets: []api_v1beta1.BackupTargetStatus{
			{
				TotalHosts:        totalHosts,
				Ref:               target.Ref,
				PreBackupActions:  preBackupActions,
				PostBackupActions: postBackupActions,
			},
		},
	})
}

func getBackupJobName(session *invoker.BackupSessionHandler, index string) string {
	return meta.ValidNameWithPrefixNSuffix(apis.PrefixStashBackup, strings.ReplaceAll(session.GetObjectMeta().Name, ".", "-"), index)
}

func getVolumeSnapshotterJobName(targetRef api_v1beta1.TargetRef, name string) string {
	parts := strings.Split(name, "-")
	suffix := parts[len(parts)-1]
	return meta.ValidNameWithPrefix(apis.PrefixStashVolumeSnapshot, fmt.Sprintf("%s-%s-%s", util.ResourceKindShortForm(targetRef.Kind), targetRef.Name, suffix))
}

func backupHistoryCleaned(conditions []kmapi.Condition) bool {
	return kmapi.HasCondition(conditions, api_v1beta1.BackupHistoryCleaned)
}

// cleanupBackupHistory deletes old BackupSessions and theirs associate resources according to BackupHistoryLimit
func (c *StashController) cleanupBackupHistory(session *invoker.BackupSessionHandler, backupHistoryLimit *int32) error {
	// default history limit is 1
	historyLimit := int32(1)
	if backupHistoryLimit != nil {
		historyLimit = *backupHistoryLimit
	}

	// BackupSession use BackupConfiguration name as label. We can use this label as selector to list only the BackupSession
	// of this particular BackupConfiguration.
	label := metav1.LabelSelector{
		MatchLabels: map[string]string{
			apis.LabelInvokerType: session.GetInvokerRef().Kind,
			apis.LabelInvokerName: session.GetInvokerRef().Name,
		},
	}
	selector, err := metav1.LabelSelectorAsSelector(&label)
	if err != nil {
		return err
	}

	// list all the BackupSessions of this particular BackupConfiguration
	bsList, err := c.backupSessionLister.BackupSessions(session.GetObjectMeta().Namespace).List(selector)
	if err != nil {
		return err
	}

	// sort BackupSession according to creation timestamp. keep latest BackupSession first.
	sort.Slice(bsList, func(i, j int) bool {
		return bsList[i].CreationTimestamp.After(bsList[j].CreationTimestamp.Time)
	})

	var lastCompletedSession string
	for i := range bsList {
		if bsList[i].Status.Phase == api_v1beta1.BackupSessionSucceeded || bsList[i].Status.Phase == api_v1beta1.BackupSessionFailed {
			lastCompletedSession = bsList[i].Name
			break
		}
	}
	// delete the BackupSession that does not fit within the history limit
	for i := int(historyLimit); i < len(bsList); i++ {
		if invoker.IsBackupCompleted(bsList[i].Status.Phase) && !(bsList[i].Name == lastCompletedSession && historyLimit > 0) {
			err = c.stashClient.StashV1beta1().BackupSessions(session.GetObjectMeta().Namespace).Delete(context.TODO(), bsList[i].Name, meta.DeleteInBackground())
			if err != nil && !(kerr.IsNotFound(err) || kerr.IsGone(err)) {
				return err
			}
		}
	}
	return conditions.SetBackupHistoryCleanedConditionToTrue(session)
}

func backupExecutor(inv invoker.BackupInvoker, targetInfo invoker.BackupTargetInfo) string {
	if inv.GetDriver() == api_v1beta1.ResticSnapshotter &&
		util.BackupModel(targetInfo.Target.Ref.Kind, targetInfo.Task.Name) == apis.ModelSidecar {
		return BackupExecutorSidecar
	}
	if inv.GetDriver() == api_v1beta1.VolumeSnapshotter {
		return BackupExecutorCSIDriver
	}
	return BackupExecutorJob
}

func postBackupHooksExecuted(targetInfo []invoker.BackupTargetInfo, targetStatus []api_v1beta1.BackupTargetStatus) bool {
	for _, target := range targetInfo {
		if target.Hooks != nil && target.Hooks.PostBackup != nil {
			if !postBackupHookExecutedForTarget(target, targetStatus) {
				return false
			}
		}
	}
	return true
}

func postBackupHookExecutedForTarget(targetInfo invoker.BackupTargetInfo, targetStatus []api_v1beta1.BackupTargetStatus) bool {
	if targetInfo.Target == nil {
		return true
	}

	for _, s := range targetStatus {
		if invoker.TargetMatched(s.Ref, targetInfo.Target.Ref) {
			if kmapi.HasCondition(s.Conditions, api_v1beta1.PostBackupHookExecutionSucceeded) {
				return true
			}
		}
	}
	return false
}

func globalPostBackupHookExecuted(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler) bool {
	backupHooks := inv.GetGlobalHooks()
	if backupHooks == nil ||
		backupHooks.PostBackup == nil ||
		backupHooks.PostBackup.Handler == nil {
		return true
	}
	return kmapi.HasCondition(session.GetConditions(), api_v1beta1.GlobalPostBackupHookSucceeded) &&
		kmapi.IsConditionTrue(session.GetConditions(), api_v1beta1.GlobalPostBackupHookSucceeded)
}

func (c *StashController) executeGlobalPostBackupHook(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler) error {
	summary := inv.GetSummary(api_v1beta1.TargetRef{}, kmapi.ObjectReference{
		Namespace: session.GetObjectMeta().Namespace,
		Name:      session.GetObjectMeta().Name,
	})
	summary.Status.Phase = string(session.GetStatus().Phase)
	summary.Status.Duration = session.GetStatus().SessionDuration

	hookExecutor := stashHooks.HookExecutor{
		Config: c.clientConfig,
		Hook:   inv.GetGlobalHooks().PostBackup.Handler,
		ExecutorPod: kmapi.ObjectReference{
			Namespace: meta.PodNamespace(),
			Name:      meta.PodName(),
		},
		Summary: summary,
	}

	executionPolicy := inv.GetGlobalHooks().PostBackup.ExecutionPolicy
	if executionPolicy == "" {
		executionPolicy = api_v1beta1.ExecuteAlways
	}

	if !stashHooks.IsAllowedByExecutionPolicy(executionPolicy, summary) {
		reason := fmt.Sprintf("Skipping executing %s. Reason: executionPolicy is %q but phase is %q.",
			apis.PostBackupHook,
			executionPolicy,
			summary.Status.Phase,
		)
		return conditions.SetGlobalPostBackupHookSucceededConditionToTrueWithMsg(session, reason)
	}
	if err := hookExecutor.Execute(); err != nil {
		return err
	}
	return conditions.SetGlobalPostBackupHookSucceededConditionToTrue(session)
}

func globalPreBackupHookExecuted(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler) bool {
	backupHooks := inv.GetGlobalHooks()
	if backupHooks == nil || backupHooks.PreBackup == nil {
		return true
	}
	return kmapi.HasCondition(session.GetConditions(), api_v1beta1.GlobalPreBackupHookSucceeded) &&
		kmapi.IsConditionTrue(session.GetConditions(), api_v1beta1.GlobalPreBackupHookSucceeded)
}

func (c *StashController) executeGlobalPreBackupHook(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler) error {
	hookExecutor := stashHooks.HookExecutor{
		Config: c.clientConfig,
		Hook:   inv.GetGlobalHooks().PreBackup,
		ExecutorPod: kmapi.ObjectReference{
			Namespace: meta.PodNamespace(),
			Name:      meta.PodName(),
		},
		Summary: inv.GetSummary(api_v1beta1.TargetRef{}, kmapi.ObjectReference{
			Namespace: session.GetObjectMeta().Namespace,
			Name:      session.GetObjectMeta().Name,
		}),
	}
	if err := hookExecutor.Execute(); err != nil {
		return err
	}
	return conditions.SetGlobalPreBackupHookSucceededConditionToTrue(session)
}

func (c *StashController) checkIfBackupShouldBeSkipped(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler) (string, error) {
	// Skip if the respective backup invoker is not in ready state
	if inv.GetPhase() != api_v1beta1.BackupInvokerReady {
		return fmt.Sprintf("Skipped taking backup. Reason: %s %s/%s is not ready.",
			inv.GetTypeMeta().Kind,
			inv.GetObjectMeta().Namespace,
			inv.GetObjectMeta().Name,
		), nil
	}

	// Skip taking backup if there is another running BackupSession
	runningBS, err := c.checkForAnotherRunningBackupSessionWithSameInvoker(inv, session)
	if err != nil {
		return "", err
	}
	if isBackupPending(session) && runningBS != nil {
		return fmt.Sprintf("Skipped taking new backup. Reason: Previous BackupSession: %s is %q.",
			runningBS.Name,
			runningBS.Status.Phase,
		), nil
	}

	// do not skip
	return "", nil
}

func (c *StashController) checkIfBackupShouldBePending(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler, targetRef api_v1beta1.TargetRef) (string, error) {
	// Keep backup pending if the target is not in next in order
	if inv.GetExecutionOrder() == api_v1beta1.Sequential &&
		!inv.NextInOrder(targetRef, session.GetTargetStatus()) {
		return "Backup order is sequential and some previous targets hasn't completed their backup process.", nil
	}

	yes, err := c.hasMultipleBackupInvokers(targetRef)
	if err != nil || !yes {
		return "", err
	}

	otherSession, err := c.checkForAnotherIncompleteBackupSessionWithDifferentInvoker(inv, targetRef)
	if err != nil {
		return "", err
	}
	if otherSession == nil {
		if time.Since(session.GetObjectMeta().CreationTimestamp.Time) < 2*time.Second {
			return "Multiple backup invoker found. Will be processed in the next requeue to handle concurrent session.", nil
		}
		return "", nil
	}

	if shouldKeepCurrentSessionPending(session.GetBackupSession(), otherSession) {
		return fmt.Sprintf("Found another incomplete BackupSession %s/%s invoked by %s/%s",
			otherSession.Namespace,
			otherSession.Name,
			otherSession.Spec.Invoker.Kind,
			otherSession.Spec.Invoker.Name,
		), nil
	}
	return "", nil
}

func shouldKeepCurrentSessionPending(cur, other *api_v1beta1.BackupSession) bool {
	if other.Status.Phase == api_v1beta1.BackupSessionRunning {
		return true
	}

	if cur.CreationTimestamp.Equal(other.CreationTimestamp.DeepCopy()) {
		return other.Name < cur.Name
	}
	return cur.CreationTimestamp.After(other.CreationTimestamp.Time)
}

func (c *StashController) hasMultipleBackupInvokers(targetRef api_v1beta1.TargetRef) (bool, error) {
	invokers, err := util.FindBackupInvokers(c.bcLister, targetRef)
	if err != nil {
		return false, err
	}

	return len(invokers) > 1, nil
}

func isSessionSkipped(session *invoker.BackupSessionHandler) bool {
	return kmapi.IsConditionTrue(session.GetConditions(), api_v1beta1.BackupSkipped)
}

func isBackupRunning(session *invoker.BackupSessionHandler) bool {
	return session.GetStatus().Phase == api_v1beta1.BackupSessionRunning
}

func isBackupPending(session *invoker.BackupSessionHandler) bool {
	return session.GetStatus().Phase == "" || session.GetStatus().Phase == api_v1beta1.BackupSessionPending
}

func (c *StashController) checkForAnotherRunningBackupSessionWithSameInvoker(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler) (*api_v1beta1.BackupSession, error) {
	runningBS, err := c.getRunningBackupSessionForInvoker(inv)
	if err != nil {
		return nil, err
	}
	if runningBS != nil && runningBS.Name != session.GetObjectMeta().Name {
		return runningBS, nil
	}
	return nil, nil
}

func (c *StashController) getRunningBackupSessionForInvoker(inv invoker.BackupInvoker) (*api_v1beta1.BackupSession, error) {
	backupSessions, err := c.backupSessionLister.BackupSessions(inv.GetObjectMeta().Namespace).List(labels.SelectorFromSet(map[string]string{
		apis.LabelInvokerName: inv.GetObjectMeta().Name,
		apis.LabelInvokerType: inv.GetTypeMeta().Kind,
	}))
	if err != nil {
		return nil, err
	}
	for i := range backupSessions {
		if backupSessions[i].Status.Phase == api_v1beta1.BackupSessionRunning {
			return backupSessions[i], nil
		}
	}
	return nil, nil
}

func (c *StashController) checkForAnotherIncompleteBackupSessionWithDifferentInvoker(inv invoker.BackupInvoker, targetRef api_v1beta1.TargetRef) (*api_v1beta1.BackupSession, error) {
	sessions, err := c.getIncompleteBackupSessionForTarget(inv, targetRef)
	if err != nil {
		return nil, err
	}

	if len(sessions) > 1 {
		sort.Slice(sessions, func(i, j int) bool {
			if sessions[i].Status.Phase == sessions[j].Status.Phase {
				return sessions[i].Name < sessions[j].Name
			}
			if sessions[i].Status.Phase == api_v1beta1.BackupSessionRunning && sessions[j].Status.Phase != api_v1beta1.BackupSessionRunning {
				return true
			}
			return false
		})
	}

	for i, s := range sessions {
		if invokedByDifferentInvoker(inv, s.Spec.Invoker) {
			s := sessions[i]
			return &s, nil
		}
	}
	return nil, nil
}

func invokedByDifferentInvoker(inv invoker.BackupInvoker, invRef api_v1beta1.BackupInvokerRef) bool {
	return invRef.Kind != inv.GetTypeMeta().Kind ||
		invRef.Name != inv.GetObjectMeta().Name
}

func (c *StashController) getIncompleteBackupSessionForTarget(inv invoker.BackupInvoker, targetRef api_v1beta1.TargetRef) ([]api_v1beta1.BackupSession, error) {
	selector := labels.SelectorFromSet(map[string]string{
		apis.LabelTargetKind:      targetRef.Kind,
		apis.LabelTargetName:      targetRef.Name,
		apis.LabelTargetNamespace: targetRef.Namespace,
	})
	bsList, err := c.stashClient.StashV1beta1().BackupSessions(inv.GetObjectMeta().Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}

	sessions := make([]api_v1beta1.BackupSession, 0)
	for i := range bsList.Items {
		if bsList.Items[i].Status.Phase == api_v1beta1.BackupSessionRunning ||
			bsList.Items[i].Status.Phase == api_v1beta1.BackupSessionPending ||
			bsList.Items[i].Status.Phase == "" {
			sessions = append(sessions, bsList.Items[i])
		}
	}
	return sessions, nil
}

func explicitInputs(params []appcat.Param) map[string]string {
	inputs := make(map[string]string)
	for _, param := range params {
		inputs[param.Name] = param.Value
	}
	return inputs
}

func (c *StashController) getBackupRBACOptions(inv invoker.BackupInvoker, index *int) (stash_rbac.RBACOptions, error) {
	invMeta := inv.GetObjectMeta()
	repo := inv.GetRepoRef()

	rbacOptions := stash_rbac.RBACOptions{
		KubeClient: c.kubeClient,
		Invoker: stash_rbac.InvokerOptions{
			ObjectMeta: invMeta,
			TypeMeta:   inv.GetTypeMeta(),
		},
		Owner:          inv.GetOwnerRef(),
		OffshootLabels: inv.GetLabels(),
		ServiceAccount: metav1.ObjectMeta{
			Namespace: invMeta.Namespace,
		},
	}

	if repo.Namespace != invMeta.Namespace {
		repository, err := c.repoLister.Repositories(repo.Namespace).Get(repo.Name)
		if err != nil {
			if kerr.IsNotFound(err) {
				return rbacOptions, nil
			}
			return rbacOptions, err
		}

		rbacOptions.CrossNamespaceResources = &stash_rbac.CrossNamespaceResources{
			Repository: repo.Name,
			Namespace:  repo.Namespace,
			Secret:     repository.Spec.Backend.StorageSecretName,
		}
	}
	rbacOptions.Suffix = "0"
	if index != nil {
		rbacOptions.Suffix = fmt.Sprintf("%d", *index)
	}

	return rbacOptions, nil
}

func backupMetricPushed(conditions []kmapi.Condition) bool {
	return kmapi.HasCondition(conditions, api_v1beta1.MetricsPushed)
}

func (c *StashController) sendBackupMetrics(inv invoker.BackupInvoker, session *invoker.BackupSessionHandler) error {
	metricsOpt := &metrics.MetricsOptions{
		Enabled:        true,
		PushgatewayURL: metrics.GetPushgatewayURL(),
		JobName:        fmt.Sprintf("%s-%s-%s", strings.ToLower(inv.GetTypeMeta().Kind), inv.GetObjectMeta().Namespace, inv.GetObjectMeta().Name),
	}

	status := session.GetStatus()
	if status.SessionDuration == "" {
		status.SessionDuration = time.Since(session.GetObjectMeta().CreationTimestamp.Time).Round(time.Second).String()
	}

	// send backup session related metrics
	err := metricsOpt.SendBackupSessionMetrics(inv, status)
	if err != nil {
		return err
	}
	// send target related metrics
	for _, target := range session.GetTargetStatus() {
		err = metricsOpt.SendBackupTargetMetrics(c.clientConfig, inv, target.Ref, session.GetStatus())
		if err != nil {
			return err
		}
	}

	return conditions.SetBackupMetricsPushedConditionToTrue(session)
}
