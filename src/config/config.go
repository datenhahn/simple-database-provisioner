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

package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type AppConfig struct {
	DbmsServers []DbmsServerConfig `yaml:"dbmsServers"`
}

type DbmsServerConfig struct {
	Name       string       `yaml:"name"`
	Type       string       `yaml:"type"`
	FromSecret SecretConfig `yaml:"fromSecret"`
}

type SecretConfig struct {
	Namespace string `yaml:"namespace"`
	Secret    string `yaml:"secret"`
}

func GetDbmsServer(config AppConfig, dbmsServerName string) (DbmsServerConfig, error) {
	for _, dbmsServer := range config.DbmsServers {
		if dbmsServer.Name == dbmsServerName {
			return dbmsServer, nil
		}
	}
	return DbmsServerConfig{}, fmt.Errorf("Could not find DBMS server with name '%s' in config", dbmsServerName)
}

func ReadConfig(filename string) (AppConfig, error) {

	config := AppConfig{}

	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return config, err
	}

	return config, nil
}
