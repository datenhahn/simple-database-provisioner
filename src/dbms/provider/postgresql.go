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
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"simple-database-provisioner/src/dbms"
	"simple-database-provisioner/src/util"
	"strings"
)

const DBNAME_REGEX = "[A-Za-z0-9_-]+"

var dbnameRegex = regexp.MustCompile(DBNAME_REGEX)

type PostgresqlDbmsProvider struct {
}

func isValidDatabaseName(dbname string) bool {
	return dbnameRegex.MatchString(dbname)
}

func QuoteValue(param string) string {
	escaped := strings.Replace(param, "'", "''", -1)
	return fmt.Sprintf("'%s'", escaped)
}

func QuoteIdentifier(param string) string {
	escaped := strings.Replace(param, "\"", "\"\"", -1)
	return fmt.Sprintf("\"%s\"", escaped)
}

func connect(credentials dbms.DatabaseCredentials) (*sql.DB, error) {

	mode := "require"
	if !credentials.Ssl {
		mode = "disable"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", credentials.User, credentials.Password, credentials.Host, credentials.Port, credentials.Dbname, mode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {

		return nil, fmt.Errorf("Could not connect to database: %s", credentials.String())
	}

	return db, nil

}

func (this *PostgresqlDbmsProvider) CreateDatabaseInstance(dbmsServerId string, dbmsServerCredentials dbms.DatabaseCredentials, databaseInstanceName string) (dbms.DatabaseCredentials, error) {
	db, err := connect(dbmsServerCredentials)

	if err != nil {
		return dbms.DatabaseCredentials{}, err
	}

	if !isValidDatabaseName(databaseInstanceName) {
		return dbms.DatabaseCredentials{}, fmt.Errorf("Database name '%s' must match regex: %s", databaseInstanceName, DBNAME_REGEX)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", databaseInstanceName))

	if err != nil {
		return dbms.DatabaseCredentials{}, err
	}

	passwd := util.GeneratePassword(20)

	instanceCreds := dbms.DatabaseCredentials{}

	instanceCreds = dbmsServerCredentials
	instanceCreds.Password = passwd
	instanceCreds.User = databaseInstanceName
	instanceCreds.Dbname = databaseInstanceName

	_, err = db.Exec(fmt.Sprintf("CREATE ROLE %s WITH PASSWORD %s LOGIN VALID UNTIL 'infinity';", QuoteIdentifier(databaseInstanceName), QuoteValue(passwd)))

	if err != nil {
		return dbms.DatabaseCredentials{}, err
	}

	_, err = db.Query(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", QuoteIdentifier(databaseInstanceName), QuoteIdentifier(databaseInstanceName)))

	if err != nil {
		return dbms.DatabaseCredentials{}, err
	}

	err = db.Close()

	if err != nil {
		return dbms.DatabaseCredentials{}, err
	}

	return instanceCreds, nil

}

func (this *PostgresqlDbmsProvider) ExistsDatabaseInstance(dbmsServerId string, dbmsServerCredentials dbms.DatabaseCredentials, databaseInstanceName string) (bool, error) {
	db, err := connect(dbmsServerCredentials)

	if err != nil {
		return false, err
	}

	stmt, err := db.Prepare("SELECT 1 FROM pg_database WHERE datname=$1")

	if err != nil {
		return false, err
	}

	row := stmt.QueryRow(databaseInstanceName)

	if err != nil {
		return false, err
	}

	var hasDb int

	err = row.Scan(&hasDb)

	if err != nil {

		// todo JHA: check for error type not string
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		} else {
			return false, err
		}
	}

	if hasDb == 1 {
		return true, nil
	} else {
		return false, nil
	}

}

func (this *PostgresqlDbmsProvider) DeleteDatabaseInstance(dbmsServerId string, dbmsServerCredentials dbms.DatabaseCredentials, databaseInstanceName string) error {

	db, err := connect(dbmsServerCredentials)

	if err != nil {
		return err
	}

	_, err = db.Query(fmt.Sprintf("DROP DATABASE %s;", QuoteIdentifier(databaseInstanceName)))

	if err != nil {
		return err
	}

	_, err = db.Query(fmt.Sprintf("DROP ROLE %s;", QuoteIdentifier(databaseInstanceName)))

	if err != nil {
		return err
	}

	err = db.Close()

	if err != nil {
		return err
	}

	return nil
}

func (this *PostgresqlDbmsProvider) Type() string {
	return "postgresql"
}
