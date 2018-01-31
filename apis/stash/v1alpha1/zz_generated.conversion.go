// +build !ignore_autogenerated

/*
Copyright 2018 The Stash Authors.

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

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	stash "github.com/appscode/stash/apis/stash"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	unsafe "unsafe"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_AzureSpec_To_stash_AzureSpec,
		Convert_stash_AzureSpec_To_v1alpha1_AzureSpec,
		Convert_v1alpha1_B2Spec_To_stash_B2Spec,
		Convert_stash_B2Spec_To_v1alpha1_B2Spec,
		Convert_v1alpha1_Backend_To_stash_Backend,
		Convert_stash_Backend_To_v1alpha1_Backend,
		Convert_v1alpha1_FileGroup_To_stash_FileGroup,
		Convert_stash_FileGroup_To_v1alpha1_FileGroup,
		Convert_v1alpha1_GCSSpec_To_stash_GCSSpec,
		Convert_stash_GCSSpec_To_v1alpha1_GCSSpec,
		Convert_v1alpha1_LocalSpec_To_stash_LocalSpec,
		Convert_stash_LocalSpec_To_v1alpha1_LocalSpec,
		Convert_v1alpha1_LocalTypedReference_To_stash_LocalTypedReference,
		Convert_stash_LocalTypedReference_To_v1alpha1_LocalTypedReference,
		Convert_v1alpha1_Recovery_To_stash_Recovery,
		Convert_stash_Recovery_To_v1alpha1_Recovery,
		Convert_v1alpha1_RecoveryList_To_stash_RecoveryList,
		Convert_stash_RecoveryList_To_v1alpha1_RecoveryList,
		Convert_v1alpha1_RecoverySpec_To_stash_RecoverySpec,
		Convert_stash_RecoverySpec_To_v1alpha1_RecoverySpec,
		Convert_v1alpha1_RecoveryStatus_To_stash_RecoveryStatus,
		Convert_stash_RecoveryStatus_To_v1alpha1_RecoveryStatus,
		Convert_v1alpha1_RestServerSpec_To_stash_RestServerSpec,
		Convert_stash_RestServerSpec_To_v1alpha1_RestServerSpec,
		Convert_v1alpha1_Restic_To_stash_Restic,
		Convert_stash_Restic_To_v1alpha1_Restic,
		Convert_v1alpha1_ResticList_To_stash_ResticList,
		Convert_stash_ResticList_To_v1alpha1_ResticList,
		Convert_v1alpha1_ResticSpec_To_stash_ResticSpec,
		Convert_stash_ResticSpec_To_v1alpha1_ResticSpec,
		Convert_v1alpha1_ResticStatus_To_stash_ResticStatus,
		Convert_stash_ResticStatus_To_v1alpha1_ResticStatus,
		Convert_v1alpha1_RestoreStats_To_stash_RestoreStats,
		Convert_stash_RestoreStats_To_v1alpha1_RestoreStats,
		Convert_v1alpha1_RetentionPolicy_To_stash_RetentionPolicy,
		Convert_stash_RetentionPolicy_To_v1alpha1_RetentionPolicy,
		Convert_v1alpha1_S3Spec_To_stash_S3Spec,
		Convert_stash_S3Spec_To_v1alpha1_S3Spec,
		Convert_v1alpha1_SwiftSpec_To_stash_SwiftSpec,
		Convert_stash_SwiftSpec_To_v1alpha1_SwiftSpec,
	)
}

func autoConvert_v1alpha1_AzureSpec_To_stash_AzureSpec(in *AzureSpec, out *stash.AzureSpec, s conversion.Scope) error {
	out.Container = in.Container
	out.Prefix = in.Prefix
	return nil
}

// Convert_v1alpha1_AzureSpec_To_stash_AzureSpec is an autogenerated conversion function.
func Convert_v1alpha1_AzureSpec_To_stash_AzureSpec(in *AzureSpec, out *stash.AzureSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_AzureSpec_To_stash_AzureSpec(in, out, s)
}

func autoConvert_stash_AzureSpec_To_v1alpha1_AzureSpec(in *stash.AzureSpec, out *AzureSpec, s conversion.Scope) error {
	out.Container = in.Container
	out.Prefix = in.Prefix
	return nil
}

// Convert_stash_AzureSpec_To_v1alpha1_AzureSpec is an autogenerated conversion function.
func Convert_stash_AzureSpec_To_v1alpha1_AzureSpec(in *stash.AzureSpec, out *AzureSpec, s conversion.Scope) error {
	return autoConvert_stash_AzureSpec_To_v1alpha1_AzureSpec(in, out, s)
}

func autoConvert_v1alpha1_B2Spec_To_stash_B2Spec(in *B2Spec, out *stash.B2Spec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.Prefix = in.Prefix
	return nil
}

// Convert_v1alpha1_B2Spec_To_stash_B2Spec is an autogenerated conversion function.
func Convert_v1alpha1_B2Spec_To_stash_B2Spec(in *B2Spec, out *stash.B2Spec, s conversion.Scope) error {
	return autoConvert_v1alpha1_B2Spec_To_stash_B2Spec(in, out, s)
}

func autoConvert_stash_B2Spec_To_v1alpha1_B2Spec(in *stash.B2Spec, out *B2Spec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.Prefix = in.Prefix
	return nil
}

// Convert_stash_B2Spec_To_v1alpha1_B2Spec is an autogenerated conversion function.
func Convert_stash_B2Spec_To_v1alpha1_B2Spec(in *stash.B2Spec, out *B2Spec, s conversion.Scope) error {
	return autoConvert_stash_B2Spec_To_v1alpha1_B2Spec(in, out, s)
}

func autoConvert_v1alpha1_Backend_To_stash_Backend(in *Backend, out *stash.Backend, s conversion.Scope) error {
	out.StorageSecretName = in.StorageSecretName
	out.Local = (*stash.LocalSpec)(unsafe.Pointer(in.Local))
	out.S3 = (*stash.S3Spec)(unsafe.Pointer(in.S3))
	out.GCS = (*stash.GCSSpec)(unsafe.Pointer(in.GCS))
	out.Azure = (*stash.AzureSpec)(unsafe.Pointer(in.Azure))
	out.Swift = (*stash.SwiftSpec)(unsafe.Pointer(in.Swift))
	out.B2 = (*stash.B2Spec)(unsafe.Pointer(in.B2))
	return nil
}

// Convert_v1alpha1_Backend_To_stash_Backend is an autogenerated conversion function.
func Convert_v1alpha1_Backend_To_stash_Backend(in *Backend, out *stash.Backend, s conversion.Scope) error {
	return autoConvert_v1alpha1_Backend_To_stash_Backend(in, out, s)
}

func autoConvert_stash_Backend_To_v1alpha1_Backend(in *stash.Backend, out *Backend, s conversion.Scope) error {
	out.StorageSecretName = in.StorageSecretName
	out.Local = (*LocalSpec)(unsafe.Pointer(in.Local))
	out.S3 = (*S3Spec)(unsafe.Pointer(in.S3))
	out.GCS = (*GCSSpec)(unsafe.Pointer(in.GCS))
	out.Azure = (*AzureSpec)(unsafe.Pointer(in.Azure))
	out.Swift = (*SwiftSpec)(unsafe.Pointer(in.Swift))
	out.B2 = (*B2Spec)(unsafe.Pointer(in.B2))
	return nil
}

// Convert_stash_Backend_To_v1alpha1_Backend is an autogenerated conversion function.
func Convert_stash_Backend_To_v1alpha1_Backend(in *stash.Backend, out *Backend, s conversion.Scope) error {
	return autoConvert_stash_Backend_To_v1alpha1_Backend(in, out, s)
}

func autoConvert_v1alpha1_FileGroup_To_stash_FileGroup(in *FileGroup, out *stash.FileGroup, s conversion.Scope) error {
	out.Path = in.Path
	out.Tags = *(*[]string)(unsafe.Pointer(&in.Tags))
	out.RetentionPolicyName = in.RetentionPolicyName
	return nil
}

// Convert_v1alpha1_FileGroup_To_stash_FileGroup is an autogenerated conversion function.
func Convert_v1alpha1_FileGroup_To_stash_FileGroup(in *FileGroup, out *stash.FileGroup, s conversion.Scope) error {
	return autoConvert_v1alpha1_FileGroup_To_stash_FileGroup(in, out, s)
}

func autoConvert_stash_FileGroup_To_v1alpha1_FileGroup(in *stash.FileGroup, out *FileGroup, s conversion.Scope) error {
	out.Path = in.Path
	out.Tags = *(*[]string)(unsafe.Pointer(&in.Tags))
	out.RetentionPolicyName = in.RetentionPolicyName
	return nil
}

// Convert_stash_FileGroup_To_v1alpha1_FileGroup is an autogenerated conversion function.
func Convert_stash_FileGroup_To_v1alpha1_FileGroup(in *stash.FileGroup, out *FileGroup, s conversion.Scope) error {
	return autoConvert_stash_FileGroup_To_v1alpha1_FileGroup(in, out, s)
}

func autoConvert_v1alpha1_GCSSpec_To_stash_GCSSpec(in *GCSSpec, out *stash.GCSSpec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.Prefix = in.Prefix
	return nil
}

// Convert_v1alpha1_GCSSpec_To_stash_GCSSpec is an autogenerated conversion function.
func Convert_v1alpha1_GCSSpec_To_stash_GCSSpec(in *GCSSpec, out *stash.GCSSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_GCSSpec_To_stash_GCSSpec(in, out, s)
}

func autoConvert_stash_GCSSpec_To_v1alpha1_GCSSpec(in *stash.GCSSpec, out *GCSSpec, s conversion.Scope) error {
	out.Bucket = in.Bucket
	out.Prefix = in.Prefix
	return nil
}

// Convert_stash_GCSSpec_To_v1alpha1_GCSSpec is an autogenerated conversion function.
func Convert_stash_GCSSpec_To_v1alpha1_GCSSpec(in *stash.GCSSpec, out *GCSSpec, s conversion.Scope) error {
	return autoConvert_stash_GCSSpec_To_v1alpha1_GCSSpec(in, out, s)
}

func autoConvert_v1alpha1_LocalSpec_To_stash_LocalSpec(in *LocalSpec, out *stash.LocalSpec, s conversion.Scope) error {
	out.VolumeSource = in.VolumeSource
	out.MountPath = in.MountPath
	out.SubPath = in.SubPath
	return nil
}

// Convert_v1alpha1_LocalSpec_To_stash_LocalSpec is an autogenerated conversion function.
func Convert_v1alpha1_LocalSpec_To_stash_LocalSpec(in *LocalSpec, out *stash.LocalSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_LocalSpec_To_stash_LocalSpec(in, out, s)
}

func autoConvert_stash_LocalSpec_To_v1alpha1_LocalSpec(in *stash.LocalSpec, out *LocalSpec, s conversion.Scope) error {
	out.VolumeSource = in.VolumeSource
	out.MountPath = in.MountPath
	out.SubPath = in.SubPath
	return nil
}

// Convert_stash_LocalSpec_To_v1alpha1_LocalSpec is an autogenerated conversion function.
func Convert_stash_LocalSpec_To_v1alpha1_LocalSpec(in *stash.LocalSpec, out *LocalSpec, s conversion.Scope) error {
	return autoConvert_stash_LocalSpec_To_v1alpha1_LocalSpec(in, out, s)
}

func autoConvert_v1alpha1_LocalTypedReference_To_stash_LocalTypedReference(in *LocalTypedReference, out *stash.LocalTypedReference, s conversion.Scope) error {
	out.Kind = in.Kind
	out.Name = in.Name
	out.APIVersion = in.APIVersion
	return nil
}

// Convert_v1alpha1_LocalTypedReference_To_stash_LocalTypedReference is an autogenerated conversion function.
func Convert_v1alpha1_LocalTypedReference_To_stash_LocalTypedReference(in *LocalTypedReference, out *stash.LocalTypedReference, s conversion.Scope) error {
	return autoConvert_v1alpha1_LocalTypedReference_To_stash_LocalTypedReference(in, out, s)
}

func autoConvert_stash_LocalTypedReference_To_v1alpha1_LocalTypedReference(in *stash.LocalTypedReference, out *LocalTypedReference, s conversion.Scope) error {
	out.Kind = in.Kind
	out.Name = in.Name
	out.APIVersion = in.APIVersion
	return nil
}

// Convert_stash_LocalTypedReference_To_v1alpha1_LocalTypedReference is an autogenerated conversion function.
func Convert_stash_LocalTypedReference_To_v1alpha1_LocalTypedReference(in *stash.LocalTypedReference, out *LocalTypedReference, s conversion.Scope) error {
	return autoConvert_stash_LocalTypedReference_To_v1alpha1_LocalTypedReference(in, out, s)
}

func autoConvert_v1alpha1_Recovery_To_stash_Recovery(in *Recovery, out *stash.Recovery, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_RecoverySpec_To_stash_RecoverySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_RecoveryStatus_To_stash_RecoveryStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Recovery_To_stash_Recovery is an autogenerated conversion function.
func Convert_v1alpha1_Recovery_To_stash_Recovery(in *Recovery, out *stash.Recovery, s conversion.Scope) error {
	return autoConvert_v1alpha1_Recovery_To_stash_Recovery(in, out, s)
}

func autoConvert_stash_Recovery_To_v1alpha1_Recovery(in *stash.Recovery, out *Recovery, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_stash_RecoverySpec_To_v1alpha1_RecoverySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_stash_RecoveryStatus_To_v1alpha1_RecoveryStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_stash_Recovery_To_v1alpha1_Recovery is an autogenerated conversion function.
func Convert_stash_Recovery_To_v1alpha1_Recovery(in *stash.Recovery, out *Recovery, s conversion.Scope) error {
	return autoConvert_stash_Recovery_To_v1alpha1_Recovery(in, out, s)
}

func autoConvert_v1alpha1_RecoveryList_To_stash_RecoveryList(in *RecoveryList, out *stash.RecoveryList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]stash.Recovery)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_RecoveryList_To_stash_RecoveryList is an autogenerated conversion function.
func Convert_v1alpha1_RecoveryList_To_stash_RecoveryList(in *RecoveryList, out *stash.RecoveryList, s conversion.Scope) error {
	return autoConvert_v1alpha1_RecoveryList_To_stash_RecoveryList(in, out, s)
}

func autoConvert_stash_RecoveryList_To_v1alpha1_RecoveryList(in *stash.RecoveryList, out *RecoveryList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Recovery)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_stash_RecoveryList_To_v1alpha1_RecoveryList is an autogenerated conversion function.
func Convert_stash_RecoveryList_To_v1alpha1_RecoveryList(in *stash.RecoveryList, out *RecoveryList, s conversion.Scope) error {
	return autoConvert_stash_RecoveryList_To_v1alpha1_RecoveryList(in, out, s)
}

func autoConvert_v1alpha1_RecoverySpec_To_stash_RecoverySpec(in *RecoverySpec, out *stash.RecoverySpec, s conversion.Scope) error {
	if err := Convert_v1alpha1_Backend_To_stash_Backend(&in.Backend, &out.Backend, s); err != nil {
		return err
	}
	out.Paths = *(*[]string)(unsafe.Pointer(&in.Paths))
	if err := Convert_v1alpha1_LocalTypedReference_To_stash_LocalTypedReference(&in.Workload, &out.Workload, s); err != nil {
		return err
	}
	out.PodOrdinal = in.PodOrdinal
	out.NodeName = in.NodeName
	out.RecoveredVolumes = *(*[]stash.LocalSpec)(unsafe.Pointer(&in.RecoveredVolumes))
	out.ImagePullSecrets = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	return nil
}

// Convert_v1alpha1_RecoverySpec_To_stash_RecoverySpec is an autogenerated conversion function.
func Convert_v1alpha1_RecoverySpec_To_stash_RecoverySpec(in *RecoverySpec, out *stash.RecoverySpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_RecoverySpec_To_stash_RecoverySpec(in, out, s)
}

func autoConvert_stash_RecoverySpec_To_v1alpha1_RecoverySpec(in *stash.RecoverySpec, out *RecoverySpec, s conversion.Scope) error {
	if err := Convert_stash_Backend_To_v1alpha1_Backend(&in.Backend, &out.Backend, s); err != nil {
		return err
	}
	out.Paths = *(*[]string)(unsafe.Pointer(&in.Paths))
	if err := Convert_stash_LocalTypedReference_To_v1alpha1_LocalTypedReference(&in.Workload, &out.Workload, s); err != nil {
		return err
	}
	out.PodOrdinal = in.PodOrdinal
	out.NodeName = in.NodeName
	out.RecoveredVolumes = *(*[]LocalSpec)(unsafe.Pointer(&in.RecoveredVolumes))
	out.ImagePullSecrets = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	return nil
}

// Convert_stash_RecoverySpec_To_v1alpha1_RecoverySpec is an autogenerated conversion function.
func Convert_stash_RecoverySpec_To_v1alpha1_RecoverySpec(in *stash.RecoverySpec, out *RecoverySpec, s conversion.Scope) error {
	return autoConvert_stash_RecoverySpec_To_v1alpha1_RecoverySpec(in, out, s)
}

func autoConvert_v1alpha1_RecoveryStatus_To_stash_RecoveryStatus(in *RecoveryStatus, out *stash.RecoveryStatus, s conversion.Scope) error {
	out.Phase = stash.RecoveryPhase(in.Phase)
	out.Stats = *(*[]stash.RestoreStats)(unsafe.Pointer(&in.Stats))
	return nil
}

// Convert_v1alpha1_RecoveryStatus_To_stash_RecoveryStatus is an autogenerated conversion function.
func Convert_v1alpha1_RecoveryStatus_To_stash_RecoveryStatus(in *RecoveryStatus, out *stash.RecoveryStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_RecoveryStatus_To_stash_RecoveryStatus(in, out, s)
}

func autoConvert_stash_RecoveryStatus_To_v1alpha1_RecoveryStatus(in *stash.RecoveryStatus, out *RecoveryStatus, s conversion.Scope) error {
	out.Phase = RecoveryPhase(in.Phase)
	out.Stats = *(*[]RestoreStats)(unsafe.Pointer(&in.Stats))
	return nil
}

// Convert_stash_RecoveryStatus_To_v1alpha1_RecoveryStatus is an autogenerated conversion function.
func Convert_stash_RecoveryStatus_To_v1alpha1_RecoveryStatus(in *stash.RecoveryStatus, out *RecoveryStatus, s conversion.Scope) error {
	return autoConvert_stash_RecoveryStatus_To_v1alpha1_RecoveryStatus(in, out, s)
}

func autoConvert_v1alpha1_RestServerSpec_To_stash_RestServerSpec(in *RestServerSpec, out *stash.RestServerSpec, s conversion.Scope) error {
	out.URL = in.URL
	return nil
}

// Convert_v1alpha1_RestServerSpec_To_stash_RestServerSpec is an autogenerated conversion function.
func Convert_v1alpha1_RestServerSpec_To_stash_RestServerSpec(in *RestServerSpec, out *stash.RestServerSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_RestServerSpec_To_stash_RestServerSpec(in, out, s)
}

func autoConvert_stash_RestServerSpec_To_v1alpha1_RestServerSpec(in *stash.RestServerSpec, out *RestServerSpec, s conversion.Scope) error {
	out.URL = in.URL
	return nil
}

// Convert_stash_RestServerSpec_To_v1alpha1_RestServerSpec is an autogenerated conversion function.
func Convert_stash_RestServerSpec_To_v1alpha1_RestServerSpec(in *stash.RestServerSpec, out *RestServerSpec, s conversion.Scope) error {
	return autoConvert_stash_RestServerSpec_To_v1alpha1_RestServerSpec(in, out, s)
}

func autoConvert_v1alpha1_Restic_To_stash_Restic(in *Restic, out *stash.Restic, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_ResticSpec_To_stash_ResticSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ResticStatus_To_stash_ResticStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Restic_To_stash_Restic is an autogenerated conversion function.
func Convert_v1alpha1_Restic_To_stash_Restic(in *Restic, out *stash.Restic, s conversion.Scope) error {
	return autoConvert_v1alpha1_Restic_To_stash_Restic(in, out, s)
}

func autoConvert_stash_Restic_To_v1alpha1_Restic(in *stash.Restic, out *Restic, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_stash_ResticSpec_To_v1alpha1_ResticSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_stash_ResticStatus_To_v1alpha1_ResticStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_stash_Restic_To_v1alpha1_Restic is an autogenerated conversion function.
func Convert_stash_Restic_To_v1alpha1_Restic(in *stash.Restic, out *Restic, s conversion.Scope) error {
	return autoConvert_stash_Restic_To_v1alpha1_Restic(in, out, s)
}

func autoConvert_v1alpha1_ResticList_To_stash_ResticList(in *ResticList, out *stash.ResticList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]stash.Restic)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_ResticList_To_stash_ResticList is an autogenerated conversion function.
func Convert_v1alpha1_ResticList_To_stash_ResticList(in *ResticList, out *stash.ResticList, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResticList_To_stash_ResticList(in, out, s)
}

func autoConvert_stash_ResticList_To_v1alpha1_ResticList(in *stash.ResticList, out *ResticList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Restic)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_stash_ResticList_To_v1alpha1_ResticList is an autogenerated conversion function.
func Convert_stash_ResticList_To_v1alpha1_ResticList(in *stash.ResticList, out *ResticList, s conversion.Scope) error {
	return autoConvert_stash_ResticList_To_v1alpha1_ResticList(in, out, s)
}

func autoConvert_v1alpha1_ResticSpec_To_stash_ResticSpec(in *ResticSpec, out *stash.ResticSpec, s conversion.Scope) error {
	out.Selector = in.Selector
	out.FileGroups = *(*[]stash.FileGroup)(unsafe.Pointer(&in.FileGroups))
	if err := Convert_v1alpha1_Backend_To_stash_Backend(&in.Backend, &out.Backend, s); err != nil {
		return err
	}
	out.Schedule = in.Schedule
	out.VolumeMounts = *(*[]v1.VolumeMount)(unsafe.Pointer(&in.VolumeMounts))
	out.Resources = in.Resources
	out.RetentionPolicies = *(*[]stash.RetentionPolicy)(unsafe.Pointer(&in.RetentionPolicies))
	out.Type = stash.BackupType(in.Type)
	out.Paused = in.Paused
	out.ImagePullSecrets = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	return nil
}

// Convert_v1alpha1_ResticSpec_To_stash_ResticSpec is an autogenerated conversion function.
func Convert_v1alpha1_ResticSpec_To_stash_ResticSpec(in *ResticSpec, out *stash.ResticSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResticSpec_To_stash_ResticSpec(in, out, s)
}

func autoConvert_stash_ResticSpec_To_v1alpha1_ResticSpec(in *stash.ResticSpec, out *ResticSpec, s conversion.Scope) error {
	out.Selector = in.Selector
	out.FileGroups = *(*[]FileGroup)(unsafe.Pointer(&in.FileGroups))
	if err := Convert_stash_Backend_To_v1alpha1_Backend(&in.Backend, &out.Backend, s); err != nil {
		return err
	}
	out.Schedule = in.Schedule
	out.VolumeMounts = *(*[]v1.VolumeMount)(unsafe.Pointer(&in.VolumeMounts))
	out.Resources = in.Resources
	out.RetentionPolicies = *(*[]RetentionPolicy)(unsafe.Pointer(&in.RetentionPolicies))
	out.Type = BackupType(in.Type)
	out.Paused = in.Paused
	out.ImagePullSecrets = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	return nil
}

// Convert_stash_ResticSpec_To_v1alpha1_ResticSpec is an autogenerated conversion function.
func Convert_stash_ResticSpec_To_v1alpha1_ResticSpec(in *stash.ResticSpec, out *ResticSpec, s conversion.Scope) error {
	return autoConvert_stash_ResticSpec_To_v1alpha1_ResticSpec(in, out, s)
}

func autoConvert_v1alpha1_ResticStatus_To_stash_ResticStatus(in *ResticStatus, out *stash.ResticStatus, s conversion.Scope) error {
	out.FirstBackupTime = (*meta_v1.Time)(unsafe.Pointer(in.FirstBackupTime))
	out.LastBackupTime = (*meta_v1.Time)(unsafe.Pointer(in.LastBackupTime))
	out.LastSuccessfulBackupTime = (*meta_v1.Time)(unsafe.Pointer(in.LastSuccessfulBackupTime))
	out.LastBackupDuration = in.LastBackupDuration
	out.BackupCount = in.BackupCount
	return nil
}

// Convert_v1alpha1_ResticStatus_To_stash_ResticStatus is an autogenerated conversion function.
func Convert_v1alpha1_ResticStatus_To_stash_ResticStatus(in *ResticStatus, out *stash.ResticStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResticStatus_To_stash_ResticStatus(in, out, s)
}

func autoConvert_stash_ResticStatus_To_v1alpha1_ResticStatus(in *stash.ResticStatus, out *ResticStatus, s conversion.Scope) error {
	out.FirstBackupTime = (*meta_v1.Time)(unsafe.Pointer(in.FirstBackupTime))
	out.LastBackupTime = (*meta_v1.Time)(unsafe.Pointer(in.LastBackupTime))
	out.LastSuccessfulBackupTime = (*meta_v1.Time)(unsafe.Pointer(in.LastSuccessfulBackupTime))
	out.LastBackupDuration = in.LastBackupDuration
	out.BackupCount = in.BackupCount
	return nil
}

// Convert_stash_ResticStatus_To_v1alpha1_ResticStatus is an autogenerated conversion function.
func Convert_stash_ResticStatus_To_v1alpha1_ResticStatus(in *stash.ResticStatus, out *ResticStatus, s conversion.Scope) error {
	return autoConvert_stash_ResticStatus_To_v1alpha1_ResticStatus(in, out, s)
}

func autoConvert_v1alpha1_RestoreStats_To_stash_RestoreStats(in *RestoreStats, out *stash.RestoreStats, s conversion.Scope) error {
	out.Path = in.Path
	out.Phase = stash.RecoveryPhase(in.Phase)
	out.Duration = in.Duration
	return nil
}

// Convert_v1alpha1_RestoreStats_To_stash_RestoreStats is an autogenerated conversion function.
func Convert_v1alpha1_RestoreStats_To_stash_RestoreStats(in *RestoreStats, out *stash.RestoreStats, s conversion.Scope) error {
	return autoConvert_v1alpha1_RestoreStats_To_stash_RestoreStats(in, out, s)
}

func autoConvert_stash_RestoreStats_To_v1alpha1_RestoreStats(in *stash.RestoreStats, out *RestoreStats, s conversion.Scope) error {
	out.Path = in.Path
	out.Phase = RecoveryPhase(in.Phase)
	out.Duration = in.Duration
	return nil
}

// Convert_stash_RestoreStats_To_v1alpha1_RestoreStats is an autogenerated conversion function.
func Convert_stash_RestoreStats_To_v1alpha1_RestoreStats(in *stash.RestoreStats, out *RestoreStats, s conversion.Scope) error {
	return autoConvert_stash_RestoreStats_To_v1alpha1_RestoreStats(in, out, s)
}

func autoConvert_v1alpha1_RetentionPolicy_To_stash_RetentionPolicy(in *RetentionPolicy, out *stash.RetentionPolicy, s conversion.Scope) error {
	out.Name = in.Name
	out.KeepLast = in.KeepLast
	out.KeepHourly = in.KeepHourly
	out.KeepDaily = in.KeepDaily
	out.KeepWeekly = in.KeepWeekly
	out.KeepMonthly = in.KeepMonthly
	out.KeepYearly = in.KeepYearly
	out.KeepTags = *(*[]string)(unsafe.Pointer(&in.KeepTags))
	out.Prune = in.Prune
	out.DryRun = in.DryRun
	return nil
}

// Convert_v1alpha1_RetentionPolicy_To_stash_RetentionPolicy is an autogenerated conversion function.
func Convert_v1alpha1_RetentionPolicy_To_stash_RetentionPolicy(in *RetentionPolicy, out *stash.RetentionPolicy, s conversion.Scope) error {
	return autoConvert_v1alpha1_RetentionPolicy_To_stash_RetentionPolicy(in, out, s)
}

func autoConvert_stash_RetentionPolicy_To_v1alpha1_RetentionPolicy(in *stash.RetentionPolicy, out *RetentionPolicy, s conversion.Scope) error {
	out.Name = in.Name
	out.KeepLast = in.KeepLast
	out.KeepHourly = in.KeepHourly
	out.KeepDaily = in.KeepDaily
	out.KeepWeekly = in.KeepWeekly
	out.KeepMonthly = in.KeepMonthly
	out.KeepYearly = in.KeepYearly
	out.KeepTags = *(*[]string)(unsafe.Pointer(&in.KeepTags))
	out.Prune = in.Prune
	out.DryRun = in.DryRun
	return nil
}

// Convert_stash_RetentionPolicy_To_v1alpha1_RetentionPolicy is an autogenerated conversion function.
func Convert_stash_RetentionPolicy_To_v1alpha1_RetentionPolicy(in *stash.RetentionPolicy, out *RetentionPolicy, s conversion.Scope) error {
	return autoConvert_stash_RetentionPolicy_To_v1alpha1_RetentionPolicy(in, out, s)
}

func autoConvert_v1alpha1_S3Spec_To_stash_S3Spec(in *S3Spec, out *stash.S3Spec, s conversion.Scope) error {
	out.Endpoint = in.Endpoint
	out.Bucket = in.Bucket
	out.Prefix = in.Prefix
	return nil
}

// Convert_v1alpha1_S3Spec_To_stash_S3Spec is an autogenerated conversion function.
func Convert_v1alpha1_S3Spec_To_stash_S3Spec(in *S3Spec, out *stash.S3Spec, s conversion.Scope) error {
	return autoConvert_v1alpha1_S3Spec_To_stash_S3Spec(in, out, s)
}

func autoConvert_stash_S3Spec_To_v1alpha1_S3Spec(in *stash.S3Spec, out *S3Spec, s conversion.Scope) error {
	out.Endpoint = in.Endpoint
	out.Bucket = in.Bucket
	out.Prefix = in.Prefix
	return nil
}

// Convert_stash_S3Spec_To_v1alpha1_S3Spec is an autogenerated conversion function.
func Convert_stash_S3Spec_To_v1alpha1_S3Spec(in *stash.S3Spec, out *S3Spec, s conversion.Scope) error {
	return autoConvert_stash_S3Spec_To_v1alpha1_S3Spec(in, out, s)
}

func autoConvert_v1alpha1_SwiftSpec_To_stash_SwiftSpec(in *SwiftSpec, out *stash.SwiftSpec, s conversion.Scope) error {
	out.Container = in.Container
	out.Prefix = in.Prefix
	return nil
}

// Convert_v1alpha1_SwiftSpec_To_stash_SwiftSpec is an autogenerated conversion function.
func Convert_v1alpha1_SwiftSpec_To_stash_SwiftSpec(in *SwiftSpec, out *stash.SwiftSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_SwiftSpec_To_stash_SwiftSpec(in, out, s)
}

func autoConvert_stash_SwiftSpec_To_v1alpha1_SwiftSpec(in *stash.SwiftSpec, out *SwiftSpec, s conversion.Scope) error {
	out.Container = in.Container
	out.Prefix = in.Prefix
	return nil
}

// Convert_stash_SwiftSpec_To_v1alpha1_SwiftSpec is an autogenerated conversion function.
func Convert_stash_SwiftSpec_To_v1alpha1_SwiftSpec(in *stash.SwiftSpec, out *SwiftSpec, s conversion.Scope) error {
	return autoConvert_stash_SwiftSpec_To_v1alpha1_SwiftSpec(in, out, s)
}
