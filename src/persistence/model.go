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
	"simple-database-provisioner/src/util"
	"time"
)

type DatabaseInstance struct {
	NamespaceUniqueId NamespaceUniqueId `storm:"id"`
	K8sName           string
	DbmsServer        string
	DatabaseName      string
	Namespace         string
	Credentials       map[string][]byte
	Meta              Meta
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

func (this DatabaseInstance) GetNamespaceUniqueId() NamespaceUniqueId {
	return NamespaceUniqueId(fmt.Sprintf("%s-%s", this.Namespace, this.K8sName))
}

type Event struct {
	Id string `storm:"id"`
}

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
	Previous State `storm:"inline"`
	Current  State `storm:"inline"`
}

type NamespaceUniqueId string

func NewNamespaceUniqueId(namespace, k8sName string) NamespaceUniqueId {
	return NamespaceUniqueId(fmt.Sprintf("%s-%s", namespace, k8sName))
}

type DatabaseBinding struct {
	NamespaceUniqueId  NamespaceUniqueId `storm:"id"`
	K8sName            string
	DatabaseInstanceId NamespaceUniqueId
	SecretName         string
	Namespace          string
	Meta               Meta `storm:"inline"`
}

func (this DatabaseBinding) GetNamespaceUniqueId() NamespaceUniqueId {
	return NamespaceUniqueId(fmt.Sprintf("%s-%s", this.Namespace, this.K8sName))
}
