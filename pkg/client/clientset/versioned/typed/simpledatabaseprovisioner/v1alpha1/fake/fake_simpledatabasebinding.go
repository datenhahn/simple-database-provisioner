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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSimpleDatabaseBindings implements SimpleDatabaseBindingInterface
type FakeSimpleDatabaseBindings struct {
	Fake *FakeSimpledatabaseprovisionerV1alpha1
	ns   string
}

var simpledatabasebindingsResource = schema.GroupVersionResource{Group: "simpledatabaseprovisioner.k8s.ecodia.de", Version: "v1alpha1", Resource: "simpledatabasebindings"}

var simpledatabasebindingsKind = schema.GroupVersionKind{Group: "simpledatabaseprovisioner.k8s.ecodia.de", Version: "v1alpha1", Kind: "SimpleDatabaseBinding"}

// Get takes name of the simpleDatabaseBinding, and returns the corresponding simpleDatabaseBinding object, and an error if there is any.
func (c *FakeSimpleDatabaseBindings) Get(name string, options v1.GetOptions) (result *v1alpha1.SimpleDatabaseBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(simpledatabasebindingsResource, c.ns, name), &v1alpha1.SimpleDatabaseBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SimpleDatabaseBinding), err
}

// List takes label and field selectors, and returns the list of SimpleDatabaseBindings that match those selectors.
func (c *FakeSimpleDatabaseBindings) List(opts v1.ListOptions) (result *v1alpha1.SimpleDatabaseBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(simpledatabasebindingsResource, simpledatabasebindingsKind, c.ns, opts), &v1alpha1.SimpleDatabaseBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SimpleDatabaseBindingList{ListMeta: obj.(*v1alpha1.SimpleDatabaseBindingList).ListMeta}
	for _, item := range obj.(*v1alpha1.SimpleDatabaseBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested simpleDatabaseBindings.
func (c *FakeSimpleDatabaseBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(simpledatabasebindingsResource, c.ns, opts))

}

// Create takes the representation of a simpleDatabaseBinding and creates it.  Returns the server's representation of the simpleDatabaseBinding, and an error, if there is any.
func (c *FakeSimpleDatabaseBindings) Create(simpleDatabaseBinding *v1alpha1.SimpleDatabaseBinding) (result *v1alpha1.SimpleDatabaseBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(simpledatabasebindingsResource, c.ns, simpleDatabaseBinding), &v1alpha1.SimpleDatabaseBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SimpleDatabaseBinding), err
}

// Update takes the representation of a simpleDatabaseBinding and updates it. Returns the server's representation of the simpleDatabaseBinding, and an error, if there is any.
func (c *FakeSimpleDatabaseBindings) Update(simpleDatabaseBinding *v1alpha1.SimpleDatabaseBinding) (result *v1alpha1.SimpleDatabaseBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(simpledatabasebindingsResource, c.ns, simpleDatabaseBinding), &v1alpha1.SimpleDatabaseBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SimpleDatabaseBinding), err
}

// Delete takes name of the simpleDatabaseBinding and deletes it. Returns an error if one occurs.
func (c *FakeSimpleDatabaseBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(simpledatabasebindingsResource, c.ns, name), &v1alpha1.SimpleDatabaseBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSimpleDatabaseBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(simpledatabasebindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SimpleDatabaseBindingList{})
	return err
}

// Patch applies the patch and returns the patched simpleDatabaseBinding.
func (c *FakeSimpleDatabaseBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SimpleDatabaseBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(simpledatabasebindingsResource, c.ns, name, pt, data, subresources...), &v1alpha1.SimpleDatabaseBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SimpleDatabaseBinding), err
}
