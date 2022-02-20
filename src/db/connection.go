package db

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	var conn *sqlx.DB
	var err error

	switch config.Conf.Database.DBMS {

	case "sqlite3":

		conn, err = sqlx.Open(config.Conf.Database.DBMS, config.Conf.Database.Address)

	case "mysql", "mariadb":

		conn, err = sqlx.Open(config.Conf.Database.DBMS, fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
			config.Conf.Database.Username,
			config.Conf.Database.Password,
			"tcp",
			config.Conf.Database.Address,
			config.Conf.Database.Port,
			config.Conf.Database.Schema))

	case "sqlserver":

		conn, err = sqlx.Open(config.Conf.Database.DBMS, fmt.Sprintf("%s://%s:%s@%s:%s?database=%s",
			config.Conf.Database.DBMS,
			config.Conf.Database.Username,
			config.Conf.Database.Password,
			config.Conf.Database.Address,
			config.Conf.Database.Port,
			config.Conf.Database.Schema))

	case "postgres":

		conn, err = sqlx.Open(config.Conf.Database.DBMS, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Conf.Database.Address,
			config.Conf.Database.Port,
			config.Conf.Database.Username,
			config.Conf.Database.Password,
			config.Conf.Database.Schema))

	}

	if err != nil || conn == nil {
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

	var home string
	var err error

	log.Println("Opening database file")

	// Get user home
	if home, err = os.UserHomeDir(); err != nil {

		log.Fatal("Impossible to get user home directory")

	}

	var mainFolder = home + "/.grest"

	if _, err = os.Stat(mainFolder); os.IsNotExist(err) {

		log.Println(mainFolder + "does not exists, trying to create")

		if err = os.Mkdir(mainFolder, 0660); err != nil {

			log.Fatal(err.Error())

		}

		log.Println("\t└──Success")

	}

	dbname := mainFolder + "/database.db"
	db, err := os.Open(dbname)

	if err != nil {

		log.Println(dbname, "was not found!")
		log.Println("\t└──Trying to create", dbname)
		_, err := os.Create(dbname)

		if err != nil {
			log.Println("\t\t└──Fail!")
			panic(err.Error())
		} else {
			log.Println("\t\t└──Success")
		}

	} else {

		log.Println(dbname, "found!")

	}

	db.Close()

	conn, err := sqlx.Open("sqlite3", dbname)
	if err != nil {
		panic(err.Error())
	}
	return conn
}
