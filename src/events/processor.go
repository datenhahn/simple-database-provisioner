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
	"strings"
	"time"
)

type ProvisioningEventProcessor interface {
	ProcessEvents()
}

type PollingEventProcessor struct {
	pollInterval  time.Duration
	appConfig     config.AppConfig
	crdService    service.CustomResourceService
	dbmsProviders []dbms.DbmsProvider
	apiclient     k8sclient.K8sClient
}

func NewPollingEventProcessor(pollInterval time.Duration,
	appConfig config.AppConfig, crdService service.CustomResourceService, apiclient k8sclient.K8sClient, dbmsProviders []dbms.DbmsProvider) ProvisioningEventProcessor {
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

		err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.K8sName)
		}

		return
	}

	provider, err := this.getDbmsProvider(dbmsServer.Type)

	if err != nil {
		errorState := CreateErrorState(instance.Meta.Current.Action, err.Error())

		err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.K8sName)
		}
		return
	}

	if instance.Meta.Current.Action == db.CREATE {

		credentials, err := this.getDbmsCredentials(dbmsServer)

		if err != nil {

			message := fmt.Sprintf("Could not get database credentials for server '%s'", dbmsServer.Name)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.K8sName, err)
			}

			return
		}

		exists, err := provider.ExistsDatabaseInstance(dbmsServer.Name, credentials, instance.PrefixedDatabaseName())

		if err != nil {
			message := fmt.Sprintf("Could not check if database exists for server '%s' - %v", dbmsServer.Name, err)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.K8sName, err)
			}

			return
		}

		if exists {
			newState := CreateOkState(instance.Meta.Current.Action)
			newState.Message = "Database already existed, keeping existing db"
			err = this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), newState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", instance.K8sName)
			}

			logrus.Infof("Database already exists instance: %s", instance.K8sName)
			return
		}

		instanceCredentials, err := provider.CreateDatabaseInstance(dbmsServer.Name, credentials, instance.PrefixedDatabaseName())

		if err != nil {
			message := fmt.Sprintf("Could not create database instance '%s' %v", instance.K8sName, err)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.K8sName, err)
			}
			return
		}

		secretData, err := instanceCredentials.ToSecretData()

		if err != nil {
			message := fmt.Sprintf("Could not create database instance secret from credentials '%s' %v", instance.K8sName, err)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.K8sName, err)
			}
			return
		}

		err = this.crdService.UpdateDatabaseInstanceCredentials(instance.NamespaceUniqueId(), secretData)

		if err != nil {
			message := fmt.Sprintf("Could not update database credentials for instance '%s'", instance.K8sName)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", instance.K8sName)
			}
			return
		}

		err = this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), CreateOkState(instance.Meta.Current.Action))
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.K8sName)
		}

		logrus.Infof("Successfully created instance: namespace=%s, instance=%s", instance.Namespace, instance.K8sName)

	} else if instance.Meta.Current.Action == db.DELETE {

		credentials, err := this.getDbmsCredentials(dbmsServer)

		if err != nil {

			message := fmt.Sprintf("Could not get database credentials for server '%s'", dbmsServer.Name)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.K8sName, err)
			}

			return
		}

		err = provider.DeleteDatabaseInstance(dbmsServer.Name, credentials, instance.PrefixedDatabaseName())

		//TODO: replace string checks with err type checks
		if err != nil && !strings.Contains(err.Error(), "does not exist") {

			message := fmt.Sprintf("Could not delete databaseInstance '%s': %v", instance.K8sName, err)

			errorState := CreateErrorState(instance.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s, %v", instance.NamespaceUniqueId(), err)
			}

			return
		}

		err = this.crdService.DeleteDatabaseInstance(instance.NamespaceUniqueId())

		if err != nil {
			logrus.Errorf("There was an error deleting instance: %s - %v", instance.NamespaceUniqueId(), err)
		}

		logrus.Infof("Successfully deleted instance: %s", instance.NamespaceUniqueId())

	} else {

		message := fmt.Sprintf("Could not handle action '%s' for databaseInstance '%s'", instance.Meta.Current.Action, instance.NamespaceUniqueId())

		errorState := CreateErrorState(instance.Meta.Current.Action, message)

		err := this.crdService.UpdateDatabaseInstanceState(instance.NamespaceUniqueId(), errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", instance.NamespaceUniqueId())
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

			err := this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
			}
			return
		}

		if dbInstance.Meta.Current.State != db.READY {
			message := fmt.Sprintf("Database Instance '%s' is not ready yet", binding.DatabaseInstanceId)
			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
			}
			return
		}

		err = this.apiclient.CreateSecret(binding.Namespace, binding.SecretName, dbInstance.Credentials)

		if err != nil {

			//TODO: replace string checks with err type checks
			if strings.Contains(err.Error(), "already exists") {
				newState := CreateOkState(binding.Meta.Current.Action)
				newState.Message = "Secret already existed, using existing secret"
				err = this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), newState)
				if err != nil {
					logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
				}

				logrus.Infof("Secret already exists: %s", binding.NamespaceUniqueId())
				return
			}

			message := fmt.Sprintf("Could not create secret for binding '%s': %v", binding.NamespaceUniqueId(), err)

			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
			}
			return
		}

		err = this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), CreateOkState(binding.Meta.Current.Action))
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
		}

		logrus.Infof("Successfully created binding: namespace=%s, binding=%s", binding.Namespace, binding.NamespaceUniqueId())

	} else if binding.Meta.Current.Action == db.DELETE {

		err := this.apiclient.DeleteSecret(binding.Namespace, binding.SecretName)

		//TODO: replace string checks with err type checks
		if err != nil && !strings.Contains(err.Error(), "not found") {
			message := fmt.Sprintf("Could not delete binding: %s, %v", binding.NamespaceUniqueId(), err)

			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
			}
			return
		}

		err = this.crdService.DeleteDatabaseBinding(binding.NamespaceUniqueId())
		if err != nil {
			message := fmt.Sprintf("Could not delete binding: %s, %v", binding.NamespaceUniqueId(), err)

			errorState := CreateErrorState(binding.Meta.Current.Action, message)

			err := this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), errorState)
			if err != nil {
				logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
			}
			return
		}

		err = this.crdService.DeleteDatabaseBinding(binding.NamespaceUniqueId())

		if err != nil {
			logrus.Errorf("There was an error deleting instance: %s - %v", binding.NamespaceUniqueId(), err)
		}

		logrus.Infof("Successfully deleted binding: %s", binding.NamespaceUniqueId())

	} else {

		message := fmt.Sprintf("Could not handle action '%s' for databaseInstance '%s'", binding.Meta.Current.Action, binding.NamespaceUniqueId())

		errorState := CreateErrorState(binding.Meta.Current.Action, message)

		err := this.crdService.UpdateDatabaseBindingState(binding.NamespaceUniqueId(), errorState)
		if err != nil {
			logrus.Errorf("There was an error updating state of instance: %s", binding.NamespaceUniqueId())
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
