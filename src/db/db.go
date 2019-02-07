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

package db

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"simple-database-provisioner/src/util"
	"sync"
	"time"
)

type ProvisioningState string
type ProvisioningAction string

const (
	CREATE ProvisioningAction = "create"
	DELETE ProvisioningAction = "delete"
)

func (this ProvisioningAction) String() string {
	return string(this)
}

const (
	PENDING ProvisioningState = "pending"
	READY   ProvisioningState = "ready"
	ERROR   ProvisioningState = "error"
)

func (this ProvisioningState) String() string {
	return string(this)
}

type State struct {
	Action     ProvisioningAction
	State      ProvisioningState
	Message    string
	LastUpdate time.Time
}

func (this State) String() string {
	return fmt.Sprintf("{ action: '%s', state: '%s', message: '%s', lastUpdate. '%s' }",
		this.Action, this.State, this.Message, this.LastUpdate)
}

type Meta struct {
	Previous State
	Current  State
}

type NamespaceUniqueId string

func NewNamespaceUniqueId(namespace, k8sName string) NamespaceUniqueId {
	return NamespaceUniqueId(fmt.Sprintf("%s-%s", namespace, k8sName))
}

type DatabaseInstance struct {
	K8sName      string
	DbmsServer   string
	DatabaseName string
	Namespace    string
	Credentials  map[string][]byte
	Meta         Meta
}

func (this DatabaseInstance) PrefixedDatabaseName() string {

	fullName := fmt.Sprintf("%s-%s", this.Namespace, this.DatabaseName)

	sliceEnd := 54

	if len(fullName) < 54 {
		sliceEnd = len(fullName)
	}

	safeName := fmt.Sprintf("%s-%s", fullName[:sliceEnd], util.Md5Short(fullName))

	return safeName
}

func (this DatabaseInstance) NamespaceUniqueId() NamespaceUniqueId {
	return NamespaceUniqueId(fmt.Sprintf("%s-%s", this.Namespace, this.K8sName))
}

type DatabaseBinding struct {
	K8sName            string
	DatabaseInstanceId NamespaceUniqueId
	SecretName         string
	Namespace          string
	Meta               Meta
}

func (this DatabaseBinding) NamespaceUniqueId() NamespaceUniqueId {
	return NamespaceUniqueId(fmt.Sprintf("%s-%s", this.Namespace, this.K8sName))
}

type DbData struct {
	DatabaseInstances []DatabaseInstance
	DatabaseBindings  []DatabaseBinding
	ProcessedEvents   []string
}

type AppDatabase interface {
	AddDatabaseBinding(binding DatabaseBinding) error
	UpdateDatabaseBindingState(bindingId NamespaceUniqueId, newState State) error
	UpdateDatabaseInstanceState(instanceId NamespaceUniqueId, newState State) error
	AddDatabaseInstance(instance DatabaseInstance) error
	UpdateDatabaseInstanceCredentials(instanceId NamespaceUniqueId, newCredentials map[string][]byte) error
	DeleteDatabaseBinding(bindingId NamespaceUniqueId) error
	DeleteDatabaseInstance(bindingInstance NamespaceUniqueId) error
	WasProcessed(eventId string) bool
	MarkProcessed(eventId string)
	FindDatabaseInstanceById(instanceId NamespaceUniqueId) (DatabaseInstance, error)
	FindAllDatabaseInstances() []DatabaseInstance
	FindAllDatabaseBindings() []DatabaseBinding
}

type YamlAppDatabase struct {
	yamlFile string
	mutex    *sync.Mutex
}

func NewYamlAppDatabase(yamlFile string) AppDatabase {

	this := &YamlAppDatabase{}

	this.yamlFile = yamlFile
	this.mutex = &sync.Mutex{}

	return this
}

func (this *YamlAppDatabase) WasProcessed(eventId string) bool {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		panic(err)
	}

	for _, processedUid := range data.ProcessedEvents {
		if processedUid == eventId {
			return true
		}
	}

	return false
}

func (this *YamlAppDatabase) MarkProcessed(eventId string) {

	if this.WasProcessed(eventId) {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		panic(err)
	}

	data.ProcessedEvents = append(data.ProcessedEvents, eventId)

	err = this.save(data)

	if err != nil {
		panic(err)
	}
}

func (this *YamlAppDatabase) FindAllDatabaseInstances() []DatabaseInstance {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return []DatabaseInstance{}
	}

	return data.DatabaseInstances
}

func (this *YamlAppDatabase) FindDatabaseInstanceById(instanceId NamespaceUniqueId) (DatabaseInstance, error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return DatabaseInstance{}, err
	}

	for _, instance := range data.DatabaseInstances {
		if instance.NamespaceUniqueId() == instanceId {
			return instance, nil
		}
	}

	return DatabaseInstance{}, fmt.Errorf("Could not find instance with id: %s", instanceId)
}

func (this *YamlAppDatabase) FindAllDatabaseBindings() []DatabaseBinding {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return []DatabaseBinding{}
	}

	return data.DatabaseBindings
}

