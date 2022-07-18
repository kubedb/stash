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

package rbac

import (
	"context"
	"strings"

	"stash.appscode.dev/apimachinery/apis"
	api_v1alpha1 "stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	api_v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	crdv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1beta1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	storage_api_v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/meta"
	meta_util "kmodules.xyz/client-go/meta"
	rbac_util "kmodules.xyz/client-go/rbac/v1"
)

func (opt *RBACOptions) EnsureVolumeSnapshotterJobRBAC() error {
	if opt.ServiceAccount.Name == "" {
		opt.ServiceAccount.Name = meta.ValidNameWithPrefixNSuffix(strings.ToLower(opt.Invoker.Kind), opt.Invoker.Name, opt.Suffix)
		err := opt.ensureServiceAccount()
		if err != nil {
			return err
		}
	}
	// ensure ClusterRole for VolumeSnapshot job
	err := opt.ensureVolumeSnapshotterJobClusterRole()
	if err != nil {
		return err
	}

	// ensure RoleBinding for VolumeSnapshot job
	err = opt.ensureVolumeSnapshotterJobRoleBinding()
	if err != nil {
		return err
	}

	return nil
}

func (opt *RBACOptions) ensureVolumeSnapshotterJobClusterRole() error {
	meta := metav1.ObjectMeta{
		Name:   apis.StashVolumeSnapshotterClusterRole,
		Labels: opt.OffshootLabels,
	}
	_, _, err := rbac_util.CreateOrPatchClusterRole(context.TODO(), opt.KubeClient, meta, func(in *rbac.ClusterRole) *rbac.ClusterRole {
		in.Rules = []rbac.PolicyRule{
			{
				APIGroups: []string{api_v1beta1.SchemeGroupVersion.Group},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{api_v1alpha1.SchemeGroupVersion.Group},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"events"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"pods"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"pods/exec"},
				Verbs:     []string{"get", "create"},
			},
			{
				APIGroups: []string{apps.GroupName},
				Resources: []string{"deployments", "statefulsets"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{apps.GroupName},
				Resources: []string{"daemonsets", "replicasets"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"replicationcontrollers"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{crdv1.GroupName},
				Resources: []string{"volumesnapshots", "volumesnapshotcontents", "volumesnapshotclasses"},
				Verbs:     []string{"create", "get", "list", "watch", "patch", "delete"},
			},
		}
		return in
	}, metav1.PatchOptions{})
	return err
}

func (opt *RBACOptions) ensureVolumeSnapshotterJobRoleBinding() error {
	meta := metav1.ObjectMeta{
		Name:      opt.getRoleBindingName(),
		Namespace: opt.Invoker.Namespace,
		Labels:    opt.OffshootLabels,
	}
	_, _, err := rbac_util.CreateOrPatchRoleBinding(context.TODO(), opt.KubeClient, meta, func(in *rbac.RoleBinding) *rbac.RoleBinding {
		core_util.EnsureOwnerReference(&in.ObjectMeta, opt.Owner)

		in.RoleRef = rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     apis.KindClusterRole,
			Name:     apis.StashVolumeSnapshotterClusterRole,
		}
		in.Subjects = []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind,
				Name:      opt.ServiceAccount.Name,
				Namespace: opt.ServiceAccount.Namespace,
			},
		}
		return in
	}, metav1.PatchOptions{})
	return err
}

func (opt *RBACOptions) EnsureVolumeSnapshotRestorerJobRBAC() error {
	if opt.ServiceAccount.Name == "" {
		opt.ServiceAccount.Name = meta.ValidNameWithPrefixNSuffix(strings.ToLower(opt.Invoker.Kind), opt.Invoker.Name, opt.Suffix)
		err := opt.ensureServiceAccount()
		if err != nil {
			return err
		}
	}
	// ensure ClusterRole for restore job
	err := opt.ensureVolumeSnapshotRestorerJobClusterRole()
	if err != nil {
		return err
	}

	// ensure RoleBinding for restore job
	err = opt.ensureVolumeSnapshotRestorerJobRoleBinding()
	if err != nil {
		return err
	}

	// ensure storageClass ClusterRole for restore job
	err = opt.ensureStorageReaderClassClusterRole()
	if err != nil {
		return err
	}

	// ensure storageClass ClusterRoleBinding for restore job
	err = opt.ensureStorageClassReaderClusterRoleBinding()
	if err != nil {
		return err
	}

	return nil
}

