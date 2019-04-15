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
	"github.com/asdine/storm"
	"github.com/stretchr/testify/assert"
	"os"
	"simple-database-provisioner/src/util"
	"testing"
	"time"
)

func TestStormPersistenceBackend_AddDatabaseBinding(t *testing.T) {

	dbFile := util.CreateTempFile()
	defer os.Remove(dbFile)

	db, err := storm.Open(dbFile)
	assert.Nil(t, err)

	backend := NewStormPersistenceBackend(db)

	updated, err := time.Parse("Mon Jan 2 15:04:05", "Mon Jan 2 15:04:05")
	assert.Nil(t, err)

	binding := DatabaseBinding{
		NamespaceUniqueId:  "foo-bar",
		Namespace:          "foo",
		SecretName:         "mytest",
		K8sName:            "foo-bar-binding",
		DatabaseInstanceId: "myId",
		Meta: Meta{
			Current: State{
				State:      PENDING,
				Action:     CREATE,
				LastUpdate: updated,
				Message:    "create db",
			},
		},
	}

	err = backend.AddDatabaseBinding(binding)
	assert.Nil(t, err)

	bindings := backend.FindAllDatabaseBindings()

	assert.Equal(t, []DatabaseBinding{binding}, bindings)

	defer db.Close()
}
