package db

import (
	"database/sql"
	"fmt"
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

			parsedDatum[k] = fmt.Sprintf("%v", v)
		}

		parsedData = append(parsedData, parsedDatum)

	}
	return

}

func TableExists(tableName string, driverName string) bool {

	var rows *sql.Rows

	switch driverName {

	case "sqlite3-config":

		rows, _ = LocalConn.Query("SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?", tableName)
		defer rows.Close()

		if rows != nil {

			for rows.Next() {

				return true
			}

		}

	case "mysql":

		rows, _ = RemoteConn.Query("SELECT TABLE_NAME FROM information_schema.TABLES where TABLE_NAME = ?", tableName)
		defer rows.Close()

		if rows != nil {

			for rows.Next() {

				return true
			}

		}
	}
	return false

}

func ColumnExists(tableName, columnName, driverName string) bool {

	switch driverName {

	case "sqlite3":

		rows, err := RemoteConn.Query("SELECT ? FROM ?", columnName, tableName)

		if err != nil {

			return true

		} else {

			return false
		}

		rows.Close()

	case "mysql":

		rows, err := RemoteConn.Query("SELECT column_name FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND COLUMN_NAME = ?", tableName, columnName)

		if err != nil {
			return false
		}

		defer rows.Close()

		if rows != nil {

			return true

		} else {

			return false

		}

	}

	return false
}
