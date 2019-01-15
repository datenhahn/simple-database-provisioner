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

package dbms

import (
	"fmt"
	"strconv"
)

// DbmsProvider provides methods for database creation and deletion.
type DbmsProvider interface {
	// CreateDatabaseInstance creates a new database instance on the dbms server with the
	// id dbmsServerId and the provided credentials.
	// The databaseInstanceName is used as name for the new database which is created
	// on the DBMS server.
	// It returns the credentials of the database instance which was created or an
	// error if the creation wasn't successful.
	CreateDatabaseInstance(dbmsServerId string, dbmsServerCredentials DatabaseCredentials, databaseInstanceName string) (DatabaseCredentials, error)

	// ExistsDatabaseInstance checks if a database instance is existing on a DBMS server.
	// Returns true if the database instance exists, false if not.
	// Returns an error if the check fails e.g. due to connection problems.
	ExistsDatabaseInstance(dbmsServerId string, dbmsServerCredentials DatabaseCredentials, databaseInstanceName string) (bool, error)

	// DeleteDatabaseInstance deletes a database instance form a DBMS server. Returns an
	// error if deletion failed, nil if the deletion was successful.
	DeleteDatabaseInstance(dbmsServerId string, dbmsServerCredentials DatabaseCredentials, databaseInstanceName string) error

	// Type returns the type of the DbmsProvider (e.g. postgresql). The type then can be used
	// in the config to specify which DbmsProvider should be used for a certain DBMS server.
	Type() string
}

// DatabaseCredentials contain all information needed to connect to a database.
type DatabaseCredentials struct {
	Host     string
	User     string
	Password string
	Port     int
	Ssl      bool
	Dbname   string
}

// String returns a human readable string representation of the credentials
func (this DatabaseCredentials) String() string {
	return fmt.Sprintf("host=%s, user=%s, password=****, port=%d, ssl=%t, database=%s", this.Host, this.User, this.Port, this.Ssl, this.Dbname)
}

// ToSecretData converts the credentials into a key-value map which can be used
// to store the credentials in a kubernetes secret.
func (this DatabaseCredentials) ToSecretData() (map[string][]byte, error) {

	secretData := make(map[string][]byte)

	secretData["host"] = []byte(this.Host)

	if string(secretData["host"]) == "" {
		return nil, fmt.Errorf("Host may not be empty when creating a secret: %s", this.Host)
	}

	secretData["user"] = []byte(this.User)

	if string(secretData["user"]) == "" {
		return nil, fmt.Errorf("User may not be empty when creating a secret: %s", this.User)
	}

	secretData["password"] = []byte(this.Password)

	if string(secretData["password"]) == "" {
		return nil, fmt.Errorf("Password may not be empty when creating a secret: %s", this.Password)
	}

	secretData["port"] = []byte(strconv.Itoa(this.Port))

	if string(secretData["port"]) == "" {
		return nil, fmt.Errorf("Port may not be empty when creating a secret: %d", this.Port)
	}

	if this.Ssl {

		secretData["ssl"] = []byte("true")
	} else {
		secretData["ssl"] = []byte("false")
	}

	secretData["database"] = []byte(this.Dbname)

	if string(secretData["database"]) == "" {
		return nil, fmt.Errorf("Database may not be empty when creating a secret: %s", this.Dbname)
	}

	return secretData, nil
}

// CreateCredentialsFromSecretData creates database credentials from a key-value map
// the key-value can be extracted from a kubernetes secret data.
func CreateCredentialsFromSecretData(secret map[string][]byte) (DatabaseCredentials, error) {

	credentials := DatabaseCredentials{}

	credentials.Host = string(secret["host"])
	credentials.User = string(secret["user"])
	credentials.Password = string(secret["password"])

	portInt, err := strconv.Atoi(string(secret["port"]))

	if err != nil {
		return DatabaseCredentials{}, err
	}

	credentials.Port = portInt

	sslBool, err := strconv.ParseBool(string(secret["ssl"]))

	if err != nil {
		return DatabaseCredentials{}, err
	}

	credentials.Ssl = sslBool
	credentials.Dbname = string(secret["database"])

	return credentials, nil
}
