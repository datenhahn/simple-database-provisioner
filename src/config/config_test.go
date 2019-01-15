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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {

	expected := AppConfig{
		DbmsServers: []DbmsServerConfig{
			DbmsServerConfig{
				Name: "dbms-dev-postgres",
				FromSecret: SecretConfig{
					Namespace: "default",
					Secret:    "dbms-dev-postgres-secret",
				},
			},
		},
	}

	conf, err := ReadConfig("testdata/config.yaml")

	assert.Nil(t, err)
	assert.Equal(t, expected, conf)

	dbmsConfig, err := GetDbmsServer(conf, "dbms-dev-postgres")
	assert.Nil(t, err)
	assert.Equal(t, "dbms-dev-postgres", dbmsConfig.Name)

	_, err = GetDbmsServer(conf, "dbms-non-postgres")
	assert.Error(t, err)
}

func TestReadConfigFileNotFound(t *testing.T) {

	conf, err := ReadConfig("testdata/doesnotexist.yaml")

	assert.Error(t, err)
	assert.Equal(t, AppConfig{}, conf)
}

func TestReadConfigCannotParse(t *testing.T) {

	conf, err := ReadConfig("testdata/cannotparse.yaml")

	assert.Error(t, err)
	assert.Equal(t, AppConfig{}, conf)
}
