package db

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

func CreateTables() (err error) {

	transaction, err := LocalConn.Begin()
	if err != nil {
		return
	}

	// SQLite Implementation of Conf Database.
	statement1, _ := transaction.Prepare("CREATE TABLE IF NOT EXISTS path_mapping (path_mapping_id	INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL,`table` TEXT NOT NULL);")
	statement2, _ := transaction.Prepare(`CREATE TABLE IF NOT EXISTS key_mapping (
		key_mapping_id	INTEGER PRIMARY KEY AUTOINCREMENT,
		key			TEXT NOT NULL,
		column		TEXT NOT NULL
	);`)
	statement3, _ := transaction.Prepare(`CREATE TABLE IF NOT EXISTS behavior (
		behavior_id		INTEGER PRIMARY KEY AUTOINCREMENT,
		path_mapping_id	INTEGER NOT NULL,
		key_mapping_id	INTEGER NOT NULL
	);`)
	statement4, _ := transaction.Prepare(`CREATE TABLE IF NOT EXISTS config (
		name	TEXT NOT NULL,
		value	TEXT NOT NULL
	);`)

	_, err = statement1.Exec()
	if err != nil {
		log.Println(err)
		panic(err.Error())
	}

	_, err = statement2.Exec()
	if err != nil {

		err = transaction.Rollback()
		if err != nil {
			return err
		}

	}

	_, err = statement3.Exec()
	if err != nil {

		err = transaction.Rollback()
		if err != nil {
			return err
		}

	}

	_, err = statement4.Exec()
	if err != nil {

		err = transaction.Rollback()
		if err != nil {
			return err
		}

	}

	err = transaction.Commit()
	if err != nil {
		return
	}

	return
}

func FactoryReset() (err error) {

	transaction, err := LocalConn.Begin()

	if err != nil {
		return err
	}

	statement1, err := transaction.Prepare("DROP TABLE behavior;")
	statement2, err := transaction.Prepare("DROP TABLE path_mapping;")
	statement3, err := transaction.Prepare("DROP TABLE key_mapping;")
	statement4, err := transaction.Prepare("DROP TABLE config;")

	_, err = statement1.Exec()
	if err != nil {

		err = transaction.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	_, err = statement2.Exec()
	if err != nil {

		err = transaction.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	_, err = statement3.Exec()
	if err != nil {

		err = transaction.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	_, err = statement4.Exec()
	if err != nil {

		err = transaction.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	err = CreateTables()
	if err != nil {
		return err
	}

	return

}

// getPathMapping returns a PathMapping object present in the database based on given id
func getPathMapping(id int) (PathMapping, error) {

	convertToPathMapping := func(data map[string]string) PathMapping {

		return PathMapping{Path: data["path"], Table: data["table"]}
	}

	var pathMapping PathMapping

	result, err := Read("path_mapping", map[string]string{"path_mapping_id": strconv.Itoa(id)}, "sqlite3-config")

	if err != nil || result == nil {
		return PathMapping{}, err
	}

	pathMapping = convertToPathMapping(result[0])

	return pathMapping, nil

}

// getKeyMapping returns a KeyMapping object present in the database based on given id
func getKeyMapping(id int) (KeyMapping, error) {

	convertToKeyMapping := func(data map[string]string) KeyMapping {

		return KeyMapping{Key: data["key"], Column: data["column"]}
	}

	var keyMapping KeyMapping

	result, err := Read("key_mapping", map[string]string{"key_mapping_id": strconv.Itoa(id)}, "sqlite3-config")

	if err != nil {
		return KeyMapping{}, err
	}

	keyMapping = convertToKeyMapping(result[0])

	return keyMapping, nil

}

// GetBehaviors return an array of behaviors present in the database
func GetBehaviors() ([]Behavior, error) {

	var behaviors []Behavior

	rows, _ := LocalConn.Query("SELECT path_mapping_id FROM behavior GROUP BY path_mapping_id;")

	if rows != nil {

		var pathMappingIds []int

		for rows.Next() {

			var pathMappingId int

			err := rows.Scan(&pathMappingId)

			if err != nil {
				return []Behavior{}, err
			}

			pathMappingIds = append(pathMappingIds, pathMappingId)

		}

		rows.Close()

		for i := range pathMappingIds {

			rows, err := LocalConn.Queryx("SELECT b.key_mapping_id FROM behavior b INNER JOIN key_mapping km ON b.key_mapping_id = km.key_mapping_id WHERE b.path_mapping_id = ?", pathMappingIds[i])

			if err != nil {
				return []Behavior{}, err
			}

			if rows != nil {

				var pathMapping PathMapping
				var keyMappings []KeyMapping
				var keyMappingIds []int
				var keyMappingId int

				for rows.Next() {

					err = rows.Scan(&keyMappingId)

					if err != nil {
						return []Behavior{}, err
					}

					keyMappingIds = append(keyMappingIds, keyMappingId)

				}

				rows.Close()

				for j := range keyMappingIds {

					keyMapping, err := getKeyMapping(keyMappingIds[j])
					if err != nil {
						return []Behavior{}, err
					}

					keyMappings = append(keyMappings, keyMapping)

				}

				pathMapping, err := getPathMapping(pathMappingIds[i])

				if err != nil {
					return []Behavior{}, err
				}

				behavior := Behavior{pathMapping, keyMappings}

				behaviors = append(behaviors, behavior)

			}

		}

	}

	return behaviors, nil

}

// getPathMappings return an array of pathMapping present in the database
func getPathMappings() ([]PathMapping, error) {

	var pathMappings []PathMapping

	var (
		path  string
		table string
	)

	rows, err := LocalConn.Query("SELECT path, `table` FROM path_mapping;")

	if err != nil {
		return []PathMapping{}, err
	}

	if rows != nil {

		defer rows.Close()

		for rows.Next() {

			err = rows.Scan(&path, &table)

			if err != nil {
				return []PathMapping{}, err
			}

			pathMappings = append(pathMappings, PathMapping{path, table})

		}

	}

	return pathMappings, nil

}

// getKeyMappings return an array of keyMapping present in the database
func getKeyMappings() ([]KeyMapping, error) {

	var keyMappings []KeyMapping

	var (
		key    string
		column string
	)

	rows, err := LocalConn.Query("SELECT `key`, `column`  FROM key_mapping;")

	if err != nil {
		return []KeyMapping{}, err
	}

	if rows != nil {

		defer rows.Close()

		for rows.Next() {

			err = rows.Scan(&key, &column)

			if err != nil {
				return []KeyMapping{}, err
			}

			keyMappings = append(keyMappings, KeyMapping{key, column})

		}

	}

	return keyMappings, nil

}
