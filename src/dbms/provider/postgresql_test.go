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

package provider

import (
	"github.com/stretchr/testify/assert"
	"simple-database-provisioner/src/dbms"
	"testing"
)

/**
 * This test requires a postgres database running with the following connection data.
 * e.g. Spawn with:
 *
 *    docker run --rm --name sdp-postgres-testdb -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres postgres
 */

var TEST_CREDENTIALS = dbms.DatabaseCredentials{
	Dbname:   "postgres",
	User:     "postgres",
	Password: "postgres",
	Port:     5432,
	Host:     "localhost",
	Ssl:      false,
}

var UNREACHABLE_HOST_CREDENTIALS = dbms.DatabaseCredentials{
	Dbname:   "postgres",
	User:     "postgres",
	Password: "postgres",
	Port:     5432,
	Host:     "hostdoesnotexist",
	Ssl:      false,
}

var WRONG_CREDENTIALS = dbms.DatabaseCredentials{
	Dbname:   "postgres",
	User:     "postgres",
	Password: "WRONGPASSWORD",
	Port:     5432,
	Host:     "localhost",
	Ssl:      false,
}

func TestPostgresqlDbmsProvider_WrongPassword(t *testing.T) {

	postgresqlProvider := &PostgresqlDbmsProvider{}

	_, err := postgresqlProvider.CreateDatabaseInstance("golang-test-persistence", WRONG_CREDENTIALS, "foo")

	assert.True(t, err != nil)
	assert.Regexp(t, "password authentication failed", err.Error())
}

func TestPostgresqlDbmsProvider_CouldNotConnect(t *testing.T) {

	postgresqlProvider := &PostgresqlDbmsProvider{}

	_, err := postgresqlProvider.CreateDatabaseInstance("golang-test-persistence", UNREACHABLE_HOST_CREDENTIALS, "foo")

	assert.NotNil(t, err)
	assert.Regexp(t, "dial tcp: lookup", err.Error())
}

// creates a test database and ensures it exists
func createDatabase(t *testing.T, postgresqlProvider dbms.DbmsProvider, dbname string) {

	cred, err := postgresqlProvider.CreateDatabaseInstance("golang-test-persistence", TEST_CREDENTIALS, dbname)

	assert.Nil(t, err)
	assert.Equal(t, "golang-test-persistence", cred.Dbname)
	assert.Equal(t, "golang-test-persistence@localhost", cred.User)
	assert.Equal(t, 5432, cred.Port)
	assert.Equal(t, "localhost", cred.Host)
	assert.Equal(t, false, cred.Ssl)
	assert.NotEqual(t, TEST_CREDENTIALS.Password, cred.Password)

	doesExist, err := postgresqlProvider.ExistsDatabaseInstance("golang-test-persistence", TEST_CREDENTIALS, "golang-test-persistence")

	assert.Nil(t, err)
	assert.True(t, doesExist)
}

// Deletes the test database and ensures it is deleted
func deleteDatabase(t *testing.T, postgresqlProvider dbms.DbmsProvider, dbname string) {

	err := postgresqlProvider.DeleteDatabaseInstance("golang-test-persistence", TEST_CREDENTIALS, "golang-test-persistence")

	assert.Nil(t, err)

	doesExist, err := postgresqlProvider.ExistsDatabaseInstance("golang-test-persistence", TEST_CREDENTIALS, dbname)

	assert.Nil(t, err)
	assert.False(t, doesExist)
}

func TestPostgresqlDbmsProvider_CreateDatabaseInstance(t *testing.T) {

	postgresqlProvider := &PostgresqlDbmsProvider{}

	createDatabase(t, postgresqlProvider, "golang-test-persistence")
	deleteDatabase(t, postgresqlProvider, "golang-test-persistence")

}

func TestPostgresqlDbmsProvider_CreateDatabaseTwoTimes(t *testing.T) {

	postgresqlProvider := &PostgresqlDbmsProvider{}
	dbName := "golang-test-persistence"

	createDatabase(t, postgresqlProvider, dbName)
	_, err := postgresqlProvider.CreateDatabaseInstance("golang-test-persistence", TEST_CREDENTIALS, dbName)

	assert.NotNil(t, err)
	assert.Regexp(t, "database \"golang-test-persistence\" already exists", err.Error())
	deleteDatabase(t, postgresqlProvider, "golang-test-persistence")

}

func TestPostgresqlDbmsProvider_CreateDatabaseInstanceWithSpecialChars(t *testing.T) {

	postgresqlProvider := &PostgresqlDbmsProvider{}
	dbName := "123$%&"
	_, err := postgresqlProvider.CreateDatabaseInstance("golang-test-persistence", TEST_CREDENTIALS, dbName)

	assert.NotNil(t, err)
	assert.Regexp(t, "Database name '123\\$%&' must match regex", err.Error())
}
