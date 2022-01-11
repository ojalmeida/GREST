package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TableDoesNotExistsError struct {
	err error
}

type ColumnDoesNotExistsError struct {
	err error
}

func (t TableDoesNotExistsError) Error() string {

	return "Table does not exists"
}

func (c ColumnDoesNotExistsError) Error() string {

	return "Column does not exists"
}

// Create inserts data in the database returning an error if it occurs
func Create(tableName string, data map[string]string, driverName string) error {

	query := fmt.Sprintf("INSERT INTO %s ", tableName)
	var columns []string
	var values []string

	for k, v := range data {

		columns = append(columns, fmt.Sprintf("`%s`", k))
		values = append(values, fmt.Sprintf("'%s'", v))

	}

	query += fmt.Sprintf("( %s ) VALUES ( %s )", strings.Join(columns, ", "), strings.Join(values, ", "))

	var err error

	switch driverName {

	case "sqlite3-config":

		_, err = LocalConn.Query(query)

	default:

		_, err = RemoteConn.Query(query)

	}

	return err

}

// Read returns an array of maps containing the results retrieved from database and an error, if it occurs
func Read(tableName string, filters map[string]string, driverName string) (result []map[string]string, err error) {

	if TableExists(tableName) {

		var unparsedResults []map[string]interface{}
		var query string
		var filtersStrings []string

		if len(filters) != 0 {

			query += fmt.Sprintf("SELECT * FROM %s WHERE ", tableName)

			for key, value := range filters {

				if ColumnExists(tableName, key) {
					filtersStrings = append(filtersStrings, fmt.Sprintf("`%s` = '%s'", key, value))
				} else {

					return nil, ColumnDoesNotExistsError{}

				}
			}

			query += strings.Join(filtersStrings, " AND ")
		} else {

			query += fmt.Sprintf("SELECT * FROM %s ", tableName)
		}

		var err error
		var rows *sqlx.Rows

		switch driverName {

		case "sqlite3-config":

			rows, err = LocalConn.Queryx(query)

		default:

			rows, err = RemoteConn.Queryx(query)

		}

		defer rows.Close()

		if err != nil {

			result = []map[string]string{}

			return result, err
		}

		for rows.Next() {

			unparsedResult := make(map[string]interface{})

			err := rows.MapScan(unparsedResult)

			if err != nil {

				result = []map[string]string{}

				return result, err
			}

			unparsedResults = append(unparsedResults, unparsedResult)

		}

		result = ToMapSlice(unparsedResults)

		return result, err

	} else {

		result = nil

		err = TableDoesNotExistsError{}

		return
	}

}

// Update changes data in the database. using the provided filters and data maps, returning an error if it occurs
func Update(tableName string, filters map[string]string, data map[string]string, driverName string) error {

	var filterSlice []string
	var dataSlice []string

	query := fmt.Sprintf("UPDATE %s ", tableName)

	for k, v := range data {

		dataSlice = append(dataSlice, fmt.Sprintf("`%s` = '%s'", k, v))

	}

	query += fmt.Sprintf("SET %s ", strings.Join(dataSlice, ", "))

	for k, v := range filters {

		filterSlice = append(filterSlice, fmt.Sprintf("%s = '%s'", k, v))

	}

	query += fmt.Sprintf("WHERE %s", strings.Join(filterSlice, ", "))

	var err error

	switch driverName {

	case "sqlite3-config":

		_, err = LocalConn.Queryx(query)

	default:

		_, err = RemoteConn.Queryx(query)

	}

	return err
}

// Delete removes a data from database, using the provided filters, returning an error if it occurs
func Delete(tableName string, filters map[string]string, driverName string) error {

	var filterSlice []string

	query := fmt.Sprintf("DELETE FROM %s ", tableName)

	for k, v := range filters {

		filterSlice = append(filterSlice, fmt.Sprintf("`%s` = '%s'", k, v))

	}

	query += fmt.Sprintf("WHERE %s", strings.Join(filterSlice, " AND "))

	var err error

	switch driverName {

	case "sqlite3-config":

		_, err = LocalConn.Queryx(query)

	default:

		_, err = RemoteConn.Queryx(query)

	}

	return err
}
