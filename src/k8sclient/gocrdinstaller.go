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
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"simple-database-provisioner/src/util"

	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func init() {
	util.PanicIfNotMinikube()
}

// GoCustomResourceDefinitionManager implements CustomResourceDefinitionManager
// interface to manage the CRDs inside the cluster.
type GoCustomResourceDefinitionManager struct {
	apiextensionClientSet *apiextension.Clientset
}

// insideClusterApiextensionClientSet constructs a clientset for the daemon
// running inside a kubernetes cluster as pod.
func insideClusterApiextensionClientSet() *apiextension.Clientset {

	apiextensionConfig, err := rest.InClusterConfig()

	if err != nil {
		logrus.Panicf("Could not build InClusterConfig: %s", err.Error())
	}

	apiextensionClientset, err := apiextension.NewForConfig(apiextensionConfig)
	if err != nil {
		logrus.Panicf("Could not build apiExtension ClientSet for InClusterConfig: %s", err.Error())
	}

	return apiextensionClientset
}

// outsideClusterApiextensionClientSet constructs a clientset for invocations
// from outside the cluster (e.g. during development). It uses ~/.kube/config
// as config file.
func outsideClusterApiextensionClientSet() (clientset *apiextension.Clientset) {

	var kubeconfig string

	if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		logrus.Panic("Could not read HOME directory, aborting config load. HOME is needed to load kubernetes config outside the cluster")
	}

	// use the current context in kubeconfig
	apiextensionConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		logrus.Panicf("Could not build clientconfig: %s", err.Error())
	}

	// create the clientset
	apiextensionClientSet, err := apiextension.NewForConfig(apiextensionConfig)
	if err != nil {
		logrus.Panicf("Could not build clientset: %s", err.Error())
	}

	return apiextensionClientSet
}

// NewGoCustomResourceDefinitionManager creates a new GoCustomResourceDefinitionManager
// instance.
//
// Set the isRunningInsideCluster parameter to true if running inside a kubernetes
// cluster pod. Set it to false if invoking from outside a cluster (e.g. during
// development).
func NewGoCustomResourceDefinitionManager(isRunningInsideCluster bool) *GoCustomResourceDefinitionManager {

	this := &GoCustomResourceDefinitionManager{}

	if isRunningInsideCluster {

		this.apiextensionClientSet = insideClusterApiextensionClientSet()

	} else {
		this.apiextensionClientSet = outsideClusterApiextensionClientSet()
	}

	return this
}

// isCustomResourceDefinitionInstalled checks if a CRD is installed in the cluster
// or not.
//
// Use the fully qualified name of the CustomResourceDefinition, e.g.:
//
//   crdManager.isCustomResourceDefinitionInstalled("simpledatabasebindings.k8s.ecodia.de")
func (this *GoCustomResourceDefinitionManager) isCustomResourceDefinitionInstalled(name string) bool {
	options := v1.GetOptions{}
	crd, err := this.apiextensionClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, options)

	if err != nil {
		logrus.Debugf("Could not find CustomResourceDefinition '%s': %s", name, err)
		return false
	} else {
		logrus.Debugf("Found CustomResourceDefinition '%s': %s", name, crd.ObjectMeta.Name)
		return true
	}
}

// deleteCustomResourceDefinition deletes a CRD from the cluster.
//
// Use the fully qualified name of the CustomResourceDefinition, e.g.:
//
//   crdManager.deleteCustomResourceDefinition("simpledatabasebindings.k8s.ecodia.de")
func (this *GoCustomResourceDefinitionManager) deleteCustomResourceDefinition(name string) error {
	util.PanicIfNotMinikube()
	options := v1.DeleteOptions{}
	err := this.apiextensionClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(name, &options)

	if err != nil {
		logrus.Debugf("Could not delete CustomResourceDefinition '%s': %s", name, err)
		return err
	} else {
		logrus.Debugf("Deleted CustomResourceDefinition '%s'", name)
		return nil
	}
}

// InstallCustomResourceDefinition installs a CustomResourceDefinition-yaml file
// into the cluster, if it does not already exist.
//
// If a CustomResourceDefinition with the same name already exits the operation
// does nothing. The existing definition is NOT updated.
func (this *GoCustomResourceDefinitionManager) InstallCustomResourceDefinition(crdYamlFile string) error {
	logrus.Infof("Install Custom Resource Definition: '%s'", crdYamlFile)

	stream, err := os.Open(crdYamlFile)

	if err != nil {
		errMessage := fmt.Sprintf("Failed to read CRD file: %s", err.Error())
		logrus.Errorf(errMessage)
		return fmt.Errorf(errMessage)
	}

	crdStruct := v1beta1.CustomResourceDefinition{}
	parseErr := yaml.NewYAMLOrJSONDecoder(stream, 4096).Decode(&crdStruct)

	if parseErr != nil {
		errMessage := fmt.Sprintf("Failed to parse CRD file: %s", parseErr.Error())
		logrus.Errorf(errMessage)
		return fmt.Errorf(errMessage)
	}

	if this.isCustomResourceDefinitionInstalled(crdStruct.ObjectMeta.Name) {
		logrus.Infof("Found CRD '%s' already installed, skipping installation", crdStruct.ObjectMeta.Name)
	} else {
		_, createErr := this.apiextensionClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&crdStruct)

		if createErr != nil {
			errMessage := fmt.Sprintf("Could not create CustomResourceDefinition '%s': %s", crdStruct.Name, createErr.Error())
			logrus.Errorf(errMessage)
			return fmt.Errorf(errMessage)
		} else {
			logrus.Infof("Successfully created custom resource definition '%s'", crdStruct.Name)
		}
	}

	return nil
}
