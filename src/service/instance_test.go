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
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	"simple-database-provisioner/src/persistence"
	"simple-database-provisioner/src/persistence/mocks"
	"testing"
	"time"
)

func TestPersistentDatabaseInstanceService_CreateDatabaseInstanceService(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	expectedInstance := &persistence.DatabaseInstance{
		Namespace:         "testns",
		DatabaseName:      "mytest-database",
		NamespaceUniqueId: persistence.NamespaceUniqueId("testns-mytest-database-instance"),
		DbmsServer:        "my-dbms",
		K8sName:           "mytest-database-instance",
	}

	mockRepo.On("AddDatabaseInstance", mock.MatchedBy(func(instance persistence.DatabaseInstance) bool {
		assert.Equal(t, expectedInstance.Namespace, instance.Namespace)
		assert.Equal(t, expectedInstance.DatabaseName, instance.DatabaseName)
		assert.Equal(t, expectedInstance.NamespaceUniqueId, instance.NamespaceUniqueId)
		assert.Equal(t, expectedInstance.DbmsServer, instance.DbmsServer)
		assert.Equal(t, expectedInstance.K8sName, instance.K8sName)
		return true
	})).Return(nil)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)

	instance := &v1alpha1.SimpleDatabaseInstance{}
	instance.Name = "mytest-database-instance"
	instance.Namespace = "testns"
	instance.Spec = v1alpha1.SimpleDatabaseInstanceSpec{
		DatabaseName: "mytest-database",
		DbmsServer:   "my-dbms",
	}

	err := instanceService.CreateDatabaseInstance(instance)

	assert.Nil(t, err)
}

func TestPersistentDatabaseInstanceService_MarkDatabaseInstanceForDeletion(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	mockRepo.On("UpdateDatabaseInstanceState", mock.MatchedBy(func(instanceId persistence.NamespaceUniqueId) bool {
		assert.Equal(t, persistence.NamespaceUniqueId("testns-mytest-database-instance"), instanceId)
		return true
	}), mock.MatchedBy(func(state persistence.State) bool {
		assert.Equal(t, persistence.DELETE, state.Action)
		assert.Equal(t, persistence.PENDING, state.State)
		return true
	})).Return(nil)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)
	err := instanceService.MarkDatabaseInstanceForDeletion(persistence.NamespaceUniqueId("testns-mytest-database-instance"))

	assert.Nil(t, err)

}

func TestPersistentDatabaseInstanceService_DeleteDatabaseInstance(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	mockRepo.On("DeleteDatabaseInstance", mock.MatchedBy(func(instanceId persistence.NamespaceUniqueId) bool {
		assert.Equal(t, persistence.NamespaceUniqueId("testns-mytest-database-instance"), instanceId)
		return true
	})).Return(nil)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)
	err := instanceService.DeleteDatabaseInstance(persistence.NamespaceUniqueId("testns-mytest-database-instance"))

	assert.Nil(t, err)

}

func TestPersistentDatabaseInstanceService_FindAllDatabaseInstances(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	expectedInstance := persistence.DatabaseInstance{
		Namespace:         "testns",
		DatabaseName:      "mytest-database",
		NamespaceUniqueId: persistence.NamespaceUniqueId("testns-mytest-database-instance"),
		DbmsServer:        "my-dbms",
		K8sName:           "mytest-database-instance",
	}

	instances := []persistence.DatabaseInstance{expectedInstance}

	mockRepo.On("FindAllDatabaseInstances").Return(instances)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)
	result := instanceService.FindAllDatabaseInstances()

	assert.Equal(t, instances, result)
}

