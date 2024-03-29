/*
 * Copyright (c) 2019 Ecodia GmbH & Co. KG <opensource@ecodia.de>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	simpledatabaseprovisionerv1alpha1 "simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	versioned "simple-database-provisioner/pkg/client/clientset/versioned"
	internalinterfaces "simple-database-provisioner/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "simple-database-provisioner/pkg/client/listers/simpledatabaseprovisioner/v1alpha1"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SimpleDatabaseBindingInformer provides access to a shared informer and lister for
// SimpleDatabaseBindings.
type SimpleDatabaseBindingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.SimpleDatabaseBindingLister
}

type simpleDatabaseBindingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSimpleDatabaseBindingInformer constructs a new informer for SimpleDatabaseBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSimpleDatabaseBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSimpleDatabaseBindingInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSimpleDatabaseBindingInformer constructs a new informer for SimpleDatabaseBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSimpleDatabaseBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SimpledatabaseprovisionerV1alpha1().SimpleDatabaseBindings(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SimpledatabaseprovisionerV1alpha1().SimpleDatabaseBindings(namespace).Watch(options)
			},
		},
		&simpledatabaseprovisionerv1alpha1.SimpleDatabaseBinding{},
		resyncPeriod,
		indexers,
	)
}

func (f *simpleDatabaseBindingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSimpleDatabaseBindingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *simpleDatabaseBindingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&simpledatabaseprovisionerv1alpha1.SimpleDatabaseBinding{}, f.defaultInformer)
}

func (f *simpleDatabaseBindingInformer) Lister() v1alpha1.SimpleDatabaseBindingLister {
	return v1alpha1.NewSimpleDatabaseBindingLister(f.Informer().GetIndexer())
}
