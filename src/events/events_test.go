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

package events

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	"simple-database-provisioner/src/events/mocks"
	"simple-database-provisioner/src/persistence"
	"simple-database-provisioner/src/service"
	"simple-database-provisioner/src/util"
	"testing"
	"time"
)

func createBinding(name string) *v1alpha1.SimpleDatabaseBinding {
	return &v1alpha1.SimpleDatabaseBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "mytest",
			UID:       types.UID(util.Md5(fmt.Sprintf("%d", time.Now().UnixNano()))),
		},
	}
}

func createInstance(name string) *v1alpha1.SimpleDatabaseInstance {
	return &v1alpha1.SimpleDatabaseInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			UID:  types.UID(util.Md5(fmt.Sprintf("%d", time.Now().UnixNano()))),
		},
	}
}

func TestEventHandler_fast_subsequent_binding_add_and_removes(t *testing.T) {

	dbName := "/tmp/TestNewGoSimpleDatabaseProvisionerEventHandler.db"
	db, err := storm.Open(dbName)
	assert.Nil(t, err)
	defer os.Remove(dbName)

	stormDb := persistence.NewStormPersistenceBackend(db)

	processor := &mocks.ProvisioningEventProcessor{}
	processor.On("ProcessEvents").Return()

	eventService := service.NewPersistentEventService(stormDb)
	bindingService := service.NewPersistentDatabaseBindingService(stormDb)
	instanceService := service.NewPersistentDatabaseInstanceService(stormDb)
	handler := NewGoSimpleDatabaseProvisionerEventHandler(eventService, bindingService, instanceService, processor)

	alpha := createBinding("alpha")
	alpha2 := createBinding("alpha")

	handler.OnAddDatabaseBinding(alpha)
	bindings := bindingService.FindAllDatabaseBindings()
	assert.Equal(t, alpha.Name, bindings[0].K8sName)
	assert.Equal(t, string(persistence.CREATE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(bindings[0].Meta.Current.State))

	handler.OnDeleteDatabaseBinding(alpha)
	time.Sleep(1 * time.Second)
	bindings = bindingService.FindAllDatabaseBindings()
	assert.Equal(t, alpha.Name, bindings[0].K8sName)
	assert.Equal(t, string(persistence.DELETE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(bindings[0].Meta.Current.State))

	handler.OnAddDatabaseBinding(alpha2)
	bindings = bindingService.FindAllDatabaseBindings()
	assert.Equal(t, alpha2.Name, bindings[0].K8sName)
	assert.Equal(t, string(persistence.CREATE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(bindings[0].Meta.Current.State))

	handler.OnDeleteDatabaseBinding(alpha2)
	bindings = bindingService.FindAllDatabaseBindings()
	assert.Equal(t, alpha2.Name, bindings[0].K8sName)
	assert.Equal(t, string(persistence.DELETE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(bindings[0].Meta.Current.State))

}

func TestEventHandler_reprocess_same_event(t *testing.T) {

	dbName := "/tmp/TestNewGoSimpleDatabaseProvisionerEventHandler.db"
	db, err := storm.Open(dbName)
	assert.Nil(t, err)
	defer os.Remove(dbName)

	stormDb := persistence.NewStormPersistenceBackend(db)

	processor := &mocks.ProvisioningEventProcessor{}
	processor.On("ProcessEvents").Return()

	eventService := service.NewPersistentEventService(stormDb)
	bindingService := service.NewPersistentDatabaseBindingService(stormDb)
	instanceService := service.NewPersistentDatabaseInstanceService(stormDb)
	handler := NewGoSimpleDatabaseProvisionerEventHandler(eventService, bindingService, instanceService, processor)

	alpha := createBinding("alpha")

	handler.OnAddDatabaseBinding(alpha)
	handler.OnAddDatabaseBinding(alpha)
	bindings := bindingService.FindAllDatabaseBindings()
	assert.Equal(t, alpha.Name, bindings[0].K8sName)
	assert.Equal(t, string(persistence.CREATE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(bindings[0].Meta.Current.State))
}

func TestEventHandler_reprocess_same_instance_event(t *testing.T) {

	dbName := "/tmp/TestNewGoSimpleDatabaseProvisionerEventHandler.db"
	db, err := storm.Open(dbName)
	assert.Nil(t, err)
	defer os.Remove(dbName)

	stormDb := persistence.NewStormPersistenceBackend(db)

	processor := &mocks.ProvisioningEventProcessor{}
	processor.On("ProcessEvents").Return()

	eventService := service.NewPersistentEventService(stormDb)
	bindingService := service.NewPersistentDatabaseBindingService(stormDb)
	instanceService := service.NewPersistentDatabaseInstanceService(stormDb)
	handler := NewGoSimpleDatabaseProvisionerEventHandler(eventService, bindingService, instanceService, processor)

	beta := createInstance("beta")

	handler.OnAddDatabaseInstance(beta)
	handler.OnAddDatabaseInstance(beta)
	instances := instanceService.FindAllDatabaseInstances()
	assert.Equal(t, beta.Name, instances[0].K8sName)
	assert.Equal(t, string(persistence.CREATE), string(instances[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(instances[0].Meta.Current.State))
}

func TestEventHandler_fast_subsequent_instance_add_and_removes(t *testing.T) {

	dbName := "/tmp/TestNewGoSimpleDatabaseProvisionerEventHandler.db"
	db, err := storm.Open(dbName)
	assert.Nil(t, err)
	defer os.Remove(dbName)

	stormDb := persistence.NewStormPersistenceBackend(db)

	processor := &mocks.ProvisioningEventProcessor{}
	processor.On("ProcessEvents").Return()

	eventService := service.NewPersistentEventService(stormDb)
	bindingService := service.NewPersistentDatabaseBindingService(stormDb)
	instanceService := service.NewPersistentDatabaseInstanceService(stormDb)
	handler := NewGoSimpleDatabaseProvisionerEventHandler(eventService, bindingService, instanceService, processor)

	beta := createInstance("beta")
	beta2 := createInstance("beta")

	handler.OnAddDatabaseInstance(beta)
	instances := instanceService.FindAllDatabaseInstances()
	assert.Equal(t, beta.Name, instances[0].K8sName)
	assert.Equal(t, string(persistence.CREATE), string(instances[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(instances[0].Meta.Current.State))

	handler.OnDeleteDatabaseInstance(beta)
	instances = instanceService.FindAllDatabaseInstances()
	assert.Equal(t, beta.Name, instances[0].K8sName)
	assert.Equal(t, string(persistence.DELETE), string(instances[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(instances[0].Meta.Current.State))

	handler.OnAddDatabaseInstance(beta2)
	instances = instanceService.FindAllDatabaseInstances()
	assert.Equal(t, beta2.Name, instances[0].K8sName)
	assert.Equal(t, string(persistence.CREATE), string(instances[0].Meta.Current.Action))
	assert.Equal(t, string(persistence.PENDING), string(instances[0].Meta.Current.State))

}
