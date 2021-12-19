package db

func TableExists(tableName string) bool {

	rows, _ := connection.Query("SELECT TABLE_NAME FROM information_schema.TABLES where TABLE_NAME = ?", tableName)

	if rows != nil {

		return true

	} else {

		return false
	}

}

func ColumnExists(tableName, columnName string) bool {

	rows, _ := connection.Query("SELECT column_name FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND COLUMN_NAME = ?", tableName, columnName)

	if rows != nil {

		return true
	} else {

		return false
	}

}
