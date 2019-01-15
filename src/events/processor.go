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
	"simple-database-provisioner/src/config"
	"simple-database-provisioner/src/db"
	"simple-database-provisioner/src/dbms"
	"simple-database-provisioner/src/k8sclient"
	"simple-database-provisioner/src/service"
	"time"
)

type ProvisioningEventProcessor interface {
	ProcessEvents()
}

type PollingEventProcessor struct {
	pollInterval  time.Duration
	appConfig     config.AppConfig
	crdService    service.CustomResourceDefinitionService
	dbmsProviders []dbms.DbmsProvider
	apiclient     k8sclient.K8sClient
}

func NewPollingEventProcessor(pollInterval time.Duration,
	appConfig config.AppConfig, crdService service.CustomResourceDefinitionService, apiclient k8sclient.K8sClient, dbmsProviders []dbms.DbmsProvider) ProvisioningEventProcessor {
	this := &PollingEventProcessor{}

	this.pollInterval = pollInterval
	this.appConfig = appConfig
	this.crdService = crdService
	this.dbmsProviders = dbmsProviders
	this.apiclient = apiclient

	return this
}

func (this *PollingEventProcessor) getDbmsProvider(dbmsType string) (dbms.DbmsProvider, error) {
	for _, provider := range this.dbmsProviders {

		providerType := provider.Type()
		if providerType == dbmsType {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("Could not find dbmsProvider for dbms type: %s", dbmsType)
}

func (this *PollingEventProcessor) getDbmsCredentials(dbmsServer config.DbmsServerConfig) (dbms.DatabaseCredentials, error) {

	secret, err := this.apiclient.ReadSecret(dbmsServer.FromSecret.Namespace, dbmsServer.FromSecret.Secret)

	if err != nil {
		return dbms.DatabaseCredentials{}, err
	}

	return dbms.CreateCredentialsFromSecretData(secret)
}

func CreateErrorState(action db.ProvisioningAction, message string) db.State {

	logrus.Errorf("Error: %s", message)
	return db.State{
		State:      db.ERROR,
		Action:     action,
		Message:    message,
		LastUpdate: time.Now().Round(time.Second),
	}
}

func CreateOkState(action db.ProvisioningAction) db.State {
	return db.State{
		State:      db.READY,
		Action:     action,
		Message:    "ok",
		LastUpdate: time.Now().Round(time.Second),
	}
}

func (this *PollingEventProcessor) processInstance(instance db.DatabaseInstance) {

	dbmsServer, err := config.GetDbmsServer(this.appConfig, instance.DbmsServer)

	if err != nil {

		errorState := CreateErrorState(instance.Meta.Current.Action, err.Error())

		err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.Id)
		}

		return
	}

	provider, err := this.getDbmsProvider(dbmsServer.Type)

	if err != nil {
		errorState := CreateErrorState(instance.Meta.Current.Action, err.Error())

		err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.Id)
		}
		return
	}

	if instance.Meta.Current.Action == db.CREATE {

		credentials, err := this.getDbmsCredentials(dbmsServer)

		if err != nil {

			message := fmt.Sprintf("Could not get database credentials for server '%s'", dbmsServer.Name)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.Id, err)
			}

			return
		}

		instanceCredentials, err := provider.CreateDatabaseInstance(dbmsServer.Name, credentials, instance.DatabaseName)

		if err != nil {
			message := fmt.Sprintf("Could not create database instance '%s' %v", instance.Id, err)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.Id, err)
			}
			return
		}

		secretData, err := instanceCredentials.ToSecretData()

		if err != nil {
			message := fmt.Sprintf("Could not create database instance secret from credentials '%s' %v", instance.Id, err)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.Id, err)
			}
			return
		}

		err = this.crdService.UpdateDatabaseInstanceCredentials(instance.Id, secretData)

		if err != nil {
			message := fmt.Sprintf("Could not update database credentials for instance '%s'", instance.Id)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", instance.Id)
			}
			return
		}

		err = this.crdService.UpdateDatabaseInstanceState(instance.Id, CreateOkState(instance.Meta.Current.Action))
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.Id)
		}

		logrus.Infof("Successfully created instance: %s", instance.Id)

	} else if instance.Meta.Current.Action == db.DELETE {

		credentials, err := this.getDbmsCredentials(dbmsServer)

		if err != nil {

			message := fmt.Sprintf("Could not get database credentials for server '%s'", dbmsServer.Name)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.Id, err)
			}

			return
		}

		err = provider.DeleteDatabaseInstance(dbmsServer.Name, credentials, instance.DatabaseName)

		if err != nil {

			message := fmt.Sprintf("Could not delete databaseInstance '%s': %v", instance.Id, err)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.Id, err)
			}

			return
		}

		err = this.crdService.DeleteDatabaseInstance(instance.Id)

		if err != nil {
			logrus.Errorf("There was an error deleting instance: %s - %v", instance.Id, err)
		}

		logrus.Infof("Successfully deleted instance: %s", instance.Id)

	} else {

		message := fmt.Sprintf("Could not handle action '%s' for databaseInstance '%s'", instance.Meta.Current.Action, instance.Id)

		errorState := CreateErrorState(instance.Meta.Current.Action, message)

		err := this.crdService.UpdateDatabaseInstanceState(instance.Id, errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.Id)
		}

		return
	}

}

