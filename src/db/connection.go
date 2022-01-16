package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)



//type databaseAttributes struct {
//	dbms     string
//	username string
//	password string
//	ip       string
//	protocol string
//	port     string
//	db       string
//}

/*
	RemoteDB is the connection with the user's database (MySQL)
	This func needs e host, port and database to create the connection...
*/
func RemoteDB() *sqlx.DB {
	log.Println("Establishing connection to local database")
	//attributes := databaseAttributes{
	//	dbms:     Conf.Database.DBMS,
	//	username: Conf.Database.Username,
	//	password: Conf.Database.Password,
	//	protocol: "tcp",
	//	ip:       Conf.Database.Address,
	//	port:     Conf.Database.Port,
	//	db:       Conf.Database.Database,
	//}
	conn, err := sqlx.Open(Conf.Database.DBMS, fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		Conf.Database.Username,
		Conf.Database.Password,
		"tcp",
		Conf.Database.Address,
		Conf.Database.Port,
		Conf.Database.Database))

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
