// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by conversion-gen. DO NOT EDIT.

package v1beta1

import (
	unsafe "unsafe"

	v1beta1 "k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	apps "k8s.io/kubernetes/pkg/apis/apps"
	autoscaling "k8s.io/kubernetes/pkg/apis/autoscaling"
	core "k8s.io/kubernetes/pkg/apis/core"
	core_v1 "k8s.io/kubernetes/pkg/apis/core/v1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1beta1_ControllerRevision_To_apps_ControllerRevision,
		Convert_apps_ControllerRevision_To_v1beta1_ControllerRevision,
		Convert_v1beta1_ControllerRevisionList_To_apps_ControllerRevisionList,
		Convert_apps_ControllerRevisionList_To_v1beta1_ControllerRevisionList,
		Convert_v1beta1_Deployment_To_extensions_Deployment,
		Convert_extensions_Deployment_To_v1beta1_Deployment,
		Convert_v1beta1_DeploymentCondition_To_extensions_DeploymentCondition,
		Convert_extensions_DeploymentCondition_To_v1beta1_DeploymentCondition,
		Convert_v1beta1_DeploymentList_To_extensions_DeploymentList,
		Convert_extensions_DeploymentList_To_v1beta1_DeploymentList,
		Convert_v1beta1_DeploymentRollback_To_extensions_DeploymentRollback,
		Convert_extensions_DeploymentRollback_To_v1beta1_DeploymentRollback,
		Convert_v1beta1_DeploymentSpec_To_extensions_DeploymentSpec,
		Convert_extensions_DeploymentSpec_To_v1beta1_DeploymentSpec,
		Convert_v1beta1_DeploymentStatus_To_extensions_DeploymentStatus,
		Convert_extensions_DeploymentStatus_To_v1beta1_DeploymentStatus,
		Convert_v1beta1_DeploymentStrategy_To_extensions_DeploymentStrategy,
		Convert_extensions_DeploymentStrategy_To_v1beta1_DeploymentStrategy,
		Convert_v1beta1_RollbackConfig_To_extensions_RollbackConfig,
		Convert_extensions_RollbackConfig_To_v1beta1_RollbackConfig,
		Convert_v1beta1_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment,
		Convert_extensions_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment,
		Convert_v1beta1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy,
		Convert_apps_RollingUpdateStatefulSetStrategy_To_v1beta1_RollingUpdateStatefulSetStrategy,
		Convert_v1beta1_Scale_To_autoscaling_Scale,
		Convert_autoscaling_Scale_To_v1beta1_Scale,
		Convert_v1beta1_ScaleSpec_To_autoscaling_ScaleSpec,
		Convert_autoscaling_ScaleSpec_To_v1beta1_ScaleSpec,
		Convert_v1beta1_ScaleStatus_To_autoscaling_ScaleStatus,
		Convert_autoscaling_ScaleStatus_To_v1beta1_ScaleStatus,
		Convert_v1beta1_StatefulSet_To_apps_StatefulSet,
		Convert_apps_StatefulSet_To_v1beta1_StatefulSet,
		Convert_v1beta1_StatefulSetCondition_To_apps_StatefulSetCondition,
		Convert_apps_StatefulSetCondition_To_v1beta1_StatefulSetCondition,
		Convert_v1beta1_StatefulSetList_To_apps_StatefulSetList,
		Convert_apps_StatefulSetList_To_v1beta1_StatefulSetList,
		Convert_v1beta1_StatefulSetSpec_To_apps_StatefulSetSpec,
		Convert_apps_StatefulSetSpec_To_v1beta1_StatefulSetSpec,
		Convert_v1beta1_StatefulSetStatus_To_apps_StatefulSetStatus,
		Convert_apps_StatefulSetStatus_To_v1beta1_StatefulSetStatus,
		Convert_v1beta1_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy,
		Convert_apps_StatefulSetUpdateStrategy_To_v1beta1_StatefulSetUpdateStrategy,
	)
}

func autoConvert_v1beta1_ControllerRevision_To_apps_ControllerRevision(in *v1beta1.ControllerRevision, out *apps.ControllerRevision, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := runtime.Convert_runtime_RawExtension_To_runtime_Object(&in.Data, &out.Data, s); err != nil {
		return err
	}
	out.Revision = in.Revision
	return nil
}

