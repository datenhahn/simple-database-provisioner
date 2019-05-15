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
	"simple-database-provisioner/src/persistence/mocks"
	"testing"
)

func TestPersistentEventService_MarkProcessed(t *testing.T) {

	mockRepo := &mocks.EventRepository{}
	mockRepo.On("MarkProcessed", "some-event-id").Return()

	eventService := NewPersistentEventService(mockRepo)

	eventService.MarkProcessed("some-event-id")
}

func TestPersistentEventService_WasProcessed(t *testing.T) {

	mockRepo := &mocks.EventRepository{}
	mockRepo.On("WasProcessed", "some-event-id").Return(true)
	mockRepo.On("WasProcessed", "other-event-id").Return(false)

	eventService := NewPersistentEventService(mockRepo)

	assert.True(t, eventService.WasProcessed("some-event-id"))
	assert.False(t, eventService.WasProcessed("other-event-id"))
}
