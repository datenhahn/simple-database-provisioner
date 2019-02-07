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
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	"simple-database-provisioner/src/db"
	"simple-database-provisioner/src/events/mocks"
	"simple-database-provisioner/src/service"
	"simple-database-provisioner/src/util"
	"testing"
	"time"
)

func createBinding(name string) *v1alpha1.SimpleDatabaseBinding {
	return &v1alpha1.SimpleDatabaseBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			UID:  types.UID(util.Md5(fmt.Sprintf("%d", time.Now().UnixNano()))),
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

	dbName := "/tmp/TestNewGoSimpleDatabaseProvisionerEventHandler.yaml"

	os.Remove(dbName)
	defer os.Remove(dbName)

	yamlDb := db.NewYamlAppDatabase(dbName)

	processor := &mocks.ProvisioningEventProcessor{}
	processor.On("ProcessEvents").Return()

	crdService := service.NewPersistentCustomResourceService(yamlDb)
	handler := NewGoSimpleDatabaseProvisionerEventHandler(crdService, processor)

	alpha := createBinding("alpha")
	alpha2 := createBinding("alpha")

	handler.OnAddDatabaseBinding(alpha)
	bindings := crdService.FindAllDatabaseBindings()
	assert.Equal(t, alpha.Name, bindings[0].K8sName)
	assert.Equal(t, string(db.CREATE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(db.PENDING), string(bindings[0].Meta.Current.State))

	handler.OnDeleteDatabaseBinding(alpha)
	bindings = crdService.FindAllDatabaseBindings()
	assert.Equal(t, alpha.Name, bindings[0].K8sName)
	assert.Equal(t, string(db.DELETE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(db.PENDING), string(bindings[0].Meta.Current.State))

	handler.OnAddDatabaseBinding(alpha2)
	bindings = crdService.FindAllDatabaseBindings()
	assert.Equal(t, alpha2.Name, bindings[0].K8sName)
	assert.Equal(t, string(db.CREATE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(db.PENDING), string(bindings[0].Meta.Current.State))

	handler.OnDeleteDatabaseBinding(alpha2)
	bindings = crdService.FindAllDatabaseBindings()
	assert.Equal(t, alpha2.Name, bindings[0].K8sName)
	assert.Equal(t, string(db.DELETE), string(bindings[0].Meta.Current.Action))
	assert.Equal(t, string(db.PENDING), string(bindings[0].Meta.Current.State))

}

func TestEventHandler_fast_subsequent_instance_add_and_removes(t *testing.T) {

	dbName := "/tmp/TestNewGoSimpleDatabaseProvisionerEventHandler.yaml"

	os.Remove(dbName)
	defer os.Remove(dbName)

	yamlDb := db.NewYamlAppDatabase(dbName)

	processor := &mocks.ProvisioningEventProcessor{}
	processor.On("ProcessEvents").Return()

	crdService := service.NewPersistentCustomResourceService(yamlDb)
	handler := NewGoSimpleDatabaseProvisionerEventHandler(crdService, processor)

	beta := createInstance("beta")
	beta2 := createInstance("beta")

	handler.OnAddDatabaseInstance(beta)
	instances := crdService.FindAllDatabaseInstances()
	assert.Equal(t, beta.Name, instances[0].K8sName)
	assert.Equal(t, string(db.CREATE), string(instances[0].Meta.Current.Action))
	assert.Equal(t, string(db.PENDING), string(instances[0].Meta.Current.State))

	handler.OnDeleteDatabaseInstance(beta)
	instances = crdService.FindAllDatabaseInstances()
	assert.Equal(t, beta.Name, instances[0].K8sName)
	assert.Equal(t, string(db.DELETE), string(instances[0].Meta.Current.Action))
	assert.Equal(t, string(db.PENDING), string(instances[0].Meta.Current.State))

	handler.OnAddDatabaseInstance(beta2)
	instances = crdService.FindAllDatabaseInstances()
	assert.Equal(t, beta2.Name, instances[0].K8sName)
	assert.Equal(t, string(db.CREATE), string(instances[0].Meta.Current.Action))
	assert.Equal(t, string(db.PENDING), string(instances[0].Meta.Current.State))

}
