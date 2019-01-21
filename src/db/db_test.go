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

package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabaseInstance_PrefixedDatabaseName(t *testing.T) {

	longInstance := DatabaseInstance{
		Namespace:    "mysuperlongnamespacename",
		DatabaseName: "and-my-even-much-much-much-longer-even-much-much-longer-and-lon",
	}

	assert.Equal(t, "mysuperlongnamespacename-and-my-even-much-much-much-lo-8157f301", longInstance.PrefixedDatabaseName())
	t.Log(longInstance.PrefixedDatabaseName())

	shortInstance := DatabaseInstance{
		Namespace:    "shortns",
		DatabaseName: "short-db",
	}

	assert.Equal(t, "shortns-short-db-34ebb06c", shortInstance.PrefixedDatabaseName())
}
