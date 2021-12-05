package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type databaseAttributes struct {
	username string
	password string
	ip       string
	protocol string
	port     uint16
}

var attributes = databaseAttributes{"root", "root", "127.0.0.1", "tcp", 3306}

var connection *sql.DB

func init() {

	_connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/grest", attributes.username, attributes.password, attributes.protocol, attributes.ip, attributes.port))

	if err != nil {
		panic(err.Error())
	}

	connection = _connection
}

type Behavior struct {
	PathMapping pathMapping
	KeyMapping  keyMapping
}

type keyMapping struct {
	Path  string
	Table string
}

type pathMapping struct {
	Path  string
	Table string
}

func getPathMapping(id int) pathMapping {

	var (
		path  string
		table string
	)

	rows, err := connection.Query(`

		SELECT path_name, table_name FROM path_mappings WHERE path_mapping_id = ?;
`, id)
	defer rows.Close()

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {

		err := rows.Scan(&path, &table)

		if err != nil {
			panic(err.Error())
		}

	}

	return pathMapping{path, table}

}

func getKeyMapping(id int) keyMapping {

	var (
		key    string
		column string
	)

	rows, err := connection.Query(`

		SELECT key_name, column_name FROM key_mappings WHERE key_mapping_id = ?;
`, id)
	defer rows.Close()

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {

		err := rows.Scan(&key, &column)

		if err != nil {
			panic(err.Error())
		}

	}

	return keyMapping{key, column}

}

func GetBehaviors() []Behavior {

	var behaviors []Behavior

	var (
		pathMappingId int
		keyMappingId  int
	)

	rows, err := connection.Query(`

		SELECT path_mapping_id, key_mapping_id FROM behaviors;
`)
	defer rows.Close()

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {

		err := rows.Scan(&pathMappingId, &keyMappingId)

		if err != nil {
			panic(err.Error())
		}

		behaviors = append(behaviors, Behavior{getPathMapping(pathMappingId), getKeyMapping(keyMappingId)})

	}

	return behaviors

}