// Convert_v1beta1_ControllerRevision_To_apps_ControllerRevision is an autogenerated conversion function.
func Convert_v1beta1_ControllerRevision_To_apps_ControllerRevision(in *v1beta1.ControllerRevision, out *apps.ControllerRevision, s conversion.Scope) error {
	return autoConvert_v1beta1_ControllerRevision_To_apps_ControllerRevision(in, out, s)
}

func autoConvert_apps_ControllerRevision_To_v1beta1_ControllerRevision(in *apps.ControllerRevision, out *v1beta1.ControllerRevision, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := runtime.Convert_runtime_Object_To_runtime_RawExtension(&in.Data, &out.Data, s); err != nil {
		return err
	}
	out.Revision = in.Revision
	return nil
}

// Convert_apps_ControllerRevision_To_v1beta1_ControllerRevision is an autogenerated conversion function.
func Convert_apps_ControllerRevision_To_v1beta1_ControllerRevision(in *apps.ControllerRevision, out *v1beta1.ControllerRevision, s conversion.Scope) error {
	return autoConvert_apps_ControllerRevision_To_v1beta1_ControllerRevision(in, out, s)
}

func autoConvert_v1beta1_ControllerRevisionList_To_apps_ControllerRevisionList(in *v1beta1.ControllerRevisionList, out *apps.ControllerRevisionList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]apps.ControllerRevision, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_ControllerRevision_To_apps_ControllerRevision(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta1_ControllerRevisionList_To_apps_ControllerRevisionList is an autogenerated conversion function.
func Convert_v1beta1_ControllerRevisionList_To_apps_ControllerRevisionList(in *v1beta1.ControllerRevisionList, out *apps.ControllerRevisionList, s conversion.Scope) error {
	return autoConvert_v1beta1_ControllerRevisionList_To_apps_ControllerRevisionList(in, out, s)
}

func autoConvert_apps_ControllerRevisionList_To_v1beta1_ControllerRevisionList(in *apps.ControllerRevisionList, out *v1beta1.ControllerRevisionList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1beta1.ControllerRevision, len(*in))
		for i := range *in {
			if err := Convert_apps_ControllerRevision_To_v1beta1_ControllerRevision(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_apps_ControllerRevisionList_To_v1beta1_ControllerRevisionList is an autogenerated conversion function.
func Convert_apps_ControllerRevisionList_To_v1beta1_ControllerRevisionList(in *apps.ControllerRevisionList, out *v1beta1.ControllerRevisionList, s conversion.Scope) error {
	return autoConvert_apps_ControllerRevisionList_To_v1beta1_ControllerRevisionList(in, out, s)
}

func autoConvert_v1beta1_Deployment_To_extensions_Deployment(in *v1beta1.Deployment, out *extensions.Deployment, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta1_DeploymentSpec_To_extensions_DeploymentSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_DeploymentStatus_To_extensions_DeploymentStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_Deployment_To_extensions_Deployment is an autogenerated conversion function.
func Convert_v1beta1_Deployment_To_extensions_Deployment(in *v1beta1.Deployment, out *extensions.Deployment, s conversion.Scope) error {
	return autoConvert_v1beta1_Deployment_To_extensions_Deployment(in, out, s)
}

func autoConvert_extensions_Deployment_To_v1beta1_Deployment(in *extensions.Deployment, out *v1beta1.Deployment, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_extensions_DeploymentSpec_To_v1beta1_DeploymentSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_extensions_DeploymentStatus_To_v1beta1_DeploymentStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_extensions_Deployment_To_v1beta1_Deployment is an autogenerated conversion function.
func Convert_extensions_Deployment_To_v1beta1_Deployment(in *extensions.Deployment, out *v1beta1.Deployment, s conversion.Scope) error {
	return autoConvert_extensions_Deployment_To_v1beta1_Deployment(in, out, s)
}

func autoConvert_v1beta1_DeploymentCondition_To_extensions_DeploymentCondition(in *v1beta1.DeploymentCondition, out *extensions.DeploymentCondition, s conversion.Scope) error {
	out.Type = extensions.DeploymentConditionType(in.Type)
	out.Status = core.ConditionStatus(in.Status)
	out.LastUpdateTime = in.LastUpdateTime
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

// Convert_v1beta1_DeploymentCondition_To_extensions_DeploymentCondition is an autogenerated conversion function.
func Convert_v1beta1_DeploymentCondition_To_extensions_DeploymentCondition(in *v1beta1.DeploymentCondition, out *extensions.DeploymentCondition, s conversion.Scope) error {
	return autoConvert_v1beta1_DeploymentCondition_To_extensions_DeploymentCondition(in, out, s)
}

func autoConvert_extensions_DeploymentCondition_To_v1beta1_DeploymentCondition(in *extensions.DeploymentCondition, out *v1beta1.DeploymentCondition, s conversion.Scope) error {
	out.Type = v1beta1.DeploymentConditionType(in.Type)
	out.Status = v1.ConditionStatus(in.Status)
	out.LastUpdateTime = in.LastUpdateTime
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

// Convert_extensions_DeploymentCondition_To_v1beta1_DeploymentCondition is an autogenerated conversion function.
func Convert_extensions_DeploymentCondition_To_v1beta1_DeploymentCondition(in *extensions.DeploymentCondition, out *v1beta1.DeploymentCondition, s conversion.Scope) error {
	return autoConvert_extensions_DeploymentCondition_To_v1beta1_DeploymentCondition(in, out, s)
}

func autoConvert_v1beta1_DeploymentList_To_extensions_DeploymentList(in *v1beta1.DeploymentList, out *extensions.DeploymentList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]extensions.Deployment, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_Deployment_To_extensions_Deployment(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta1_DeploymentList_To_extensions_DeploymentList is an autogenerated conversion function.
func Convert_v1beta1_DeploymentList_To_extensions_DeploymentList(in *v1beta1.DeploymentList, out *extensions.DeploymentList, s conversion.Scope) error {
	return autoConvert_v1beta1_DeploymentList_To_extensions_DeploymentList(in, out, s)
}

func autoConvert_extensions_DeploymentList_To_v1beta1_DeploymentList(in *extensions.DeploymentList, out *v1beta1.DeploymentList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1beta1.Deployment, len(*in))
		for i := range *in {
			if err := Convert_extensions_Deployment_To_v1beta1_Deployment(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_extensions_DeploymentList_To_v1beta1_DeploymentList is an autogenerated conversion function.
func Convert_extensions_DeploymentList_To_v1beta1_DeploymentList(in *extensions.DeploymentList, out *v1beta1.DeploymentList, s conversion.Scope) error {
	return autoConvert_extensions_DeploymentList_To_v1beta1_DeploymentList(in, out, s)
}

func autoConvert_v1beta1_DeploymentRollback_To_extensions_DeploymentRollback(in *v1beta1.DeploymentRollback, out *extensions.DeploymentRollback, s conversion.Scope) error {
	out.Name = in.Name
	out.UpdatedAnnotations = *(*map[string]string)(unsafe.Pointer(&in.UpdatedAnnotations))
	if err := Convert_v1beta1_RollbackConfig_To_extensions_RollbackConfig(&in.RollbackTo, &out.RollbackTo, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_DeploymentRollback_To_extensions_DeploymentRollback is an autogenerated conversion function.
func Convert_v1beta1_DeploymentRollback_To_extensions_DeploymentRollback(in *v1beta1.DeploymentRollback, out *extensions.DeploymentRollback, s conversion.Scope) error {
	return autoConvert_v1beta1_DeploymentRollback_To_extensions_DeploymentRollback(in, out, s)
}

func autoConvert_extensions_DeploymentRollback_To_v1beta1_DeploymentRollback(in *extensions.DeploymentRollback, out *v1beta1.DeploymentRollback, s conversion.Scope) error {
	out.Name = in.Name
	out.UpdatedAnnotations = *(*map[string]string)(unsafe.Pointer(&in.UpdatedAnnotations))
	if err := Convert_extensions_RollbackConfig_To_v1beta1_RollbackConfig(&in.RollbackTo, &out.RollbackTo, s); err != nil {
		return err
	}
	return nil
}

// Convert_extensions_DeploymentRollback_To_v1beta1_DeploymentRollback is an autogenerated conversion function.
func Convert_extensions_DeploymentRollback_To_v1beta1_DeploymentRollback(in *extensions.DeploymentRollback, out *v1beta1.DeploymentRollback, s conversion.Scope) error {
	return autoConvert_extensions_DeploymentRollback_To_v1beta1_DeploymentRollback(in, out, s)
}

func autoConvert_v1beta1_DeploymentSpec_To_extensions_DeploymentSpec(in *v1beta1.DeploymentSpec, out *extensions.DeploymentSpec, s conversion.Scope) error {
	if err := meta_v1.Convert_Pointer_int32_To_int32(&in.Replicas, &out.Replicas, s); err != nil {
		return err
	}
	out.Selector = (*meta_v1.LabelSelector)(unsafe.Pointer(in.Selector))
	if err := core_v1.Convert_v1_PodTemplateSpec_To_core_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_DeploymentStrategy_To_extensions_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	out.MinReadySeconds = in.MinReadySeconds
	out.RevisionHistoryLimit = (*int32)(unsafe.Pointer(in.RevisionHistoryLimit))
	out.Paused = in.Paused
	out.RollbackTo = (*extensions.RollbackConfig)(unsafe.Pointer(in.RollbackTo))
	out.ProgressDeadlineSeconds = (*int32)(unsafe.Pointer(in.ProgressDeadlineSeconds))
	return nil
}

func autoConvert_extensions_DeploymentSpec_To_v1beta1_DeploymentSpec(in *extensions.DeploymentSpec, out *v1beta1.DeploymentSpec, s conversion.Scope) error {
	if err := meta_v1.Convert_int32_To_Pointer_int32(&in.Replicas, &out.Replicas, s); err != nil {
		return err
	}
	out.Selector = (*meta_v1.LabelSelector)(unsafe.Pointer(in.Selector))
	if err := core_v1.Convert_core_PodTemplateSpec_To_v1_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_extensions_DeploymentStrategy_To_v1beta1_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	out.MinReadySeconds = in.MinReadySeconds
	out.RevisionHistoryLimit = (*int32)(unsafe.Pointer(in.RevisionHistoryLimit))
	out.Paused = in.Paused
	out.RollbackTo = (*v1beta1.RollbackConfig)(unsafe.Pointer(in.RollbackTo))
	out.ProgressDeadlineSeconds = (*int32)(unsafe.Pointer(in.ProgressDeadlineSeconds))
	return nil
}

func autoConvert_v1beta1_DeploymentStatus_To_extensions_DeploymentStatus(in *v1beta1.DeploymentStatus, out *extensions.DeploymentStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.Replicas = in.Replicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.ReadyReplicas = in.ReadyReplicas
	out.AvailableReplicas = in.AvailableReplicas
	out.UnavailableReplicas = in.UnavailableReplicas
	out.Conditions = *(*[]extensions.DeploymentCondition)(unsafe.Pointer(&in.Conditions))
	out.CollisionCount = (*int32)(unsafe.Pointer(in.CollisionCount))
	return nil
}

// Convert_v1beta1_DeploymentStatus_To_extensions_DeploymentStatus is an autogenerated conversion function.
func Convert_v1beta1_DeploymentStatus_To_extensions_DeploymentStatus(in *v1beta1.DeploymentStatus, out *extensions.DeploymentStatus, s conversion.Scope) error {
	return autoConvert_v1beta1_DeploymentStatus_To_extensions_DeploymentStatus(in, out, s)
}

func autoConvert_extensions_DeploymentStatus_To_v1beta1_DeploymentStatus(in *extensions.DeploymentStatus, out *v1beta1.DeploymentStatus, s conversion.Scope) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.Replicas = in.Replicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.ReadyReplicas = in.ReadyReplicas
	out.AvailableReplicas = in.AvailableReplicas
	out.UnavailableReplicas = in.UnavailableReplicas
	out.Conditions = *(*[]v1beta1.DeploymentCondition)(unsafe.Pointer(&in.Conditions))
	out.CollisionCount = (*int32)(unsafe.Pointer(in.CollisionCount))
	return nil
}

// Convert_extensions_DeploymentStatus_To_v1beta1_DeploymentStatus is an autogenerated conversion function.
func Convert_extensions_DeploymentStatus_To_v1beta1_DeploymentStatus(in *extensions.DeploymentStatus, out *v1beta1.DeploymentStatus, s conversion.Scope) error {
	return autoConvert_extensions_DeploymentStatus_To_v1beta1_DeploymentStatus(in, out, s)
}

func autoConvert_v1beta1_DeploymentStrategy_To_extensions_DeploymentStrategy(in *v1beta1.DeploymentStrategy, out *extensions.DeploymentStrategy, s conversion.Scope) error {
	out.Type = extensions.DeploymentStrategyType(in.Type)
	if in.RollingUpdate != nil {
		in, out := &in.RollingUpdate, &out.RollingUpdate
		*out = new(extensions.RollingUpdateDeployment)
		if err := Convert_v1beta1_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func autoConvert_extensions_DeploymentStrategy_To_v1beta1_DeploymentStrategy(in *extensions.DeploymentStrategy, out *v1beta1.DeploymentStrategy, s conversion.Scope) error {
	out.Type = v1beta1.DeploymentStrategyType(in.Type)
	if in.RollingUpdate != nil {
		in, out := &in.RollingUpdate, &out.RollingUpdate
		*out = new(v1beta1.RollingUpdateDeployment)
		if err := Convert_extensions_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func autoConvert_v1beta1_RollbackConfig_To_extensions_RollbackConfig(in *v1beta1.RollbackConfig, out *extensions.RollbackConfig, s conversion.Scope) error {
	out.Revision = in.Revision
	return nil
}

// Convert_v1beta1_RollbackConfig_To_extensions_RollbackConfig is an autogenerated conversion function.
func Convert_v1beta1_RollbackConfig_To_extensions_RollbackConfig(in *v1beta1.RollbackConfig, out *extensions.RollbackConfig, s conversion.Scope) error {
	return autoConvert_v1beta1_RollbackConfig_To_extensions_RollbackConfig(in, out, s)
}

func autoConvert_extensions_RollbackConfig_To_v1beta1_RollbackConfig(in *extensions.RollbackConfig, out *v1beta1.RollbackConfig, s conversion.Scope) error {
	out.Revision = in.Revision
	return nil
}

// Convert_extensions_RollbackConfig_To_v1beta1_RollbackConfig is an autogenerated conversion function.
func Convert_extensions_RollbackConfig_To_v1beta1_RollbackConfig(in *extensions.RollbackConfig, out *v1beta1.RollbackConfig, s conversion.Scope) error {
	return autoConvert_extensions_RollbackConfig_To_v1beta1_RollbackConfig(in, out, s)
}

func autoConvert_v1beta1_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment(in *v1beta1.RollingUpdateDeployment, out *extensions.RollingUpdateDeployment, s conversion.Scope) error {
	// WARNING: in.MaxUnavailable requires manual conversion: inconvertible types (*k8s.io/apimachinery/pkg/util/intstr.IntOrString vs k8s.io/apimachinery/pkg/util/intstr.IntOrString)
	// WARNING: in.MaxSurge requires manual conversion: inconvertible types (*k8s.io/apimachinery/pkg/util/intstr.IntOrString vs k8s.io/apimachinery/pkg/util/intstr.IntOrString)
	return nil
}

func autoConvert_extensions_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment(in *extensions.RollingUpdateDeployment, out *v1beta1.RollingUpdateDeployment, s conversion.Scope) error {
	// WARNING: in.MaxUnavailable requires manual conversion: inconvertible types (k8s.io/apimachinery/pkg/util/intstr.IntOrString vs *k8s.io/apimachinery/pkg/util/intstr.IntOrString)
	// WARNING: in.MaxSurge requires manual conversion: inconvertible types (k8s.io/apimachinery/pkg/util/intstr.IntOrString vs *k8s.io/apimachinery/pkg/util/intstr.IntOrString)
	return nil
}

func autoConvert_v1beta1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy(in *v1beta1.RollingUpdateStatefulSetStrategy, out *apps.RollingUpdateStatefulSetStrategy, s conversion.Scope) error {
	if err := meta_v1.Convert_Pointer_int32_To_int32(&in.Partition, &out.Partition, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy is an autogenerated conversion function.
func Convert_v1beta1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy(in *v1beta1.RollingUpdateStatefulSetStrategy, out *apps.RollingUpdateStatefulSetStrategy, s conversion.Scope) error {
	return autoConvert_v1beta1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy(in, out, s)
}

func autoConvert_apps_RollingUpdateStatefulSetStrategy_To_v1beta1_RollingUpdateStatefulSetStrategy(in *apps.RollingUpdateStatefulSetStrategy, out *v1beta1.RollingUpdateStatefulSetStrategy, s conversion.Scope) error {
	if err := meta_v1.Convert_int32_To_Pointer_int32(&in.Partition, &out.Partition, s); err != nil {
		return err
	}
	return nil
}

// Convert_apps_RollingUpdateStatefulSetStrategy_To_v1beta1_RollingUpdateStatefulSetStrategy is an autogenerated conversion function.
func Convert_apps_RollingUpdateStatefulSetStrategy_To_v1beta1_RollingUpdateStatefulSetStrategy(in *apps.RollingUpdateStatefulSetStrategy, out *v1beta1.RollingUpdateStatefulSetStrategy, s conversion.Scope) error {
	return autoConvert_apps_RollingUpdateStatefulSetStrategy_To_v1beta1_RollingUpdateStatefulSetStrategy(in, out, s)
}

func autoConvert_v1beta1_Scale_To_autoscaling_Scale(in *v1beta1.Scale, out *autoscaling.Scale, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta1_ScaleSpec_To_autoscaling_ScaleSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_ScaleStatus_To_autoscaling_ScaleStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_Scale_To_autoscaling_Scale is an autogenerated conversion function.
func Convert_v1beta1_Scale_To_autoscaling_Scale(in *v1beta1.Scale, out *autoscaling.Scale, s conversion.Scope) error {
	return autoConvert_v1beta1_Scale_To_autoscaling_Scale(in, out, s)
}

func autoConvert_autoscaling_Scale_To_v1beta1_Scale(in *autoscaling.Scale, out *v1beta1.Scale, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_autoscaling_ScaleSpec_To_v1beta1_ScaleSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_autoscaling_ScaleStatus_To_v1beta1_ScaleStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_autoscaling_Scale_To_v1beta1_Scale is an autogenerated conversion function.
func Convert_autoscaling_Scale_To_v1beta1_Scale(in *autoscaling.Scale, out *v1beta1.Scale, s conversion.Scope) error {
	return autoConvert_autoscaling_Scale_To_v1beta1_Scale(in, out, s)
}

func autoConvert_v1beta1_ScaleSpec_To_autoscaling_ScaleSpec(in *v1beta1.ScaleSpec, out *autoscaling.ScaleSpec, s conversion.Scope) error {
	out.Replicas = in.Replicas
	return nil
}

// Convert_v1beta1_ScaleSpec_To_autoscaling_ScaleSpec is an autogenerated conversion function.
func Convert_v1beta1_ScaleSpec_To_autoscaling_ScaleSpec(in *v1beta1.ScaleSpec, out *autoscaling.ScaleSpec, s conversion.Scope) error {
	return autoConvert_v1beta1_ScaleSpec_To_autoscaling_ScaleSpec(in, out, s)
}

func autoConvert_autoscaling_ScaleSpec_To_v1beta1_ScaleSpec(in *autoscaling.ScaleSpec, out *v1beta1.ScaleSpec, s conversion.Scope) error {
	out.Replicas = in.Replicas
	return nil
}

// Convert_autoscaling_ScaleSpec_To_v1beta1_ScaleSpec is an autogenerated conversion function.
func Convert_autoscaling_ScaleSpec_To_v1beta1_ScaleSpec(in *autoscaling.ScaleSpec, out *v1beta1.ScaleSpec, s conversion.Scope) error {
	return autoConvert_autoscaling_ScaleSpec_To_v1beta1_ScaleSpec(in, out, s)
}

func autoConvert_v1beta1_ScaleStatus_To_autoscaling_ScaleStatus(in *v1beta1.ScaleStatus, out *autoscaling.ScaleStatus, s conversion.Scope) error {
	out.Replicas = in.Replicas
	// WARNING: in.Selector requires manual conversion: inconvertible types (map[string]string vs string)
	// WARNING: in.TargetSelector requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_autoscaling_ScaleStatus_To_v1beta1_ScaleStatus(in *autoscaling.ScaleStatus, out *v1beta1.ScaleStatus, s conversion.Scope) error {
	out.Replicas = in.Replicas
	// WARNING: in.Selector requires manual conversion: inconvertible types (string vs map[string]string)
	return nil
}

func autoConvert_v1beta1_StatefulSet_To_apps_StatefulSet(in *v1beta1.StatefulSet, out *apps.StatefulSet, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta1_StatefulSetSpec_To_apps_StatefulSetSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_StatefulSetStatus_To_apps_StatefulSetStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_StatefulSet_To_apps_StatefulSet is an autogenerated conversion function.
func Convert_v1beta1_StatefulSet_To_apps_StatefulSet(in *v1beta1.StatefulSet, out *apps.StatefulSet, s conversion.Scope) error {
	return autoConvert_v1beta1_StatefulSet_To_apps_StatefulSet(in, out, s)
}

func autoConvert_apps_StatefulSet_To_v1beta1_StatefulSet(in *apps.StatefulSet, out *v1beta1.StatefulSet, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_apps_StatefulSetSpec_To_v1beta1_StatefulSetSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_apps_StatefulSetStatus_To_v1beta1_StatefulSetStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_apps_StatefulSet_To_v1beta1_StatefulSet is an autogenerated conversion function.
func Convert_apps_StatefulSet_To_v1beta1_StatefulSet(in *apps.StatefulSet, out *v1beta1.StatefulSet, s conversion.Scope) error {
	return autoConvert_apps_StatefulSet_To_v1beta1_StatefulSet(in, out, s)
}

func autoConvert_v1beta1_StatefulSetCondition_To_apps_StatefulSetCondition(in *v1beta1.StatefulSetCondition, out *apps.StatefulSetCondition, s conversion.Scope) error {
	out.Type = apps.StatefulSetConditionType(in.Type)
	out.Status = core.ConditionStatus(in.Status)
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

// Convert_v1beta1_StatefulSetCondition_To_apps_StatefulSetCondition is an autogenerated conversion function.
func Convert_v1beta1_StatefulSetCondition_To_apps_StatefulSetCondition(in *v1beta1.StatefulSetCondition, out *apps.StatefulSetCondition, s conversion.Scope) error {
	return autoConvert_v1beta1_StatefulSetCondition_To_apps_StatefulSetCondition(in, out, s)
}

func autoConvert_apps_StatefulSetCondition_To_v1beta1_StatefulSetCondition(in *apps.StatefulSetCondition, out *v1beta1.StatefulSetCondition, s conversion.Scope) error {
	out.Type = v1beta1.StatefulSetConditionType(in.Type)
	out.Status = v1.ConditionStatus(in.Status)
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

// Convert_apps_StatefulSetCondition_To_v1beta1_StatefulSetCondition is an autogenerated conversion function.
func Convert_apps_StatefulSetCondition_To_v1beta1_StatefulSetCondition(in *apps.StatefulSetCondition, out *v1beta1.StatefulSetCondition, s conversion.Scope) error {
	return autoConvert_apps_StatefulSetCondition_To_v1beta1_StatefulSetCondition(in, out, s)
}

func autoConvert_v1beta1_StatefulSetList_To_apps_StatefulSetList(in *v1beta1.StatefulSetList, out *apps.StatefulSetList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]apps.StatefulSet, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_StatefulSet_To_apps_StatefulSet(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta1_StatefulSetList_To_apps_StatefulSetList is an autogenerated conversion function.
func Convert_v1beta1_StatefulSetList_To_apps_StatefulSetList(in *v1beta1.StatefulSetList, out *apps.StatefulSetList, s conversion.Scope) error {
	return autoConvert_v1beta1_StatefulSetList_To_apps_StatefulSetList(in, out, s)
}

func autoConvert_apps_StatefulSetList_To_v1beta1_StatefulSetList(in *apps.StatefulSetList, out *v1beta1.StatefulSetList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1beta1.StatefulSet, len(*in))
		for i := range *in {
			if err := Convert_apps_StatefulSet_To_v1beta1_StatefulSet(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_apps_StatefulSetList_To_v1beta1_StatefulSetList is an autogenerated conversion function.
func Convert_apps_StatefulSetList_To_v1beta1_StatefulSetList(in *apps.StatefulSetList, out *v1beta1.StatefulSetList, s conversion.Scope) error {
	return autoConvert_apps_StatefulSetList_To_v1beta1_StatefulSetList(in, out, s)
}

func autoConvert_v1beta1_StatefulSetSpec_To_apps_StatefulSetSpec(in *v1beta1.StatefulSetSpec, out *apps.StatefulSetSpec, s conversion.Scope) error {
	if err := meta_v1.Convert_Pointer_int32_To_int32(&in.Replicas, &out.Replicas, s); err != nil {
		return err
	}
	out.Selector = (*meta_v1.LabelSelector)(unsafe.Pointer(in.Selector))
	if err := core_v1.Convert_v1_PodTemplateSpec_To_core_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	out.VolumeClaimTemplates = *(*[]core.PersistentVolumeClaim)(unsafe.Pointer(&in.VolumeClaimTemplates))
	out.ServiceName = in.ServiceName
	out.PodManagementPolicy = apps.PodManagementPolicyType(in.PodManagementPolicy)
	if err := Convert_v1beta1_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(&in.UpdateStrategy, &out.UpdateStrategy, s); err != nil {
		return err
	}
	out.RevisionHistoryLimit = (*int32)(unsafe.Pointer(in.RevisionHistoryLimit))
	return nil
}

func autoConvert_apps_StatefulSetSpec_To_v1beta1_StatefulSetSpec(in *apps.StatefulSetSpec, out *v1beta1.StatefulSetSpec, s conversion.Scope) error {
	if err := meta_v1.Convert_int32_To_Pointer_int32(&in.Replicas, &out.Replicas, s); err != nil {
		return err
	}
	out.Selector = (*meta_v1.LabelSelector)(unsafe.Pointer(in.Selector))
	if err := core_v1.Convert_core_PodTemplateSpec_To_v1_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	out.VolumeClaimTemplates = *(*[]v1.PersistentVolumeClaim)(unsafe.Pointer(&in.VolumeClaimTemplates))
	out.ServiceName = in.ServiceName
	out.PodManagementPolicy = v1beta1.PodManagementPolicyType(in.PodManagementPolicy)
	if err := Convert_apps_StatefulSetUpdateStrategy_To_v1beta1_StatefulSetUpdateStrategy(&in.UpdateStrategy, &out.UpdateStrategy, s); err != nil {
		return err
	}
	out.RevisionHistoryLimit = (*int32)(unsafe.Pointer(in.RevisionHistoryLimit))
	return nil
}

func autoConvert_v1beta1_StatefulSetStatus_To_apps_StatefulSetStatus(in *v1beta1.StatefulSetStatus, out *apps.StatefulSetStatus, s conversion.Scope) error {
	out.ObservedGeneration = (*int64)(unsafe.Pointer(in.ObservedGeneration))
	out.Replicas = in.Replicas
	out.ReadyReplicas = in.ReadyReplicas
	out.CurrentReplicas = in.CurrentReplicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.CurrentRevision = in.CurrentRevision
	out.UpdateRevision = in.UpdateRevision
	out.CollisionCount = (*int32)(unsafe.Pointer(in.CollisionCount))
	out.Conditions = *(*[]apps.StatefulSetCondition)(unsafe.Pointer(&in.Conditions))
	return nil
}

// Convert_v1beta1_StatefulSetStatus_To_apps_StatefulSetStatus is an autogenerated conversion function.
func Convert_v1beta1_StatefulSetStatus_To_apps_StatefulSetStatus(in *v1beta1.StatefulSetStatus, out *apps.StatefulSetStatus, s conversion.Scope) error {
	return autoConvert_v1beta1_StatefulSetStatus_To_apps_StatefulSetStatus(in, out, s)
}

func autoConvert_apps_StatefulSetStatus_To_v1beta1_StatefulSetStatus(in *apps.StatefulSetStatus, out *v1beta1.StatefulSetStatus, s conversion.Scope) error {
	out.ObservedGeneration = (*int64)(unsafe.Pointer(in.ObservedGeneration))
	out.Replicas = in.Replicas
	out.ReadyReplicas = in.ReadyReplicas
	out.CurrentReplicas = in.CurrentReplicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.CurrentRevision = in.CurrentRevision
	out.UpdateRevision = in.UpdateRevision
	out.CollisionCount = (*int32)(unsafe.Pointer(in.CollisionCount))
	out.Conditions = *(*[]v1beta1.StatefulSetCondition)(unsafe.Pointer(&in.Conditions))
	return nil
}

// Convert_apps_StatefulSetStatus_To_v1beta1_StatefulSetStatus is an autogenerated conversion function.
func Convert_apps_StatefulSetStatus_To_v1beta1_StatefulSetStatus(in *apps.StatefulSetStatus, out *v1beta1.StatefulSetStatus, s conversion.Scope) error {
	return autoConvert_apps_StatefulSetStatus_To_v1beta1_StatefulSetStatus(in, out, s)
}

func autoConvert_v1beta1_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(in *v1beta1.StatefulSetUpdateStrategy, out *apps.StatefulSetUpdateStrategy, s conversion.Scope) error {
	out.Type = apps.StatefulSetUpdateStrategyType(in.Type)
	if in.RollingUpdate != nil {
		in, out := &in.RollingUpdate, &out.RollingUpdate
		*out = new(apps.RollingUpdateStatefulSetStrategy)
		if err := Convert_v1beta1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func autoConvert_apps_StatefulSetUpdateStrategy_To_v1beta1_StatefulSetUpdateStrategy(in *apps.StatefulSetUpdateStrategy, out *v1beta1.StatefulSetUpdateStrategy, s conversion.Scope) error {
	out.Type = v1beta1.StatefulSetUpdateStrategyType(in.Type)
	if in.RollingUpdate != nil {
		in, out := &in.RollingUpdate, &out.RollingUpdate
		*out = new(v1beta1.RollingUpdateStatefulSetStrategy)
		if err := Convert_apps_RollingUpdateStatefulSetStrategy_To_v1beta1_RollingUpdateStatefulSetStrategy(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}
