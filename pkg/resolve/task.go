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

package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"stash.appscode.dev/apimachinery/apis"
	v1beta1_api "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/stash/pkg/util"

	"gomodules.xyz/envsubst"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_util "kmodules.xyz/client-go/meta"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	ofst_util "kmodules.xyz/offshoot-api/util"
)

type TaskResolver struct {
	StashClient       cs.Interface
	TaskName          string
	Inputs            map[string]string
	RuntimeSettings   ofst.RuntimeSettings
	TempDir           v1beta1_api.EmptyDirSettings
	PreTaskHookInput  map[string]string
	PostTaskHookInput map[string]string
}

func (o TaskResolver) GetPodSpec(invokerType, invokerName, targetKind, targetName string) (core.PodSpec, error) {
	task, err := o.StashClient.StashV1beta1().Tasks().Get(context.TODO(), o.TaskName, metav1.GetOptions{})
	if err != nil {
		return core.PodSpec{}, err
	}
	// resolve Task with inputs, modify in place
	if err = resolveWithInputs(task, o.Inputs); err != nil {
		return core.PodSpec{}, err
	}

	var containers []core.Container
	// User may overwrite some variables (i.e. outputDir) of hook executor container in Task params
	// We need to substitute these variables properly. Params of last Function will have higher precedence
	taskParams := make(map[string]string)

	// get Functions for Task
	for i, fn := range task.Spec.Steps {
		function, err := o.StashClient.StashV1beta1().Functions().Get(context.TODO(), fn.Name, metav1.GetOptions{})
		if err != nil {
			return core.PodSpec{}, fmt.Errorf("can't get Function %s for Task %s, reason: %s", fn.Name, task.Name, err)
		}

		// inputs from params
		inputs := make(map[string]string)
		for _, param := range fn.Params {
			inputs[param.Name] = param.Value
		}
		taskParams = meta_util.OverwriteKeys(taskParams, inputs)

		// merge/replace backup config inputs
		inputs = meta_util.OverwriteKeys(inputs, o.Inputs)

		//Add addon image as input
		inputs = meta_util.OverwriteKeys(inputs, map[string]string{
			apis.AddonImage: function.Spec.Image,
		})

		// resolve Function with inputs, modify in place
		if err = resolveWithInputs(function, inputs); err != nil {
			return core.PodSpec{}, fmt.Errorf("can't resolve Function %s for Task %s, reason: %s", fn.Name, task.Name, err)
		}

		// init ContainerRuntimeSettings to avoid nil pointer
		if function.Spec.RuntimeSettings == nil {
			function.Spec.RuntimeSettings = &ofst.ContainerRuntimeSettings{}
		}

		// container from function spec
		container := core.Container{
			Name:            fmt.Sprintf("%s-%d", strings.ReplaceAll(function.Name, ".", "-"), i), // TODO
			Image:           function.Spec.Image,
			Command:         function.Spec.Command,
			Args:            function.Spec.Args,
			WorkingDir:      function.Spec.WorkingDir,
			Ports:           function.Spec.Ports,
			EnvFrom:         function.Spec.RuntimeSettings.EnvFrom,
			Env:             function.Spec.RuntimeSettings.Env,
			VolumeMounts:    function.Spec.VolumeMounts,
			VolumeDevices:   function.Spec.VolumeDevices,
			Resources:       function.Spec.RuntimeSettings.Resources,
			LivenessProbe:   function.Spec.RuntimeSettings.LivenessProbe,
			ReadinessProbe:  function.Spec.RuntimeSettings.ReadinessProbe,
			Lifecycle:       function.Spec.RuntimeSettings.Lifecycle,
			SecurityContext: function.Spec.RuntimeSettings.SecurityContext,
			ImagePullPolicy: core.PullIfNotPresent,
		}

		// mount tmp volume
		container.VolumeMounts = util.UpsertTmpVolumeMount(container.VolumeMounts)

		// apply RuntimeSettings to Container
		if o.RuntimeSettings.Container != nil {
			container = ofst_util.ApplyContainerRuntimeSettings(container, *o.RuntimeSettings.Container)
		}

		containers = append(containers, container)
	}
	if len(containers) == 0 {
		return core.PodSpec{}, fmt.Errorf("empty steps/containers for Task %s", task.Name)
	}
	// if hook specified then, add hook executor containers
	if o.PreTaskHookInput != nil {
		// Inputs precedence:
		// 1. Inputs from BackupConfiguration/RestoreSession
		// 2. Inputs from Task params
		// 3. Default hook specific inputs
		inputs := meta_util.OverwriteKeys(taskParams, o.Inputs)
		inputs = meta_util.OverwriteKeys(o.PreTaskHookInput, inputs)
		hookExecutor := util.HookExecutorContainer(apis.PreTaskHook, containers, invokerType, invokerName, targetKind, targetName)

		if err = resolveWithInputs(&hookExecutor, inputs); err != nil {
			return core.PodSpec{}, fmt.Errorf("failed to resolve preTaskHook. Reason: %v", err)
		}

		// apply RuntimeSettings to Container
		if o.RuntimeSettings.Container != nil {
			hookExecutor = ofst_util.ApplyContainerRuntimeSettings(hookExecutor, *o.RuntimeSettings.Container)
		}

		containers = append([]core.Container{hookExecutor}, containers...)
	}

	if o.PostTaskHookInput != nil {
		inputs := meta_util.OverwriteKeys(taskParams, o.Inputs)
		inputs = meta_util.OverwriteKeys(o.PostTaskHookInput, inputs)
		hookExecutor := util.HookExecutorContainer(apis.PostTaskHook, containers, invokerType, invokerName, targetKind, targetName)

		if err = resolveWithInputs(&hookExecutor, inputs); err != nil {
			return core.PodSpec{}, fmt.Errorf("failed to resolve postTaskHook. Reason: %v", err)
		}

		// apply RuntimeSettings to Container
		if o.RuntimeSettings.Container != nil {
			hookExecutor = ofst_util.ApplyContainerRuntimeSettings(hookExecutor, *o.RuntimeSettings.Container)
		}

		lastContainer := containers[len(containers)-1]
		containers[len(containers)-1] = hookExecutor
		containers = append(containers, lastContainer)
	}
	// podSpec from task
	podSpec := core.PodSpec{
		Volumes:        task.Spec.Volumes,
		InitContainers: containers[:len(containers)-1],
		Containers:     containers[len(containers)-1:],
		RestartPolicy:  core.RestartPolicyNever, // TODO: use OnFailure ?
	}

	// The "update-status" Function uses the same image as the Stash operator which runs as user 65535.
	// However, the backup container runs as root user if no securityContext is provided. As a result
	// the "update-status" container can't access the cache created by the backup container.
	// Hence, we need to run all the containers of backup job with same user. Here, we are adding a default
	// securityContext at pod level to overcome the issue.
	// This securityContext will be overwritten by the user provided securityContext in the backup invokers.
	podSpec.SecurityContext = &core.PodSecurityContext{
		RunAsUser: pointer.Int64P(65535),
	}
	// apply RuntimeSettings to PodSpec
	if o.RuntimeSettings.Pod != nil {
		podSpec = ofst_util.ApplyPodRuntimeSettings(podSpec, *o.RuntimeSettings.Pod)
	}
	// always upsert tmp volume
	podSpec.Volumes = util.UpsertTmpVolume(podSpec.Volumes, o.TempDir)
	return podSpec, nil
}

