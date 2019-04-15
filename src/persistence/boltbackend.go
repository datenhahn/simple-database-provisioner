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

package persistence

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/sirupsen/logrus"
)

type AppDb interface {
	EventRepository
	DatabaseBindingRepository
	DatabaseInstanceRepository
}

type StormPersistenceBackend struct {
	db *storm.DB
}

func NewStormPersistenceBackend(db *storm.DB) *StormPersistenceBackend {

	this := &StormPersistenceBackend{}

	this.db = db
	err := db.Init(&DatabaseInstance{})
	if err != nil {
		logrus.Panicf("Error initializing bucket DatabaseInstance: %v", err)
	}
	err = db.Init(&DatabaseBinding{})
	if err != nil {
		logrus.Panicf("Error initializing bucket DatabaseBinding: %v", err)
	}
	err = db.Init(&Event{})
	if err != nil {
		logrus.Panicf("Error initializing bucket Event: %v", err)
	}

	return this
}

func (this *StormPersistenceBackend) FindAllDatabaseInstances() []DatabaseInstance {

	var instance []DatabaseInstance

	err := this.db.All(&instance)

	if err != nil {
		logrus.Panicf("Error in FindAllDatabaseInstances: %v", err)
	}

	return instance
}

func (this *StormPersistenceBackend) FindDatabaseInstanceById(instanceId NamespaceUniqueId) (DatabaseInstance, error) {

	var instance DatabaseInstance

	err := this.db.One("NamespaceUniqueId", instanceId, &instance)

	if err != nil {
		return DatabaseInstance{}, fmt.Errorf("Could not find instance with id: %s", instanceId)
	}

	return instance, nil
}

func (this *StormPersistenceBackend) FindAllDatabaseBindings() []DatabaseBinding {

	var binding []DatabaseBinding

	err := this.db.All(&binding)

	if err != nil {
		logrus.Panicf("Error in FindAllDatabaseBindings: %v", err)
	}

	return binding

}

func (this *StormPersistenceBackend) AddDatabaseBinding(binding DatabaseBinding) error {

	err := this.db.Save(&binding)

	if err != nil {
		if err == storm.ErrAlreadyExists {

			var existing DatabaseBinding
			err := this.db.One("NamespaceUniqueId", binding.NamespaceUniqueId, &existing)

			if err != nil {
				return fmt.Errorf("Error during adding database binding (lookup of existing binding), %v", err)
			}

			if existing.Meta.Current.Action == binding.Meta.Current.Action {
				logrus.Debugf("Binding with id '%s' already exists, skipping action '%s'", binding.NamespaceUniqueId, binding.Meta.Current.Action)
				return nil
			} else {

				err := this.db.Delete("DatabaseBinding", &binding.NamespaceUniqueId)
				if err != nil {
					return fmt.Errorf("Error deleting existing, outdated database binding during add: %v", err)
				}

				err = this.db.Save(&binding)

				if err != nil {
					return fmt.Errorf("Error during add of new Database Binding: %v", err)
				}
			}
		} else {
			return fmt.Errorf("Error during add of new Database Binding: %v", err)
		}
	}

	return nil
}

func (this *StormPersistenceBackend) WasProcessed(eventId string) bool {

	event := Event{}

	err := this.db.One("Id", eventId, &event)

	if err != nil {
		if err.Error() != "not found" {
			logrus.Panicf("Unexpected error in checking WasProcessed: %v", err)
		}

		return false
	} else {
		return true
	}
}

func (this *StormPersistenceBackend) MarkProcessed(eventId string) {

	err := this.db.Save(&Event{Id: eventId})

	if err != nil {
		logrus.Panicf("Error saving event: %v", err)
	}
}

func (this *StormPersistenceBackend) UpdateDatabaseBindingState(bindingId NamespaceUniqueId, newState State) error {

	var existing DatabaseBinding
	err := this.db.One("NamespaceUniqueId", bindingId, &existing)

	if err != nil {
		return fmt.Errorf("Error during update of database binding %s : %v", bindingId, err)
	}

	existing.Meta.Previous = existing.Meta.Current
	existing.Meta.Current = newState

	err = this.db.Update(&existing)

	if err != nil {
		return fmt.Errorf("Error during update of database binding %s : %v", bindingId, err)
	}

	logrus.Debugf("Updated state for binding with id '%s' : %s", existing.NamespaceUniqueId, newState.String())

	return nil
}

func (this *StormPersistenceBackend) DeleteDatabaseBinding(bindingId NamespaceUniqueId) error {

	err := this.db.Delete("DatabaseBinding", &bindingId)

	if err != nil {
		return fmt.Errorf("Error during deletion of binding %s : %v", bindingId, err)
	}

	return nil
}

func (this *StormPersistenceBackend) AddDatabaseInstance(instance DatabaseInstance) error {

	err := this.db.Save(&instance)

	if err == storm.ErrAlreadyExists {

		var existing DatabaseInstance
		err := this.db.One("NamespaceUniqueId", instance.NamespaceUniqueId, &existing)

		if err != nil {
			return fmt.Errorf("Error during adding database instance (lookup of existing instance), %v", err)
		}

		if existing.Meta.Current.Action == instance.Meta.Current.Action {
			logrus.Debugf("Instance with id '%s' already exists, skipping action '%s'", instance.NamespaceUniqueId, instance.Meta.Current.Action)
			return nil
		} else {

			err := this.db.Delete("DatabaseInstance", &instance.NamespaceUniqueId)
			if err != nil {
				return fmt.Errorf("Error deleting existing, outdated database instance during add: %v", err)
			}

			err = this.db.Save(&instance)

			if err != nil {
				return fmt.Errorf("Error during add of new Database Instance: %v", err)
			}
		}
	} else {
		return fmt.Errorf("Error during add of new Database Instance: %v", err)
	}

	return err
}

func (this *StormPersistenceBackend) UpdateDatabaseInstanceState(instanceId NamespaceUniqueId, newState State) error {

	var existing DatabaseInstance
	err := this.db.One("NamespaceUniqueId", instanceId, &existing)

	if err != nil {
		return fmt.Errorf("Error during update of database instance %s : %v", instanceId, err)
	}

	existing.Meta.Previous = existing.Meta.Current
	existing.Meta.Current = newState

	err = this.db.Update(&existing)

	if err != nil {
		return fmt.Errorf("Error during update of database instance %s : %v", instanceId, err)
	}

	logrus.Debugf("Updated state for instance with id '%s' : %s", existing.NamespaceUniqueId, newState.String())

	return nil
}

func (this *StormPersistenceBackend) UpdateDatabaseInstanceCredentials(instanceId NamespaceUniqueId, newCredentials map[string][]byte) error {

	var existing DatabaseInstance
	err := this.db.One("NamespaceUniqueId", instanceId, &existing)

	if err != nil {
		return fmt.Errorf("Error during update of database instance credentials %s : %v", instanceId, err)
	}

	err = this.db.UpdateField(&existing, "credentials", newCredentials)

	if err != nil {
		return fmt.Errorf("Error during update of database instance credentials %s : %v", instanceId, err)
	}

	logrus.Debugf("Updated credentials for instance with id '%s'", existing.NamespaceUniqueId)

	return nil

}

func (this *StormPersistenceBackend) DeleteDatabaseInstance(instanceId NamespaceUniqueId) error {

	err := this.db.Delete("DatabaseInstance", &instanceId)

	if err != nil {
		return fmt.Errorf("Error during deletion of instance %s : %v", instanceId, err)
	}

	return nil
}
