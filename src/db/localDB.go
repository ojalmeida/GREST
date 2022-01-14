package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var LocalConn *sqlx.DB

func init() {
	// Remote connection.
	LocalConn = LocalDB()
}

func CheckLocalDB() bool {

	transaction, err := LocalConn.Begin()

	if err != nil {

		return false

	}

	statement, err := transaction.Prepare(`

		SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'behavior'
															OR name = 'path_mapping'
															OR name = 'key_mapping'
															OR name = 'config';
	`)
	rows, err := statement.Query()

	if err != nil {

		return false
	}

	numberOfTables := 0

	for rows.Next() {

		numberOfTables += 1

	}

	rows.Close()

	_ = transaction.Commit()

	return numberOfTables == 4

}
