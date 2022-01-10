package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var attributes = databaseAttributes{"root", "root", "127.0.0.1", "tcp", 3306}

type databaseAttributes struct {
	username string
	password string
	ip       string
	protocol string
	port     uint16
}

func GetConnection() *sqlx.DB {

	conn, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/grest", attributes.username, attributes.password, attributes.protocol, attributes.ip, attributes.port))

	if err != nil {
		panic(err.Error())
	}

	return conn

}

// TODO: 1) define config database (sqlite)
// TODO: 2) implement connection with sqlite
// TODO: 3) segregate client db from config db
