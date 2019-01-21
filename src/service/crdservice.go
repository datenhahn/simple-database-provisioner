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
	"simple-database-provisioner/src/db"
	"time"
)

// CustomResourceService handles the creation and deletion of
// the custom kubernetes resources SimpleDatabaseBinding and
// SimpleDatabaseInstance.
//go:generate $GOPATH/bin/mockery -name CustomResourceService
type CustomResourceService interface {
	// WasProcessed returns true if an event was already processed
	WasProcessed(eventId string) bool

	// MarkProcessed adds an eventId to the list of processed events
	MarkProcessed(eventId string)

	// CreateDatabaseBinding creates a new database binding
	CreateDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) error

	// DeleteDatabaseBinding deletes a new database binding
	DeleteDatabaseBinding(bindingId db.NamespaceUniqueId) error

	// CreateDatabaseInstance creates a new database instance
	CreateDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) error

	// DeleteDatabaseInstance deletes a database instance
	DeleteDatabaseInstance(instanceId db.NamespaceUniqueId) error

	// UpdateDatabaseBindingState updates
	UpdateDatabaseBindingState(bindingId db.NamespaceUniqueId, newState db.State) error
	UpdateDatabaseInstanceState(instanceId db.NamespaceUniqueId, newState db.State) error
	MarkDatabaseInstanceForDeletion(instanceId db.NamespaceUniqueId) error
	MarkDatabaseBindingForDeletion(bindingId db.NamespaceUniqueId) error
	UpdateDatabaseInstanceCredentials(instanceId db.NamespaceUniqueId, newCredentials map[string][]byte) error
	FindDatabaseInstanceById(instanceId db.NamespaceUniqueId) (db.DatabaseInstance, error)
	FindAllDatabaseInstances() []db.DatabaseInstance
	FindAllDatabaseBindings() []db.DatabaseBinding
	FindInstancesByState(state db.ProvisioningState) []db.DatabaseInstance
	FindBindingsByState(state db.ProvisioningState) []db.DatabaseBinding
}

type PersistentCustomResourceService struct {
	appDatabase db.AppDatabase
}

func NewPersistentCustomResourceService(database db.AppDatabase) CustomResourceService {
	this := &PersistentCustomResourceService{}
	this.appDatabase = database
	return this
}

func (this *PersistentCustomResourceService) WasProcessed(eventId string) bool {

	return this.appDatabase.WasProcessed(eventId)
}

func (this *PersistentCustomResourceService) FindAllDatabaseBindings() []db.DatabaseBinding {

	return this.appDatabase.FindAllDatabaseBindings()
}

func (this *PersistentCustomResourceService) FindAllDatabaseInstances() []db.DatabaseInstance {

	return this.appDatabase.FindAllDatabaseInstances()
}

func (this *PersistentCustomResourceService) MarkProcessed(eventId string) {

	this.appDatabase.MarkProcessed(eventId)
}

func (this *PersistentCustomResourceService) CreateDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) error {

	newBinding := db.DatabaseBinding{
		K8sName:            binding.Name,
		Namespace:          binding.Namespace,
		SecretName:         binding.Spec.SecretName,
		DatabaseInstanceId: db.NewNamespaceUniqueId(binding.Namespace, binding.Spec.InstanceName),
		Meta: db.Meta{
			Current: db.State{
				Action:     db.CREATE,
				State:      db.PENDING,
				Message:    "",
				LastUpdate: time.Now().Round(time.Second),
			},
		},
	}

	return this.appDatabase.AddDatabaseBinding(newBinding)
}

func (this *PersistentCustomResourceService) DeleteDatabaseBinding(bindingId db.NamespaceUniqueId) error {

	return this.appDatabase.DeleteDatabaseBinding(bindingId)
}

func (this *PersistentCustomResourceService) DeleteDatabaseInstance(instanceId db.NamespaceUniqueId) error {

	return this.appDatabase.DeleteDatabaseInstance(instanceId)
}

func (this *PersistentCustomResourceService) FindDatabaseInstanceById(instanceId db.NamespaceUniqueId) (db.DatabaseInstance, error) {

	return this.appDatabase.FindDatabaseInstanceById(instanceId)
}

func (this *PersistentCustomResourceService) MarkDatabaseBindingForDeletion(bindingId db.NamespaceUniqueId) error {

	newState := db.State{
		Action:     db.DELETE,
		State:      db.PENDING,
		Message:    "",
		LastUpdate: time.Now().Round(time.Second),
	}

	return this.appDatabase.UpdateDatabaseBindingState(bindingId, newState)
}

func (this *PersistentCustomResourceService) UpdateDatabaseBindingState(instanceId db.NamespaceUniqueId, newState db.State) error {

	return this.appDatabase.UpdateDatabaseBindingState(instanceId, newState)

}

func (this *PersistentCustomResourceService) FindInstancesByState(state db.ProvisioningState) []db.DatabaseInstance {
	instances := this.appDatabase.FindAllDatabaseInstances()

	matchingInstances := []db.DatabaseInstance{}

	for _, instance := range instances {
		if instance.Meta.Current.State == state {
			matchingInstances = append(matchingInstances, instance)
		}
	}

	return matchingInstances
}

func (this *PersistentCustomResourceService) FindBindingsByState(state db.ProvisioningState) []db.DatabaseBinding {
	bindings := this.appDatabase.FindAllDatabaseBindings()

	matchingBindings := []db.DatabaseBinding{}

	for _, binding := range bindings {
		if binding.Meta.Current.State == state {
			matchingBindings = append(matchingBindings, binding)
		}
	}

	return matchingBindings
}

func (this *PersistentCustomResourceService) CreateDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) error {

	newInstance := db.DatabaseInstance{
		K8sName:      instance.Name,
		Namespace:    instance.Namespace,
		DbmsServer:   instance.Spec.DbmsServer,
		DatabaseName: instance.Spec.DatabaseName,
		Meta: db.Meta{
			Current: db.State{
				Action:     db.CREATE,
				State:      db.PENDING,
				Message:    "",
				LastUpdate: time.Now().Round(time.Second),
			},
		},
	}

	return this.appDatabase.AddDatabaseInstance(newInstance)

}

func (this *PersistentCustomResourceService) UpdateDatabaseInstanceState(instanceId db.NamespaceUniqueId, newState db.State) error {

	return this.appDatabase.UpdateDatabaseInstanceState(instanceId, newState)

}

func (this *PersistentCustomResourceService) UpdateDatabaseInstanceCredentials(instanceId db.NamespaceUniqueId, credentials map[string][]byte) error {

	return this.appDatabase.UpdateDatabaseInstanceCredentials(instanceId, credentials)

}

func (this *PersistentCustomResourceService) MarkDatabaseInstanceForDeletion(instanceId db.NamespaceUniqueId) error {

	newState := db.State{
		Action:     db.DELETE,
		State:      db.PENDING,
		Message:    "",
		LastUpdate: time.Now().Round(time.Second),
	}

	return this.appDatabase.UpdateDatabaseInstanceState(instanceId, newState)
}
