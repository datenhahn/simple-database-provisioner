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
	"github.com/gin-gonic/gin"
	"simple-database-provisioner/src/db"
	"simple-database-provisioner/src/service"
)

type CommandApi interface {
	RunServer()
}

type RestCommandApi struct {
	crdService service.CustomResourceDefinitionService
}

func NewRestCommandApi(crdService service.CustomResourceDefinitionService) CommandApi {

	this := &RestCommandApi{}

	this.crdService = crdService

	return this
}

func (this *RestCommandApi) RunServer() {
	go this.runServer()
}

func displayBindings(bindings []db.DatabaseBinding) []map[string]string {
	lines := make([]map[string]string, 1)

	for _, binding := range bindings {

		text := make(map[string]string)

		text["id"] = binding.Id
		text["namespace"] = binding.Namespace
		text["secret"] = binding.SecretName
		text["databaseId"] = binding.DatabaseInstanceId
		text["action"] = string(binding.Meta.Current.Action)
		text["status"] = string(binding.Meta.Current.State)

		lines = append(lines, text)
	}

	return lines
}

func (this *RestCommandApi) runServer() {
	r := gin.Default()
	r.GET("/list", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"instances": this.crdService.FindAllDatabaseInstances(),
			"bindings":  displayBindings(this.crdService.FindAllDatabaseBindings()),
		})
	})
	r.Run()
}
