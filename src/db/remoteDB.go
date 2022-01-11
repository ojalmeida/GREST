package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var RemoteConn *sqlx.DB

func init() {

	// Local Connection for Configuration.
	// LConn := LocalDB()
	RemoteConn = RemoteDB()

}
