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

package restapi

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"simple-database-provisioner/src/persistence"
	"simple-database-provisioner/src/service"
)

type CommandApi interface {
	RunServer(htmlPath string)
}

type RestCommandApi struct {
	bindingService  service.DatabaseBindingService
	instanceService service.DatabaseInstanceService
}

func NewRestCommandApi(bindingService service.DatabaseBindingService, instanceService service.DatabaseInstanceService) CommandApi {

	this := &RestCommandApi{}

	this.bindingService = bindingService
	this.instanceService = instanceService

	return this
}

func (this *RestCommandApi) RunServer(htmlPath string) {
	go this.runServer(htmlPath)
}

func displayBindings(bindings []persistence.DatabaseBinding) []map[string]string {
	lines := []map[string]string{}

	for _, binding := range bindings {

		text := make(map[string]string)

		text["id"] = string(binding.NamespaceUniqueId())
		text["namespace"] = binding.Namespace
		text["secret"] = binding.SecretName
		text["databaseId"] = string(binding.DatabaseInstanceId)
		text["action"] = string(binding.Meta.Current.Action)
		text["status"] = string(binding.Meta.Current.State)
		text["message"] = binding.Meta.Current.Message

		lines = append(lines, text)
	}

	return lines
}

func displayInstances(instances []persistence.DatabaseInstance) []map[string]string {
	lines := []map[string]string{}

	for _, instance := range instances {

		text := make(map[string]string)

		text["id"] = string(instance.NamespaceUniqueId())
		text["namespace"] = instance.Namespace
		text["databaseName"] = instance.DatabaseName
		text["dbmsServer"] = string(instance.DbmsServer)
		text["action"] = string(instance.Meta.Current.Action)
		text["status"] = string(instance.Meta.Current.State)
		text["message"] = instance.Meta.Current.Message

		lines = append(lines, text)
	}

	return lines
}

func (this *RestCommandApi) runServer(htmlPath string) {
	r := gin.New()
	r.Use(cors.Default())
	r.Use(static.Serve("/", static.LocalFile(htmlPath, false)))
	r.GET("/list", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"instances": displayInstances(this.instanceService.FindAllDatabaseInstances()),
			"bindings":  displayBindings(this.bindingService.FindAllDatabaseBindings()),
		})
	})

	r.Run()
}