func (opt *RBACOptions) ensureVolumeSnapshotRestorerJobClusterRole() error {
	meta := metav1.ObjectMeta{
		Name:   apis.StashVolumeSnapshotRestorerClusterRole,
		Labels: opt.OffshootLabels,
	}
	_, _, err := rbac_util.CreateOrPatchClusterRole(context.TODO(), opt.KubeClient, meta, func(in *rbac.ClusterRole) *rbac.ClusterRole {
		in.Rules = []rbac.PolicyRule{
			{
				APIGroups: []string{api_v1beta1.SchemeGroupVersion.Group},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"events"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"get", "list", "watch", "create", "patch"},
			},
			{
				APIGroups: []string{storage_api_v1.GroupName},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{crdv1.GroupName},
				Resources: []string{"volumesnapshots"},
				Verbs:     []string{"get"},
			},
		}
		return in
	}, metav1.PatchOptions{})
	return err
}

func (opt *RBACOptions) ensureVolumeSnapshotRestorerJobRoleBinding() error {
	meta := metav1.ObjectMeta{
		Name:      opt.getRoleBindingName(),
		Namespace: opt.Invoker.Namespace,
		Labels:    opt.OffshootLabels,
	}
	_, _, err := rbac_util.CreateOrPatchRoleBinding(context.TODO(), opt.KubeClient, meta, func(in *rbac.RoleBinding) *rbac.RoleBinding {
		core_util.EnsureOwnerReference(&in.ObjectMeta, opt.Owner)

		in.RoleRef = rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     apis.KindClusterRole,
			Name:     apis.StashVolumeSnapshotRestorerClusterRole,
		}
		in.Subjects = []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind,
				Name:      opt.ServiceAccount.Name,
				Namespace: opt.ServiceAccount.Namespace,
			},
		}
		return in
	}, metav1.PatchOptions{})
	return err
}

func (opt *RBACOptions) ensureStorageReaderClassClusterRole() error {
	meta := metav1.ObjectMeta{
		Name:   apis.StashStorageClassReaderClusterRole,
		Labels: opt.OffshootLabels,
	}
	_, _, err := rbac_util.CreateOrPatchClusterRole(context.TODO(), opt.KubeClient, meta, func(in *rbac.ClusterRole) *rbac.ClusterRole {
		in.Rules = []rbac.PolicyRule{
			{
				APIGroups: []string{storage_api_v1.GroupName},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{crdv1.GroupName},
				Resources: []string{"volumesnapshots"},
				Verbs:     []string{"get"},
			},
		}
		return in
	}, metav1.PatchOptions{})
	return err
}

func (opt *RBACOptions) ensureStorageClassReaderClusterRoleBinding() error {
	meta := metav1.ObjectMeta{
		Name:      meta_util.ValidCronJobNameWithSuffix(opt.getRoleBindingName(), apis.StashStorageClassReaderClusterRole),
		Namespace: opt.Invoker.Namespace,
		Labels:    opt.OffshootLabels,
	}
	_, _, err := rbac_util.CreateOrPatchClusterRoleBinding(context.TODO(), opt.KubeClient, meta, func(in *rbac.ClusterRoleBinding) *rbac.ClusterRoleBinding {
		core_util.EnsureOwnerReference(&in.ObjectMeta, opt.Owner)

		in.RoleRef = rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     apis.KindClusterRole,
			Name:     apis.StashStorageClassReaderClusterRole,
		}
		in.Subjects = []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind,
				Name:      opt.ServiceAccount.Name,
				Namespace: opt.ServiceAccount.Namespace,
			},
		}
		return in
	}, metav1.PatchOptions{})
	return err
}
