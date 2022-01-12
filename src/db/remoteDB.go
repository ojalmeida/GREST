package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var RemoteConn *sqlx.DB

func init() {

	// Remote Connection for Configuration.
	RemoteConn = RemoteDB()

}
