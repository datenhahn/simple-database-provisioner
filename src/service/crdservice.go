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

type CustomResourceDefinitionService interface {
	WasProcessed(uid string) bool
	MarkProcessed(uid string)
	CreateDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) error
	DeleteDatabaseBinding(bindingId string) error
	CreateDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) error
	DeleteDatabaseInstance(instanceId string) error
	UpdateDatabaseBindingState(bindingId string, newState db.State) error
	UpdateDatabaseInstanceState(instanceId string, newState db.State) error
	MarkDatabaseInstanceForDeletion(instanceId string) error
	MarkDatabaseBindingForDeletion(bindingId string) error
	UpdateDatabaseInstanceCredentials(instanceId string, newCredentials map[string][]byte) error
	FindDatabaseInstanceById(instanceId string) (db.DatabaseInstance, error)
	FindAllDatabaseInstances() []db.DatabaseInstance
	FindAllDatabaseBindings() []db.DatabaseBinding
	FindInstancesByState(state db.ProvisioningState) []db.DatabaseInstance
	FindBindingsByState(state db.ProvisioningState) []db.DatabaseBinding
}

type PersistentCustomResourceDefinitionService struct {
	appDatabase db.AppDatabase
}

func NewPersistentCustomResourceDefinitionService(database db.AppDatabase) CustomResourceDefinitionService {
	this := &PersistentCustomResourceDefinitionService{}
	this.appDatabase = database
	return this
}

func (this *PersistentCustomResourceDefinitionService) WasProcessed(uid string) bool {

	return this.appDatabase.WasProcessed(uid)
}

func (this *PersistentCustomResourceDefinitionService) FindAllDatabaseBindings() []db.DatabaseBinding {

	return this.appDatabase.FindAllDatabaseBindings()
}

func (this *PersistentCustomResourceDefinitionService) FindAllDatabaseInstances() []db.DatabaseInstance {

	return this.appDatabase.FindAllDatabaseInstances()
}

func (this *PersistentCustomResourceDefinitionService) MarkProcessed(uid string) {

	this.appDatabase.MarkProcessed(uid)
}

func (this *PersistentCustomResourceDefinitionService) CreateDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) error {

	newBinding := db.DatabaseBinding{
		Id:                 binding.Name,
		Namespace:          binding.Namespace,
		SecretName:         binding.Spec.SecretName,
		DatabaseInstanceId: binding.Spec.InstanceName,
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

func (this *PersistentCustomResourceDefinitionService) DeleteDatabaseBinding(bindingId string) error {

	return this.appDatabase.DeleteDatabaseBinding(bindingId)
}

func (this *PersistentCustomResourceDefinitionService) DeleteDatabaseInstance(instanceId string) error {

	return this.appDatabase.DeleteDatabaseInstance(instanceId)
}

func (this *PersistentCustomResourceDefinitionService) FindDatabaseInstanceById(instanceId string) (db.DatabaseInstance, error) {

	return this.appDatabase.FindDatabaseInstanceById(instanceId)
}

func (this *PersistentCustomResourceDefinitionService) MarkDatabaseBindingForDeletion(bindingId string) error {

	newState := db.State{
		Action:     db.DELETE,
		State:      db.PENDING,
		Message:    "",
		LastUpdate: time.Now().Round(time.Second),
	}

	return this.appDatabase.UpdateDatabaseBindingState(bindingId, newState)
}

func (this *PersistentCustomResourceDefinitionService) UpdateDatabaseBindingState(instanceId string, newState db.State) error {

	return this.appDatabase.UpdateDatabaseBindingState(instanceId, newState)

}

func (this *PersistentCustomResourceDefinitionService) FindInstancesByState(state db.ProvisioningState) []db.DatabaseInstance {
	instances := this.appDatabase.FindAllDatabaseInstances()

	matchingInstances := []db.DatabaseInstance{}

	for _, instance := range instances {
		if instance.Meta.Current.State == state {
			matchingInstances = append(matchingInstances, instance)
		}
	}

	return matchingInstances
}

func (this *PersistentCustomResourceDefinitionService) FindBindingsByState(state db.ProvisioningState) []db.DatabaseBinding {
	bindings := this.appDatabase.FindAllDatabaseBindings()

	matchingBindings := []db.DatabaseBinding{}

	for _, binding := range bindings {
		if binding.Meta.Current.State == state {
			matchingBindings = append(matchingBindings, binding)
		}
	}

	return matchingBindings
}

func (this *PersistentCustomResourceDefinitionService) CreateDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) error {

	newInstance := db.DatabaseInstance{
		Id:           instance.Name,
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

func (this *PersistentCustomResourceDefinitionService) UpdateDatabaseInstanceState(instanceId string, newState db.State) error {

	return this.appDatabase.UpdateDatabaseInstanceState(instanceId, newState)

}

func (this *PersistentCustomResourceDefinitionService) UpdateDatabaseInstanceCredentials(instanceId string, credentials map[string][]byte) error {

	return this.appDatabase.UpdateDatabaseInstanceCredentials(instanceId, credentials)

}

func (this *PersistentCustomResourceDefinitionService) MarkDatabaseInstanceForDeletion(instanceId string) error {

	newState := db.State{
		Action:     db.DELETE,
		State:      db.PENDING,
		Message:    "",
		LastUpdate: time.Now().Round(time.Second),
	}

	return this.appDatabase.UpdateDatabaseInstanceState(instanceId, newState)
}
