package db

import (
	"database/sql"
	"fmt"
	"github.com/ojalmeida/GREST/src/log"
)

func CompareBehaviors(behavior1 Behavior, behavior2 Behavior) bool {

	if !ComparePathMappings(behavior1.PathMapping, behavior2.PathMapping) {
		return false
	}

	if !(len(behavior1.KeyMappings) == len(behavior2.KeyMappings)) {
		return false
	}

	for i := range behavior1.KeyMappings {

		if !CompareKeyMappings(behavior1.KeyMappings[i], behavior2.KeyMappings[i]) {
			return false
		}

	}

	return true
}

func CompareKeyMappings(keyMapping1 KeyMapping, keyMapping2 KeyMapping) bool {

	if keyMapping1.Key == keyMapping2.Key && keyMapping1.Column == keyMapping2.Column {
		return true
	} else {
		return false
	}

}

func ComparePathMappings(pathMapping1 PathMapping, pathMapping2 PathMapping) bool {

	if pathMapping1.Path == pathMapping2.Path && pathMapping1.Table == pathMapping2.Table {
		return true
	} else {
		return false
	}

}

func ToMapSlice(unparsedData []map[string]interface{}) (parsedData []map[string]string) {

	for index := range unparsedData {

		var parsedDatum = map[string]string{}

		for k, v := range unparsedData[index] {

			switch v.(type) {

			case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:

				parsedDatum[k] = fmt.Sprintf("%d", v)

			default:

				parsedDatum[k] = fmt.Sprintf("%s", v)

			}

		}

		parsedData = append(parsedData, parsedDatum)

	}
	return

}

func TableExists(tableName string, driverName string) bool {

	var rows *sql.Rows
	var err error

	switch driverName {

	case "sqlite3-config":

		rows, _ = LocalConn.Query("SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?", tableName)
		defer rows.Close()

		if rows != nil {

			for rows.Next() {

				return true
			}

		}

	case "sqlite3":

		rows, _ = RemoteConn.Query("SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?", tableName)
		defer rows.Close()

		if rows != nil {

			for rows.Next() {

				return true
			}

		}

	case "mysql", "mariadb":

		rows, err = RemoteConn.Query("SELECT table_name FROM information_schema.TABLES WHERE table_name = ?", tableName)
		defer rows.Close()

		if err != nil {

			log.ErrorLogger.Println(err.Error())

			return false

		}

		if rows != nil {

			for rows.Next() {

				return true
			}

		}

	case "sqlserver":

		rows, err = RemoteConn.Query("SELECT table_name FROM information_schema.TABLES WHERE table_name = @p1", tableName)
		defer rows.Close()

		if err != nil {

			log.ErrorLogger.Println(err.Error())
			return false

		}

		if rows != nil {

			for rows.Next() {

				return true
			}

		}

	case "postgres":

		rows, err = RemoteConn.Query("SELECT table_name FROM information_schema.TABLES WHERE table_name = $1", tableName)
		defer rows.Close()

		if err != nil {

			log.ErrorLogger.Println(err.Error())
			return false

		}

		if rows != nil {

			for rows.Next() {

				return true
			}

		}

	}

	return false

}

func ColumnExists(tableName, columnName, driverName string) bool {

	var rows *sql.Rows
	var err error

	switch driverName {

	case "sqlite3-config":

		rows, err = LocalConn.Query("SELECT ? FROM "+tableName+" LIMIT 1", columnName)

		defer rows.Close()

		if err != nil {

			return false
		}

		return true

	case "sqlite3":

		rows, err = RemoteConn.Query("SELECT ? FROM "+tableName+" LIMIT 1", columnName)

		defer rows.Close()

		if err != nil {

			return false
		}

		return true

	case "mysql", "mariadb":

		rows, err = RemoteConn.Query("SELECT ? FROM "+tableName+" LIMIT 1", columnName)

		defer rows.Close()

		if err != nil {

			return false
		}

		return true

	case "sqlserver":

		// table name can not be a query parameter in sql server
		rows, err = RemoteConn.Query("SELECT @p1 FROM "+tableName, columnName)

		defer rows.Close()

		if err != nil {

			return false
		}

		return true

	case "postgres":

		// table name can not be a query parameter in postgres
		rows, err = RemoteConn.Query("SELECT $1 FROM "+tableName+" LIMIT 1", columnName)

		defer rows.Close()

		if err != nil {

			return false
		}

		return true

	}

	return false
}