func (this *PollingEventProcessor) processBinding(binding db.DatabaseBinding) {

	if binding.Meta.Current.Action == db.CREATE {

		dbInstance, err := this.crdService.FindDatabaseInstanceById(binding.DatabaseInstanceId)

		if err != nil {

			message := fmt.Sprintf("Could not find database Instance with ID: %s", binding.DatabaseInstanceId)

			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.Id)
			}
			return
		}

		if dbInstance.Meta.Current.State != db.READY {
			message := fmt.Sprintf("Database Instance '%s' is not ready yet", binding.DatabaseInstanceId)
			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.Id)
			}
			return
		}

		err = this.apiclient.CreateSecret(binding.Namespace, binding.SecretName, dbInstance.Credentials)

		if err != nil {
			message := fmt.Sprintf("Could not create secret for binding '%s': %v", binding.Id, err)

			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.Id)
			}
			return
		}

		err = this.crdService.UpdateDatabaseBindingState(binding.Id, CreateOkState(binding.Meta.Current.Action))
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", binding.Id)
		}

		logrus.Infof("Successfully created binding: %s", binding.Id)

	} else if binding.Meta.Current.Action == db.DELETE {

		err := this.apiclient.DeleteSecret(binding.Namespace, binding.SecretName)

		if err != nil {
			message := fmt.Sprintf("Could not delete binding: %s, %v", binding.Id, err)

			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.Id)
			}
			return
		}

		err = this.crdService.DeleteDatabaseBinding(binding.Id)
		if err != nil {
			message := fmt.Sprintf("Could not delete binding: %s, %v", binding.Id, err)

			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.Id, errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.Id)
			}
			return
		}

		err = this.crdService.DeleteDatabaseBinding(binding.Id)

		if err != nil {
			logrus.Errorf("There was an error deleting instance: %s - %v", binding.Id, err)
		}

		logrus.Infof("Successfully deleted binding: %s", binding.Id)

	} else {

		message := fmt.Sprintf("Could not handle action '%s' for databaseInstance '%s'", binding.Meta.Current.Action, binding.Id)

		errorState := CreateErrorState(binding.Meta.Current.Action, message)

		err := this.crdService.UpdateDatabaseBindingState(binding.Id, errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", binding.Id)
		}
		return
	}

}

func (this *PollingEventProcessor) ProcessEvents() {

	// first try to reprocess error instances
	errorInstances := this.crdService.FindInstancesByState(db.ERROR)

	for _, errorInstance := range errorInstances {
		this.processInstance(errorInstance)
	}

	// then try to reprocess error bindings
	errorBindings := this.crdService.FindBindingsByState(db.ERROR)
	for _, errorBinding := range errorBindings {
		this.processBinding(errorBinding)
	}

	// then process pending instances
	pendingInstances := this.crdService.FindInstancesByState(db.PENDING)
	for _, pendingInstance := range pendingInstances {
		this.processInstance(pendingInstance)

	}

	// then process pending bindings
	pendingBindings := this.crdService.FindBindingsByState(db.PENDING)
	for _, pendingBinding := range pendingBindings {
		this.processBinding(pendingBinding)
	}
}