func (this *YamlAppDatabase) AddDatabaseBinding(binding DatabaseBinding) error {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return err
	}

	for _, dbBinding := range data.DatabaseBindings {
		if dbBinding.NamespaceUniqueId() == binding.NamespaceUniqueId() {
			if dbBinding.Meta.Current.Action == binding.Meta.Current.Action {
				logrus.Debugf("Binding with id '%s' already exists, skipping action '%s'", dbBinding.NamespaceUniqueId(), binding.Meta.Current.Action)
				return nil
			} else {
				err := this.deleteDatabaseBindingNoLock(dbBinding.NamespaceUniqueId())
				if err != nil {
					return err
				}

				data, err = this.load()
				if err != nil {
					return err
				}

			}
		}
	}

	data.DatabaseBindings = append(data.DatabaseBindings, binding)

	err = this.save(data)

	return err
}

func (this *YamlAppDatabase) UpdateDatabaseBindingState(bindingId NamespaceUniqueId, newState State) error {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return err
	}

	for index, dbBinding := range data.DatabaseBindings {
		if dbBinding.NamespaceUniqueId() == bindingId {

			data.DatabaseBindings[index].Meta.Previous = dbBinding.Meta.Current
			data.DatabaseBindings[index].Meta.Current = newState

			logrus.Debugf("Updated state for binding with id '%s' : %s", dbBinding.NamespaceUniqueId(), newState.String())
		}
	}

	err = this.save(data)

	return err

}

func (this *YamlAppDatabase) deleteDatabaseBindingNoLock(bindingId NamespaceUniqueId) error {

	data, err := this.load()

	if err != nil {
		return err
	}

	newBindings := []DatabaseBinding{}

	for _, dbBinding := range data.DatabaseBindings {
		if dbBinding.NamespaceUniqueId() != bindingId {

			newBindings = append(newBindings, dbBinding)
		}
	}

	data.DatabaseBindings = newBindings

	err = this.save(data)

	return err

}

func (this *YamlAppDatabase) DeleteDatabaseBinding(bindingId NamespaceUniqueId) error {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.deleteDatabaseBindingNoLock(bindingId)

}

func (this *YamlAppDatabase) AddDatabaseInstance(instance DatabaseInstance) error {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return err
	}

	for _, dbInstance := range data.DatabaseInstances {
		if dbInstance.NamespaceUniqueId() == instance.NamespaceUniqueId() {

			if dbInstance.Meta.Current.Action == instance.Meta.Current.Action {

				logrus.Debugf("Binding with id '%s' already exists, skipping", dbInstance.NamespaceUniqueId())
				return nil
			} else {
				err := this.deleteDatabaseInstanceNoLock(dbInstance.NamespaceUniqueId())
				if err != nil {
					return err
				}

				data, err = this.load()
				if err != nil {
					return err
				}

			}

		}
	}

	data.DatabaseInstances = append(data.DatabaseInstances, instance)

	err = this.save(data)

	return err
}

func (this *YamlAppDatabase) UpdateDatabaseInstanceState(instanceId NamespaceUniqueId, newState State) error {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return err
	}

	for index, dbInstance := range data.DatabaseInstances {
		if dbInstance.NamespaceUniqueId() == instanceId {

			data.DatabaseInstances[index].Meta.Previous = dbInstance.Meta.Current
			data.DatabaseInstances[index].Meta.Current = newState

			logrus.Debugf("Updated state for binding with id '%s' : %s", dbInstance.NamespaceUniqueId(), newState.String())
		}
	}

	err = this.save(data)

	return err
}

func (this *YamlAppDatabase) UpdateDatabaseInstanceCredentials(instanceId NamespaceUniqueId, newCredentials map[string][]byte) error {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	data, err := this.load()

	if err != nil {
		return err
	}

	for index, dbInstance := range data.DatabaseInstances {
		if dbInstance.NamespaceUniqueId() == instanceId {

			data.DatabaseInstances[index].Credentials = newCredentials
		}
	}

	err = this.save(data)

	return err
}

func (this *YamlAppDatabase) deleteDatabaseInstanceNoLock(instanceId NamespaceUniqueId) error {

	data, err := this.load()

	if err != nil {
		return err
	}

	newInstances := []DatabaseInstance{}

	for _, dbInstance := range data.DatabaseInstances {
		if dbInstance.NamespaceUniqueId() != instanceId {

			newInstances = append(newInstances, dbInstance)
		}
	}

	data.DatabaseInstances = newInstances

	err = this.save(data)

	return err

}

func (this *YamlAppDatabase) DeleteDatabaseInstance(instanceId NamespaceUniqueId) error {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.deleteDatabaseInstanceNoLock(instanceId)
}

func (this *YamlAppDatabase) load() (DbData, error) {

	dbData := DbData{}

	if _, err := os.Stat(this.yamlFile); os.IsNotExist(err) {
		logrus.Infof("Database '%s' not found, creating the file", this.yamlFile)

		err := this.save(dbData)

		if err != nil {
			return dbData, err
		}
	}

	bytes, err := ioutil.ReadFile(this.yamlFile)

	if err != nil {
		logrus.Debug(err)
		return dbData, err
	}

	err = yaml.Unmarshal(bytes, &dbData)

	if err != nil {
		logrus.Debug(err)
		return dbData, err
	}

	return dbData, nil
}

func (this *YamlAppDatabase) save(data DbData) error {

	bytes, err := yaml.Marshal(data)

	if err != nil {
		logrus.Debug(err)
		return err
	}

	err = ioutil.WriteFile(this.yamlFile, bytes, 0660)

	if err != nil {
		logrus.Debug(err)
		return err
	}

	return nil
}
