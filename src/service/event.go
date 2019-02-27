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
	"simple-database-provisioner/src/persistence"
)

// PersistentEventService handles event persistence

//go:generate $GOPATH/bin/mockery -name EventService
type EventService interface {
	// WasProcessed returns true if an event was already processed
	WasProcessed(eventId string) bool

	// MarkProcessed adds an eventId to the list of processed events
	MarkProcessed(eventId string)
}

// PersistentEventService saves the performed actions in a database
// as intermediate buffer.
// The events then have to be processed by another process.
type PersistentEventService struct {
	eventRepo persistence.EventRepository
}

func NewPersistentEventService(eventRepo persistence.EventRepository) EventService {
	this := &PersistentEventService{}
	this.eventRepo = eventRepo
	return this
}

func (this *PersistentEventService) WasProcessed(eventId string) bool {

	return this.eventRepo.WasProcessed(eventId)
}

func (this *PersistentEventService) MarkProcessed(eventId string) {

	this.eventRepo.MarkProcessed(eventId)
}
