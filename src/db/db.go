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

type PathMapping struct {
	Path  string
	Table string
}

func GetPathMappings() []PathMapping {

	var pathMappings []PathMapping

	var (
		id    int
		path  string
		table string
	)

	rows, err := connection.Query(`

		SELECT * FROM path_mappings;
`)
	defer rows.Close()

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {

		err := rows.Scan(&id, &path, &table)

		if err != nil {
			panic(err.Error())
		}

		pathMappings = append(pathMappings, PathMapping{path, table})

	}

	return pathMappings

}
