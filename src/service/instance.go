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
//go:generate $GOPATH/bin/mockery -name DatabaseInstanceService
type DatabaseInstanceService interface {
	// CreateDatabaseInstance creates a new database instance
	CreateDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) error

	// DeleteDatabaseInstance deletes a database instance
	DeleteDatabaseInstance(instanceId persistence.NamespaceUniqueId) error

	// UpdateDatabaseInstanceState updates a database instance's state
	UpdateDatabaseInstanceState(instanceId persistence.NamespaceUniqueId, newState persistence.State) error

	// MarkDatabaseInstanceForDeletion marks a database instance for deletion
	MarkDatabaseInstanceForDeletion(instanceId persistence.NamespaceUniqueId) error

	// UpdateDatabaseInstanceCredentials updates a database instance's credentials in the database
	UpdateDatabaseInstanceCredentials(instanceId persistence.NamespaceUniqueId, newCredentials map[string][]byte) error

	// FindDatabaseInstanceById finds a database instance with a specific namespace unique id
	FindDatabaseInstanceById(instanceId persistence.NamespaceUniqueId) (persistence.DatabaseInstance, error)

	// FindAllDatabaseInstances returns all database instances
	FindAllDatabaseInstances() []persistence.DatabaseInstance

	// FindInstancesByState returns all database instances in a certain state
	FindInstancesByState(state persistence.ProvisioningState) []persistence.DatabaseInstance
}

type PersistentDatabaseInstanceService struct {
	instanceRepo persistence.DatabaseInstanceRepository
}

func NewPersistentDatabaseInstanceService(instanceRepo persistence.DatabaseInstanceRepository) DatabaseInstanceService {
	this := &PersistentDatabaseInstanceService{}
	this.instanceRepo = instanceRepo
	return this
}

func (this *PersistentDatabaseInstanceService) FindAllDatabaseInstances() []persistence.DatabaseInstance {

	return this.instanceRepo.FindAllDatabaseInstances()
}

func (this *PersistentDatabaseInstanceService) DeleteDatabaseInstance(instanceId persistence.NamespaceUniqueId) error {

	return this.instanceRepo.DeleteDatabaseInstance(instanceId)
}

func (this *PersistentDatabaseInstanceService) FindDatabaseInstanceById(instanceId persistence.NamespaceUniqueId) (persistence.DatabaseInstance, error) {

	return this.instanceRepo.FindDatabaseInstanceById(instanceId)
}

func (this *PersistentDatabaseInstanceService) FindInstancesByState(state persistence.ProvisioningState) []persistence.DatabaseInstance {
	instances := this.instanceRepo.FindAllDatabaseInstances()

	matchingInstances := []persistence.DatabaseInstance{}

	for _, instance := range instances {
		if instance.Meta.Current.State == state {
			matchingInstances = append(matchingInstances, instance)
		}
	}

	return matchingInstances
}

func (this *PersistentDatabaseInstanceService) CreateDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) error {

	newInstance := persistence.DatabaseInstance{
		K8sName:      instance.Name,
		Namespace:    instance.Namespace,
		DbmsServer:   instance.Spec.DbmsServer,
		DatabaseName: instance.Spec.DatabaseName,
		Meta: persistence.Meta{
			Current: persistence.State{
				Action:     persistence.CREATE,
				State:      persistence.PENDING,
				Message:    "",
				LastUpdate: time.Now().Round(time.Second),
			},
		},
	}

	newInstance.NamespaceUniqueId = newInstance.GetNamespaceUniqueId()

	return this.instanceRepo.AddDatabaseInstance(newInstance)

}

func (this *PersistentDatabaseInstanceService) UpdateDatabaseInstanceState(instanceId persistence.NamespaceUniqueId, newState persistence.State) error {

	return this.instanceRepo.UpdateDatabaseInstanceState(instanceId, newState)

}

func (this *PersistentDatabaseInstanceService) UpdateDatabaseInstanceCredentials(instanceId persistence.NamespaceUniqueId, credentials map[string][]byte) error {

	return this.instanceRepo.UpdateDatabaseInstanceCredentials(instanceId, credentials)

}

func (this *PersistentDatabaseInstanceService) MarkDatabaseInstanceForDeletion(instanceId persistence.NamespaceUniqueId) error {

	newState := persistence.State{
		Action:     persistence.DELETE,
		State:      persistence.PENDING,
		Message:    "",
		LastUpdate: time.Now().Round(time.Second),
	}

	return this.instanceRepo.UpdateDatabaseInstanceState(instanceId, newState)
}
