// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	appsv1 "kmodules.xyz/openshift/apis/apps/v1"
	versioned "kmodules.xyz/openshift/client/clientset/versioned"
	internalinterfaces "kmodules.xyz/openshift/client/informers/externalversions/internalinterfaces"
	v1 "kmodules.xyz/openshift/client/listers/apps/v1"
)

// DeploymentConfigInformer provides access to a shared informer and lister for
// DeploymentConfigs.
type DeploymentConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.DeploymentConfigLister
}

type deploymentConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewDeploymentConfigInformer constructs a new informer for DeploymentConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDeploymentConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDeploymentConfigInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredDeploymentConfigInformer constructs a new informer for DeploymentConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDeploymentConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppsV1().DeploymentConfigs(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppsV1().DeploymentConfigs(namespace).Watch(options)
			},
		},
		&appsv1.DeploymentConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *deploymentConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDeploymentConfigInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *deploymentConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&appsv1.DeploymentConfig{}, f.defaultInformer)
}

func (f *deploymentConfigInformer) Lister() v1.DeploymentConfigLister {
	return v1.NewDeploymentConfigLister(f.Informer().GetIndexer())
}
