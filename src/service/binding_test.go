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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	"simple-database-provisioner/src/persistence"
	"simple-database-provisioner/src/persistence/mocks"
	"testing"
	"time"
)

func TestNewPersistentDatabaseBindingService(t *testing.T) {

	mockRepo := &mocks.DatabaseBindingRepository{}

	bindingService := NewPersistentDatabaseBindingService(mockRepo)
	assert.NotNil(t, bindingService)
}

func TestPersistentDatabaseBindingService_CreateDatabaseBinding(t *testing.T) {
	mockRepo := &mocks.DatabaseBindingRepository{}

	expectedBinding := &persistence.DatabaseBinding{
		Namespace:          "testns",
		SecretName:         "mytest-database-secret",
		NamespaceUniqueId:  persistence.NamespaceUniqueId("testns-mytest-database-binding"),
		DatabaseInstanceId: persistence.NamespaceUniqueId("testns-mytest"),
		K8sName:            "mytest-database-binding",
	}

	mockRepo.On("AddDatabaseBinding", mock.MatchedBy(func(binding persistence.DatabaseBinding) bool {
		assert.Equal(t, expectedBinding.Namespace, binding.Namespace)
		assert.Equal(t, expectedBinding.SecretName, binding.SecretName)
		assert.Equal(t, expectedBinding.NamespaceUniqueId, binding.NamespaceUniqueId)
		assert.Equal(t, expectedBinding.DatabaseInstanceId, binding.DatabaseInstanceId)
		assert.Equal(t, expectedBinding.K8sName, binding.K8sName)
		return true
	})).Return(nil)

	bindingService := NewPersistentDatabaseBindingService(mockRepo)

	binding := &v1alpha1.SimpleDatabaseBinding{}
	binding.Name = "mytest-database-binding"
	binding.Namespace = "testns"
	binding.Spec = v1alpha1.SimpleDatabaseBindingSpec{
		SecretName:   "mytest-database-secret",
		InstanceName: "mytest",
	}

	err := bindingService.CreateDatabaseBinding(binding)

	assert.Nil(t, err)

}

func TestPersistentDatabaseBindingService_MarkDatabaseBindingForDeletion(t *testing.T) {
	mockRepo := &mocks.DatabaseBindingRepository{}

	mockRepo.On("UpdateDatabaseBindingState", mock.MatchedBy(func(bindingId persistence.NamespaceUniqueId) bool {
		assert.Equal(t, persistence.NamespaceUniqueId("testns-mytest-database-binding"), bindingId)
		return true
	}), mock.MatchedBy(func(state persistence.State) bool {
		assert.Equal(t, persistence.DELETE, state.Action)
		assert.Equal(t, persistence.PENDING, state.State)
		return true
	})).Return(nil)

	bindingService := NewPersistentDatabaseBindingService(mockRepo)
	err := bindingService.MarkDatabaseBindingForDeletion(persistence.NamespaceUniqueId("testns-mytest-database-binding"))

	assert.Nil(t, err)

}

func TestPersistentDatabaseBindingService_DeleteDatabaseBinding(t *testing.T) {
	mockRepo := &mocks.DatabaseBindingRepository{}

	mockRepo.On("DeleteDatabaseBinding", mock.MatchedBy(func(bindingId persistence.NamespaceUniqueId) bool {
		assert.Equal(t, persistence.NamespaceUniqueId("testns-mytest-database-binding"), bindingId)
		return true
	})).Return(nil)

	bindingService := NewPersistentDatabaseBindingService(mockRepo)
	err := bindingService.DeleteDatabaseBinding(persistence.NamespaceUniqueId("testns-mytest-database-binding"))

	assert.Nil(t, err)

}

func TestPersistentDatabaseBindingService_FindAllDatabaseBindings(t *testing.T) {
	mockRepo := &mocks.DatabaseBindingRepository{}

	expectedBinding := persistence.DatabaseBinding{
		Namespace:          "testns",
		SecretName:         "mytest-database-secret",
		NamespaceUniqueId:  persistence.NamespaceUniqueId("testns-mytest-database-binding"),
		DatabaseInstanceId: persistence.NamespaceUniqueId("testns-mytest"),
		K8sName:            "mytest-database-binding",
	}

	bindings := []persistence.DatabaseBinding{expectedBinding}

	mockRepo.On("FindAllDatabaseBindings").Return(bindings)

	bindingService := NewPersistentDatabaseBindingService(mockRepo)
	result := bindingService.FindAllDatabaseBindings()

	assert.Equal(t, bindings, result)
}

func TestPersistentDatabaseBindingService_FindBindingsByState(t *testing.T) {
	mockRepo := &mocks.DatabaseBindingRepository{}

	str := "2019-04-12T11:45:26.371Z"
	lastUpdate, _ := time.Parse(time.RFC3339, str)

	pendingBinding := persistence.DatabaseBinding{
		Namespace:          "testns",
		SecretName:         "mytest-database-secret",
		NamespaceUniqueId:  persistence.NamespaceUniqueId("testns-mytest-database-binding"),
		DatabaseInstanceId: persistence.NamespaceUniqueId("testns-mytest"),
		K8sName:            "mytest-database-binding",
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

	readyBinding := persistence.DatabaseBinding{
		Namespace:          "testns",
		SecretName:         "mytest2-database-secret",
		NamespaceUniqueId:  persistence.NamespaceUniqueId("testns-mytest2-database-binding"),
		DatabaseInstanceId: persistence.NamespaceUniqueId("testns-mytest2"),
		K8sName:            "mytest2-database-binding",
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

	bindings := []persistence.DatabaseBinding{readyBinding, pendingBinding}

	mockRepo.On("FindAllDatabaseBindings").Return(bindings)

	bindingService := NewPersistentDatabaseBindingService(mockRepo)
	result := bindingService.FindBindingsByState(persistence.READY)

	assert.Equal(t, []persistence.DatabaseBinding{readyBinding}, result)

	result2 := bindingService.FindBindingsByState(persistence.PENDING)

	assert.Equal(t, []persistence.DatabaseBinding{pendingBinding}, result2)
}

func TestPersistentDatabaseBindingService_UpdateDatabaseBindingState(t *testing.T) {
	mockRepo := &mocks.DatabaseBindingRepository{}

	str := "2019-04-12T11:45:26.371Z"
	lastUpdate, _ := time.Parse(time.RFC3339, str)

	newState := persistence.State{
		State:      persistence.READY,
		Action:     persistence.CREATE,
		LastUpdate: lastUpdate,
		Message:    "OK",
	}

	mockRepo.On("UpdateDatabaseBindingState", persistence.NamespaceUniqueId("testns-mytest-database-binding"), newState).Return(nil)

	bindingService := NewPersistentDatabaseBindingService(mockRepo)

	err := bindingService.UpdateDatabaseBindingState(persistence.NamespaceUniqueId("testns-mytest-database-binding"), newState)

	assert.Nil(t, err)
}
