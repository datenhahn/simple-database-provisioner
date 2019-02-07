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
	"github.com/sirupsen/logrus"
	"simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	"simple-database-provisioner/src/db"
	"simple-database-provisioner/src/service"
)

type SimpleDatabaseProvisionerEventHandler interface {
	OnAddDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding)
	OnDeleteDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding)
	OnAddDatabaseInstance(binding *v1alpha1.SimpleDatabaseInstance)
	OnDeleteDatabaseInstance(binding *v1alpha1.SimpleDatabaseInstance)
}

type GoSimpleDatabaseProvisionerEventHandler struct {
	crdService service.CustomResourceService
	processor  ProvisioningEventProcessor
}

func createEventId(action, objectUid string) string {
	return fmt.Sprintf("%s-%s", action, objectUid)
}

func NewGoSimpleDatabaseProvisionerEventHandler(crdService service.CustomResourceService, processor ProvisioningEventProcessor) SimpleDatabaseProvisionerEventHandler {

	this := &GoSimpleDatabaseProvisionerEventHandler{}
	this.crdService = crdService
	this.processor = processor

	return this
}

func (this *GoSimpleDatabaseProvisionerEventHandler) OnAddDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) {

	eventId := createEventId("ADD", string(binding.UID))

	logrus.Infof("Received AddDatabaseBinding event '%s': %s in namespace=%s", eventId, binding.Name, binding.Namespace)

	if this.crdService.WasProcessed(eventId) {
		logrus.Infof("Event '%s' - '%s' was already processed, skipping", eventId, binding.Name)
		return
	}

	err := this.crdService.CreateDatabaseBinding(binding)

	if err != nil {
		logrus.Errorf("Could not create database binding: %v", err)
	}

	this.crdService.MarkProcessed(eventId)

	go this.processor.ProcessEvents()
}

func (this *GoSimpleDatabaseProvisionerEventHandler) OnDeleteDatabaseBinding(binding *v1alpha1.SimpleDatabaseBinding) {

	eventId := createEventId("DELETE", string(binding.UID))

	logrus.Infof("Received MarkDatabaseBindingForDeletion event '%s': %s in namespace=%s", eventId, binding.Name, binding.Namespace)

	if this.crdService.WasProcessed(eventId) {
		logrus.Infof("Event '%s' - '%s' - '%s' was already processed, skipping", eventId, "DELETE", binding.Name)
		return
	}

	err := this.crdService.MarkDatabaseBindingForDeletion(db.NewNamespaceUniqueId(binding.Namespace, binding.Name))

	if err != nil {
		logrus.Errorf("Could not delete database binding '%s': %v", binding.Name, err)
	}

	this.crdService.MarkProcessed(eventId)

	go this.processor.ProcessEvents()
}

func (this *GoSimpleDatabaseProvisionerEventHandler) OnAddDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) {

	eventId := createEventId("ADD", string(instance.UID))

	logrus.Infof("Received AddDatabaseInstance event '%s': %s in namespace=%s", eventId, instance.Name, instance.Namespace)

	if this.crdService.WasProcessed(eventId) {
		logrus.Infof("Event '%s' - '%s' was already processed, skipping", eventId, instance.Name)
		return
	}

	err := this.crdService.CreateDatabaseInstance(instance)

	if err != nil {
		logrus.Errorf("Could not create database instance '%s': %v", instance.Name, err)
	}

	this.crdService.MarkProcessed(eventId)

	go this.processor.ProcessEvents()
}

func (this *GoSimpleDatabaseProvisionerEventHandler) OnDeleteDatabaseInstance(instance *v1alpha1.SimpleDatabaseInstance) {

	eventId := createEventId("DELETE", string(instance.UID))

	logrus.Infof("Received MarkDatabaseInstanceForDeletion event '%s': %s in namespace=%s", eventId, instance.Name, instance.Namespace)

	if this.crdService.WasProcessed(eventId) {
		logrus.Infof("Event '%s' - '%s' was already processed, skipping", eventId, instance.Name)
		return
	}

	err := this.crdService.MarkDatabaseInstanceForDeletion(db.NewNamespaceUniqueId(instance.Namespace, instance.Name))

	if err != nil {
		logrus.Errorf("Could not delete database instance '%s': %v", instance.Name, err)
	}

	this.crdService.MarkProcessed(eventId)

	go this.processor.ProcessEvents()
}
