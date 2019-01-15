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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	crdclientset "simple-database-provisioner/pkg/client/clientset/versioned"
	crdinformers "simple-database-provisioner/pkg/client/informers/externalversions"
	"time"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

type GoK8sClient struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

func NewGoK8sClient(isRunningInsideCluster bool) *GoK8sClient {

	this := &GoK8sClient{}

	if isRunningInsideCluster {
		this.clientset = this.insideClusterClientset()
	} else {
		this.clientset = this.outsideClusterClientset()
	}

	return this

}

func (this *GoK8sClient) insideClusterClientset() *kubernetes.Clientset {

	var err error

	this.config, err = rest.InClusterConfig()

	if err != nil {
		logrus.Panic("Could not build InClusterConfig: {}", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(this.config)
	if err != nil {
		logrus.Panic("Could not build clientset for InClusterConfig: {}", err.Error())
	}

	return clientset
}

func (this *GoK8sClient) outsideClusterClientset() *kubernetes.Clientset {

	var kubeconfig string

	if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		logrus.Panic("Could not read HOME directory, aborting config load. HOME is needed to load kubernetes config outside the cluster")
	}

	var err error
	// use the current context in kubeconfig
	this.config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		logrus.Panic("Could not build clientconfig: {}", err.Error())
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(this.config)
	if err != nil {
		logrus.Panic("Could not build clientset: {}", err.Error())
	}

	return clientSet
}

func (this *GoK8sClient) WatchSimpleDatabaseProvisionerCustomResources(eventHandler cache.ResourceEventHandler) error {

	customclientset, err := crdclientset.NewForConfig(this.config)

	if err != nil {
		panic(err)
	}

	informerFactory := crdinformers.NewSharedInformerFactory(customclientset, time.Second*30)
	informerFactory.Simpledatabaseprovisioner().V1alpha1().SimpleDatabaseBindings().Informer().AddEventHandler(eventHandler)
	informerFactory.Simpledatabaseprovisioner().V1alpha1().SimpleDatabaseInstances().Informer().AddEventHandler(eventHandler)

	stop := make(chan struct{})
	defer close(stop)

	informerFactory.Start(stop)

	for {
		time.Sleep(time.Second)
	}

	return nil
}

func (this *GoK8sClient) ReadSecret(namespace string, name string) (map[string][]byte, error) {
	secret, err := this.clientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})

	if err != nil {

		return make(map[string][]byte), err
	}

	return secret.Data, nil
}

func (this *GoK8sClient) CreateSecret(namespace string, name string, data map[string][]byte) error {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	_, err := this.clientset.CoreV1().Secrets(namespace).Create(&secret)

	return err
}

func (this *GoK8sClient) DeleteSecret(namespace string, name string) error {

	return this.clientset.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})

}
