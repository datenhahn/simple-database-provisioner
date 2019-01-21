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

package main

import "C"
import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
	"os"
	"path"
	"simple-database-provisioner/pkg/apis/simpledatabaseprovisioner/v1alpha1"
	"simple-database-provisioner/src/config"
	"simple-database-provisioner/src/db"
	"simple-database-provisioner/src/dbms"
	"simple-database-provisioner/src/dbms/provider"
	"simple-database-provisioner/src/events"
	"simple-database-provisioner/src/k8sclient"
	"simple-database-provisioner/src/restapi"
	"simple-database-provisioner/src/service"
	"time"
)

func printInfo() {

	banner := `
       _            __         __     __       __                                    _     _                  
  ___ (_)_ _  ___  / /__   ___/ /__ _/ /____ _/ /  ___ ____ ___   ___  _______ _  __(_)__ (_)__  ___  ___ ____
 (_-</ /  ' \/ _ \/ / -_) / _  / _ '/ __/ _ '/ _ \/ _ '(_-</ -_) / _ \/ __/ _ \ |/ / (_-</ / _ \/ _ \/ -_) __/
/___/_/_/_/_/ .__/_/\__/  \_,_/\_,_/\__/\_,_/_.__/\_,_/___/\__/ / .__/_/  \___/___/_/___/_/\___/_//_/\__/_/
           /_/                                                 /_/
  Copyright (c) 2019 Ecodia GmbH & Co. KG ( https://ecodia.de ) <opensource@ecodia.de>
  Licensed under the Apache License, Version 2.0

`
	fmt.Print(banner)

	// this sleep is needed to order output correctly, otherwise logrus outputs its first line first
	time.Sleep(100 * time.Millisecond)

	logrus.Info("Starting simple database provisioner controller ...")
}

func installCustomResourceDefinitions(isRunningInsideCluster bool, crdPath string) {
	crdManager := k8sclient.NewGoCustomResourceDefinitionManager(isRunningInsideCluster)

	err := crdManager.InstallCustomResourceDefinition(path.Join(crdPath, "simpledatabasebinding.yaml"))

	if err != nil {
		panic(err)
	}

	err = crdManager.InstallCustomResourceDefinition(path.Join(crdPath, "simpledatabaseinstance.yaml"))

	if err != nil {
		panic(err)
	}
}

func main() {

	printInfo()

	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Reading config file from: config.yaml")

	configFile := flag.String("configFile", "/app/config.yaml", "the configuration file")
	databaseFile := flag.String("databaseFile", "/persistence/database.yaml", "the database file where all state is stored in")
	crdPath := flag.String("crdPath", "/app/crds", "path to the simpledatabasebinding.yaml and simpledatabaseinstance.yaml custom resource definitions")
	logLevel := flag.String("logLevel", "info", "the log level (e.g. info, debug)")
	htmlPath := flag.String("htmlPath", "/app/html", "the path to the webui html directory (will be served as /)")

	flag.Parse()

	if os.Getenv("SIMPLEDATABASEPROVISIONER_CONFIGFILE") != "" {
		envConfigFile := os.Getenv("SIMPLEDATABASEPROVISIONER_CONFIGFILE")
		configFile = &envConfigFile
	}

	if os.Getenv("SIMPLEDATABASEPROVISIONER_DATABASEFILE") != "" {
		envDatabaseFile := os.Getenv("SIMPLEDATABASEPROVISIONER_DATABASEFILE")
		databaseFile = &envDatabaseFile
	}

	if os.Getenv("SIMPLEDATABASEPROVISIONER_CRDPATH") != "" {
		envCrdPath := os.Getenv("SIMPLEDATABASEPROVISIONER_CRDPATH")
		crdPath = &envCrdPath
	}

	if os.Getenv("SIMPLEDATABASEPROVISIONER_LOGLEVEL") != "" {
		envLogLevel := os.Getenv("SIMPLEDATABASEPROVISIONER_LOGLEVEL")
		logLevel = &envLogLevel
	}

	if os.Getenv("SIMPLEDATABASEPROVISIONER_HTMLPATH") != "" {
		envHtmlPath := os.Getenv("SIMPLEDATABASEPROVISIONER_HTMLPATH")
		htmlPath = &envHtmlPath
	}

	logrus.Infof("configFile: %s", *configFile)
	logrus.Infof("databaseFile: %s", *databaseFile)
	logrus.Infof("crdPath: %s", *crdPath)
	logrus.Infof("logLevel: %s", *logLevel)
	logrus.Infof("htmlPath: %s", *htmlPath)

	level, err := logrus.ParseLevel(*logLevel)

	if err != nil {
		logrus.Panicf("Could not parse loglevel: %s : %v", *logLevel, err)
	}

	logrus.SetLevel(level)

	myConfig, err := config.ReadConfig(*configFile)

	if err != nil {
		logrus.Panicf("Could not read config file: %s : %v", *configFile, err)
	}

	insideCluster := k8sclient.IsRunningInsideCluster()

	if insideCluster {
		logrus.Info("Detected running INSIDE CLUSTER")
	} else {
		logrus.Info("Detected running OUTSIDE CLUSTER")
	}

	installCustomResourceDefinitions(insideCluster, *crdPath)

	client := k8sclient.NewGoK8sClient(insideCluster)

	appDb := db.NewYamlAppDatabase(*databaseFile)
	crdService := service.NewPersistentCustomResourceService(appDb)
	crdEventHandler := events.NewGoSimpleDatabaseProvisionerEventHandler(crdService)

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

			switch v := obj.(type) {
			default:
				fmt.Printf("unexpected type %T", v)
			case *v1alpha1.SimpleDatabaseBinding:
				crdEventHandler.OnAddDatabaseBinding(v)
			case *v1alpha1.SimpleDatabaseInstance:
				crdEventHandler.OnAddDatabaseInstance(v)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			logrus.Debugf("Detected update event, but skipping: old=%v , new=%v", old, new)
		},
		DeleteFunc: func(obj interface{}) {
			switch v := obj.(type) {
			default:
				fmt.Printf("unexpected type %T", v)
			case *v1alpha1.SimpleDatabaseBinding:
				crdEventHandler.OnDeleteDatabaseBinding(v)
			case *v1alpha1.SimpleDatabaseInstance:
				crdEventHandler.OnDeleteDatabaseInstance(v)
			}
		},
	}

	postgresProvider := &provider.PostgresqlDbmsProvider{}

	pollingProcessor := events.NewPollingEventProcessor(10*time.Second, myConfig, crdService, client, []dbms.DbmsProvider{postgresProvider})

	logrus.Info("Starting polling with interval 10 secs ...")

	commandApi := restapi.NewRestCommandApi(crdService)

	commandApi.RunServer(*htmlPath)

	go func() {
		for {
			pollingProcessor.ProcessEvents()
			time.Sleep(10 * time.Second)
		}
	}()

	err = client.WatchSimpleDatabaseProvisionerCustomResources(eventHandler)

	if err != nil {
		panic(err)
	}

}
