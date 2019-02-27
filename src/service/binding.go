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

package service

import (
	"simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	"simple-database-provisioner/src/persistence"
	"time"
)

// CustomResourceService handles the creation and deletion of
// the custom kubernetes resources SimpleDatabaseBinding and
// SimpleDatabaseInstance.
//go:generate $GOPATH/bin/mockery -name DatabaseBindingService
type DatabaseBindingService interface {
	// CreateDatabaseBinding creates a new database binding
	CreateDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) error

	// DeleteDatabaseBinding deletes a new database binding
	DeleteDatabaseBinding(bindingId persistence.NamespaceUniqueId) error

	// UpdateDatabaseBindingState updates a database binding's state
	UpdateDatabaseBindingState(bindingId persistence.NamespaceUniqueId, newState persistence.State) error

	// MarkDatabaseBindingForDeletion marks a database binding for deletion
	MarkDatabaseBindingForDeletion(bindingId persistence.NamespaceUniqueId) error

	// FindAllDatabaseBindings returns all database bindings
	FindAllDatabaseBindings() []persistence.DatabaseBinding

	// FindBindingsByState returns all database bindings in a certain state
	FindBindingsByState(state persistence.ProvisioningState) []persistence.DatabaseBinding
}

type PersistentDatabaseBindingService struct {
	bindingRepo persistence.DatabaseBindingRepository
}

func NewPersistentDatabaseBindingService(bindingRepo persistence.DatabaseBindingRepository) DatabaseBindingService {
	this := &PersistentDatabaseBindingService{}
	this.bindingRepo = bindingRepo
	return this
}

func (this *PersistentDatabaseBindingService) FindAllDatabaseBindings() []persistence.DatabaseBinding {

	return this.bindingRepo.FindAllDatabaseBindings()
}

func (this *PersistentDatabaseBindingService) CreateDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) error {

	newBinding := persistence.DatabaseBinding{
		K8sName:            binding.Name,
		Namespace:          binding.Namespace,
		SecretName:         binding.Spec.SecretName,
		DatabaseInstanceId: persistence.NewNamespaceUniqueId(binding.Namespace, binding.Spec.InstanceName),
		Meta: persistence.Meta{
			Current: persistence.State{
				Action:     persistence.CREATE,
				State:      persistence.PENDING,
				Message:    "",
				LastUpdate: time.Now().Round(time.Second),
			},
		},
	}

	return this.bindingRepo.AddDatabaseBinding(newBinding)
}

func (this *PersistentDatabaseBindingService) DeleteDatabaseBinding(bindingId persistence.NamespaceUniqueId) error {

	return this.bindingRepo.DeleteDatabaseBinding(bindingId)
}

func (this *PersistentDatabaseBindingService) MarkDatabaseBindingForDeletion(bindingId persistence.NamespaceUniqueId) error {

	newState := persistence.State{
		Action:     persistence.DELETE,
		State:      persistence.PENDING,
		Message:    "",
		LastUpdate: time.Now().Round(time.Second),
	}

	return this.bindingRepo.UpdateDatabaseBindingState(bindingId, newState)
}

func (this *PersistentDatabaseBindingService) UpdateDatabaseBindingState(instanceId persistence.NamespaceUniqueId, newState persistence.State) error {

	return this.bindingRepo.UpdateDatabaseBindingState(instanceId, newState)

}

func (this *PersistentDatabaseBindingService) FindBindingsByState(state persistence.ProvisioningState) []persistence.DatabaseBinding {
	bindings := this.bindingRepo.FindAllDatabaseBindings()

	matchingBindings := []persistence.DatabaseBinding{}

	for _, binding := range bindings {
		if binding.Meta.Current.State == state {
			matchingBindings = append(matchingBindings, binding)
		}
	}

	return matchingBindings
}
