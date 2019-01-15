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

package k8sclient

//go:generate $GOPATH/bin/mockery -name CustomResourceDefinitionManager

// CustomResourceDefinitionManager interface provides methods to manage
// CustomResourceDefinitions inside the cluster.
type CustomResourceDefinitionManager interface {

	// InstallCustomResourceDefinition installs a CustomResourceDefinition-yaml file
	// into the cluster, if it does not already exist.
	//
	// If a CustomResourceDefinition with the same name already exits the operation
	// does nothing. The existing definition is NOT updated.
	InstallCustomResourceDefinition(crdYamlFilePath string) error
}

const SDP_BINDINGS_FQDN = "simpledatabasebindings.simpledatabaseprovisioner.k8s.ecodia.de"
const SDP_INSTANCES_FQDN = "simpledatabaseinstances.simpledatabaseprovisioner.k8s.ecodia.de"
