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

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func TestGoCustomResourceDefinitionInstaller_NoHomeVar(t *testing.T) {

	home := os.Getenv("HOME")

	os.Unsetenv("HOME")
	assert.Panics(t, func() {
		NewGoCustomResourceDefinitionManager(false)
	})

	os.Setenv("HOME", home)

}

func TestGoCustomResourceDefinitionInstaller_NoKubeConfigFound(t *testing.T) {

	home := os.Getenv("HOME")

	os.Setenv("HOME", "/tmp")
	assert.Panics(t, func() {
		NewGoCustomResourceDefinitionManager(false)
	})

	os.Setenv("HOME", home)
}

func TestGoCustomResourceDefinitionInstaller_CannotParseYaml(t *testing.T) {

	crdInstaller := NewGoCustomResourceDefinitionManager(false)

	err := crdInstaller.InstallCustomResourceDefinition("testdata/crds/invalidcrd.yaml")
	assert.Error(t, err)
}

func TestGoCustomResourceDefinitionInstaller_CannotCreateCrd(t *testing.T) {

	crdInstaller := NewGoCustomResourceDefinitionManager(false)

	err := crdInstaller.InstallCustomResourceDefinition("testdata/crds/notcreatablecrd.yaml")
	assert.Error(t, err)
}

func TestGoCustomResourceDefinitionInstaller_NonExistentYaml(t *testing.T) {

	crdInstaller := NewGoCustomResourceDefinitionManager(false)

	err := crdInstaller.InstallCustomResourceDefinition("testdata/crds/doesnotexist.yaml")
	assert.Error(t, err)
}

func TestGoCustomResourceDefinitionInstaller_InstallCustomResourceDefinition(t *testing.T) {

	crdInstaller := NewGoCustomResourceDefinitionManager(false)

	// Dangerous as it might lead to resource deletions
	//err := crdInstaller.deleteCustomResourceDefinition(SDP_BINDINGS_FQDN)
	//err = crdInstaller.deleteCustomResourceDefinition(SDP_INSTANCES_FQDN)
	//
	//assert.False(t, crdInstaller.isCustomResourceDefinitionInstalled(SDP_BINDINGS_FQDN))
	//assert.False(t, crdInstaller.isCustomResourceDefinitionInstalled(SDP_INSTANCES_FQDN))

	err := crdInstaller.InstallCustomResourceDefinition("../../crds/simpledatabasebinding.yaml")
	assert.Nil(t, err)

	err2 := crdInstaller.InstallCustomResourceDefinition("../../crds/simpledatabaseinstance.yaml")
	assert.Nil(t, err2)

	err3 := crdInstaller.InstallCustomResourceDefinition("../../crds/simpledatabaseinstance.yaml")
	assert.Nil(t, err3)

	assert.True(t, crdInstaller.isCustomResourceDefinitionInstalled(SDP_BINDINGS_FQDN))
	assert.True(t, crdInstaller.isCustomResourceDefinitionInstalled(SDP_INSTANCES_FQDN))

	// Dangerous as it might lead to resource deletions
	//err4 := crdInstaller.deleteCustomResourceDefinition(SDP_BINDINGS_FQDN)
	//assert.Nil(t, err4)
	//
	//err5 := crdInstaller.deleteCustomResourceDefinition(SDP_INSTANCES_FQDN)
	//assert.Nil(t, err5)
}

func TestGoCustomResourceDefinitionInstaller_InitInsideCluster(t *testing.T) {

	// Expect to panic running inside cluster config outside cluster

	assert.Panics(t, func() { NewGoCustomResourceDefinitionManager(true) })

}
