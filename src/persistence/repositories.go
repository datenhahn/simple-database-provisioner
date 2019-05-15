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

// DatabaseBindingRepository allows CRUD operations on a DatabaseBinding
//go:generate $GOPATH/bin/mockery -name DatabaseBindingRepository
type DatabaseBindingRepository interface {

	// AddDatabaseBinding adds a database binding
	AddDatabaseBinding(binding DatabaseBinding) error

	// UpdateDatabaseBindingState updates the state property of a database binding
	UpdateDatabaseBindingState(bindingId NamespaceUniqueId, newState State) error

	// DeleteDatabaseBinding deletes a database binding. The binding is identified
	// by the NamespaceUniqueId.
	DeleteDatabaseBinding(bindingId NamespaceUniqueId) error

	// FindAllDatabaseBindings returns a list of database bindings
	FindAllDatabaseBindings() []DatabaseBinding
}

// EventRepository allows to store and query all processed events
//go:generate $GOPATH/bin/mockery -name EventRepository
type EventRepository interface {

	// WasProcessed returns true if an eventId was found in the database
	// and thereby was already processed
	WasProcessed(eventId string) bool

	// MarkProcessed marks an eventId as processed by adding it to the
	// database
	MarkProcessed(eventId string)
}

// DatabaseInstanceRepository allows CRUD operations on a DatabaseInstance
//go:generate $GOPATH/bin/mockery -name DatabaseInstanceRepository
type DatabaseInstanceRepository interface {

	// UpdateDatabaseInstanceState updates the state property of a database instance
	UpdateDatabaseInstanceState(instanceId NamespaceUniqueId, newState State) error

	// AddDatabaseInstance adds a database instance
	AddDatabaseInstance(instance DatabaseInstance) error

	// UpdateDatabaseInstanceCredentials updates the database instance credentials
	UpdateDatabaseInstanceCredentials(instanceId NamespaceUniqueId, newCredentials map[string][]byte) error

	// DeleteDatabaseInstance deletes a database instance. The instance is identified
	// by the NamespaceUniqueId.
	DeleteDatabaseInstance(bindingInstance NamespaceUniqueId) error

	// FindDatabaseInstanceById finds a database instance by its NamespaceUniqueId
	FindDatabaseInstanceById(instanceId NamespaceUniqueId) (DatabaseInstance, error)

	// FindAllDatabaseInstances returns a list of database instances
	FindAllDatabaseInstances() []DatabaseInstance
}