func resolveWithInputs(obj interface{}, inputs map[string]string) error {
	// convert to JSON, apply replacements and convert back to struct
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	resolved, err := envsubst.EvalMap(string(jsonObj), inputs)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(resolved), obj)
}

func ResolveBackupBlueprint(bb *v1beta1_api.BackupBlueprint, input map[string]string) error {
	return resolveWithInputs(bb, input)
}

func ResolvePVCSpec(pvc *core.PersistentVolumeClaim, input map[string]string) error {
	return resolveWithInputs(pvc, input)
}

// GetPVCFromVolumeClaimTemplates returns list of PVCs generated according to the VolumeClaimTemplates
func GetPVCFromVolumeClaimTemplates(ordinal int32, claimTemplates []ofst.PersistentVolumeClaim) ([]core.PersistentVolumeClaim, error) {
	pvcList := make([]core.PersistentVolumeClaim, 0)
	for i := range claimTemplates {
		inputs := make(map[string]string)
		inputs[apis.KeyPodOrdinal] = strconv.Itoa(int(ordinal))
		claim := claimTemplates[i].DeepCopy().ToCorePVC()
		err := ResolvePVCSpec(claim, inputs)
		if err != nil {
			return pvcList, err
		}
		pvcList = append(pvcList, *claim)
	}
	return pvcList, nil
}
