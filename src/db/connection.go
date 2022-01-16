package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ojalmeida/GREST/src/config"
	"log"
	"os"
)

/*
	RemoteDB is the connection with the user's database (MySQL)
	This func needs e host, port and database to create the connection...
*/
func RemoteDB() *sqlx.DB {
	log.Println("Establishing connection to local database")

	conn, err := sqlx.Open(config.Conf.Database.DBMS, fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		config.Conf.Database.Username,
		config.Conf.Database.Password,
		"tcp",
		config.Conf.Database.Address,
		config.Conf.Database.Port,
		config.Conf.Database.Schema))

	if err != nil {
		log.Println("Fail!")
		panic(err.Error())
	}

	return conn
}

/*
	LocalDB will connect to a local SQLITE database which stores GREST's configuration.
*/
func LocalDB() *sqlx.DB {

	log.Println("Establishing connection to local database")

	dbname := "database.db"
	db, err := os.Open(dbname)

	if err != nil {

		log.Println(dbname, "was not found!")
		log.Println("\t└──Trying to create", dbname)
		_, err := os.Create(dbname)

		if err != nil {
			log.Println("\t\t└──Fail!")
		} else {
			log.Println("\t\t└──Success")
		}

	} else {

		log.Println(dbname, "found!")

	}

	db.Close()

	conn, err := sqlx.Connect("sqlite3", dbname)
	if err != nil {
		panic(err.Error())
	}
	return conn
}