func TestPersistentDatabaseInstanceService_FindInstancesByState(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	str := "2019-04-12T11:45:26.371Z"
	lastUpdate, _ := time.Parse(time.RFC3339, str)

	pendingInstance := persistence.DatabaseInstance{
		Namespace:         "testns",
		DatabaseName:      "mytest-database",
		NamespaceUniqueId: persistence.NamespaceUniqueId("testns-mytest-database-instance"),
		DbmsServer:        "my-dbms",
		K8sName:           "mytest-database-instance",
		Meta: persistence.Meta{
			Previous: persistence.State{
				State:      persistence.READY,
				Action:     persistence.CREATE,
				LastUpdate: lastUpdate,
				Message:    "OK",
			},
			Current: persistence.State{
				State:      persistence.PENDING,
				Action:     persistence.CREATE,
				LastUpdate: lastUpdate,
				Message:    "PENDING",
			},
		},
	}

	readyInstance := persistence.DatabaseInstance{
		Namespace:         "testns",
		DatabaseName:      "mytest-database",
		NamespaceUniqueId: persistence.NamespaceUniqueId("testns-mytest2-database-instance"),
		DbmsServer:        "my-dbms",
		K8sName:           "mytest2-database-instance",
		Meta: persistence.Meta{
			Previous: persistence.State{
				State:      persistence.PENDING,
				Action:     persistence.CREATE,
				LastUpdate: lastUpdate,
				Message:    "PENDING",
			},
			Current: persistence.State{
				State:      persistence.READY,
				Action:     persistence.CREATE,
				LastUpdate: lastUpdate,
				Message:    "OK",
			},
		},
	}

	instances := []persistence.DatabaseInstance{readyInstance, pendingInstance}

	mockRepo.On("FindAllDatabaseInstances").Return(instances)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)
	result := instanceService.FindInstancesByState(persistence.READY)

	assert.Equal(t, []persistence.DatabaseInstance{readyInstance}, result)

	result2 := instanceService.FindInstancesByState(persistence.PENDING)

	assert.Equal(t, []persistence.DatabaseInstance{pendingInstance}, result2)
}

func TestPersistentDatabaseInstanceService_FindDatabaseInstanceById(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	str := "2019-04-12T11:45:26.371Z"
	lastUpdate, _ := time.Parse(time.RFC3339, str)

	pendingInstance := persistence.DatabaseInstance{
		Namespace:         "testns",
		DatabaseName:      "mytest-database",
		NamespaceUniqueId: persistence.NamespaceUniqueId("testns-mytest-database-instance"),
		DbmsServer:        "my-dbms",
		K8sName:           "mytest-database-instance",
		Meta: persistence.Meta{
			Previous: persistence.State{
				State:      persistence.READY,
				Action:     persistence.CREATE,
				LastUpdate: lastUpdate,
				Message:    "OK",
			},
			Current: persistence.State{
				State:      persistence.PENDING,
				Action:     persistence.CREATE,
				LastUpdate: lastUpdate,
				Message:    "PENDING",
			},
		},
	}

	mockRepo.On("FindDatabaseInstanceById", persistence.NamespaceUniqueId("testns-mytest-database-instance")).Return(pendingInstance, nil)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)
	result, err := instanceService.FindDatabaseInstanceById(persistence.NamespaceUniqueId("testns-mytest-database-instance"))

	assert.Nil(t, err)
	assert.Equal(t, pendingInstance, result)

	mockRepo.On("FindDatabaseInstanceById", persistence.NamespaceUniqueId("testns-mytest-database-instance123")).Return(persistence.DatabaseInstance{}, fmt.Errorf("Some error"))

	result2, err2 := instanceService.FindDatabaseInstanceById(persistence.NamespaceUniqueId("testns-mytest-database-instance123"))

	assert.Equal(t, err2.Error(), "Some error")
	assert.Equal(t, result2, persistence.DatabaseInstance{})

}

func TestPersistentDatabaseInstanceService_UpdateDatabaseInstanceState(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	str := "2019-04-12T11:45:26.371Z"
	lastUpdate, _ := time.Parse(time.RFC3339, str)

	newState := persistence.State{
		State:      persistence.READY,
		Action:     persistence.CREATE,
		LastUpdate: lastUpdate,
		Message:    "OK",
	}

	mockRepo.On("UpdateDatabaseInstanceState", persistence.NamespaceUniqueId("testns-mytest-database-instance"), newState).Return(nil)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)

	err := instanceService.UpdateDatabaseInstanceState(persistence.NamespaceUniqueId("testns-mytest-database-instance"), newState)

	assert.Nil(t, err)
}

func TestPersistentDatabaseInstanceService_UpdateDatabaseInstanceCredentials(t *testing.T) {
	mockRepo := &mocks.DatabaseInstanceRepository{}

	newCredentials := map[string][]byte{
		"name":     []byte("my-db"),
		"username": []byte("other"),
	}

	mockRepo.On("UpdateDatabaseInstanceCredentials", persistence.NamespaceUniqueId("testns-mytest-database-instance"), newCredentials).Return(nil)

	instanceService := NewPersistentDatabaseInstanceService(mockRepo)

	err := instanceService.UpdateDatabaseInstanceCredentials(persistence.NamespaceUniqueId("testns-mytest-database-instance"), newCredentials)

	assert.Nil(t, err)
}
