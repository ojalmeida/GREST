package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"

	//"log"
	"strconv"
	"strings"
)

func createPrerequisites() (err error) {

	transaction, err := LocalConn.Begin()
	if err != nil {
		return
	}

	// MySQL Implementation of Conf Database.
	// statement1, _ := transaction.Prepare("CREATE TABLE IF NOT EXISTS path_mappings (path_mapping_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY , path VARCHAR(255) NOT NULL, `table` VARCHAR(255) NOT NULL) AUTO_INCREMENT=20000;")
	// statement2, _ := transaction.Prepare("CREATE TABLE IF NOT EXISTS key_mappings (key_mapping_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY , `key` VARCHAR(255) NOT NULL, `column` VARCHAR(255) NOT NULL) AUTO_INCREMENT=30000;")
	// statement3, _ := transaction.Prepare("CREATE TABLE IF NOT EXISTS behaviors (behavior_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY , path_mapping_id INT NOT NULL, key_mapping_id INT NOT NULL) AUTO_INCREMENT=10000;")

	// SQLite Implementation of Conf Database.
	statement1, _ := transaction.Prepare( "CREATE TABLE IF NOT EXISTS path_mappings (path_mapping_id	INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL,`table` TEXT NOT NULL);")
	statement2, _ := transaction.Prepare(`CREATE TABLE IF NOT EXISTS key_mappings (
		key_mapping_id	INTEGER PRIMARY KEY AUTOINCREMENT,
		key			TEXT NOT NULL,
		column		TEXT NOT NULL
	);`)
	statement3, _ := transaction.Prepare(`CREATE TABLE IF NOT EXISTS behaviors (
		behavior_id		INTEGER PRIMARY KEY AUTOINCREMENT,
		path_mapping_id	INTEGER NOT NULL,
		key_mapping_id	INTEGER NOT NULL
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

	statement1, err := transaction.Prepare("DROP TABLE `behaviors`;")
	statement2, err := transaction.Prepare("DROP TABLE `path_mappings`;")
	statement3, err := transaction.Prepare("DROP TABLE `key_mappings`;")

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

	err = transaction.Commit()
	if err != nil {
		return err
	}

	err = createPrerequisites()
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

	result, err := Read("path_mappings", map[string]string{"path_mapping_id": strconv.Itoa(id)}, "sqlite3-config")

	if err != nil {
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

	result, err := Read("key_mappings", map[string]string{"key_mapping_id": strconv.Itoa(id)}, "sqlite3-config")

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

		defer rows.Close()

		var pathMappingIds []int

		for rows.Next() {

			var pathMappingId int

			err := rows.Scan(&pathMappingId)

			if err != nil {
				return []Behavior{}, err
			}

			pathMappingIds = append(pathMappingIds, pathMappingId)

		}

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

	rows, err := LocalConn.Query("SELECT path, _table FROM path_mappings;")

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

	rows, err := LocalConn.Query("SELECT _key, _column  FROM key_mappings;")

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

// CheckConfigs verifies if the tables for GREST functionality exists, returning True if it exists or False if not and also slices of missing elements
func CheckConfigs() (ok bool, missingPathMappings []PathMapping, missingKeyMappings []KeyMapping, missingBehaviors []Behavior) {

	currentPathMappings, err := getPathMappings()

	if err != nil {
		return false, nil, nil, nil
	}
	currentKeyMappings, err := getKeyMappings()

	if err != nil {
		return false, nil, nil, nil
	}

	currentBehaviors, err := GetBehaviors()

	if err != nil {
		return false, nil, nil, nil
	}

	missingPathMappings = []PathMapping{}
	missingKeyMappings = []KeyMapping{}
	missingBehaviors = []Behavior{}

	// populate missingPathMappings
	for i := range requiredPathMappings {

		contains := false
		pathMapping := requiredPathMappings[i]

		for j := range currentPathMappings {

			if ComparePathMappings(currentPathMappings[j], pathMapping) {

				contains = true
				break

			}

		}

		if !contains {

			missingPathMappings = append(missingPathMappings, pathMapping)

		}

	}

	// populate missingKeyMappings
	for i := range requiredKeyMappings {

		contains := false
		keyMapping := requiredKeyMappings[i]

		for j := range currentKeyMappings {

			if CompareKeyMappings(currentKeyMappings[j], keyMapping) {

				contains = true
				break

			}

		}

		if !contains {

			missingKeyMappings = append(missingKeyMappings, keyMapping)

		}

	}

	// populate missingBehaviors
	for i := range requiredBehaviors {

		contains := false
		behavior := requiredBehaviors[i]

		for j := range currentBehaviors {

			if CompareBehaviors(currentBehaviors[j], behavior) {

				contains = true
				break

			}

		}

		if !contains {

			missingBehaviors = append(missingBehaviors, behavior)

		}

	}

	ok = len(missingPathMappings) == 0 && len(missingKeyMappings) == 0 && len(missingBehaviors) == 0

	return
}

func PopulateConfigs() error {

	err := createPrerequisites()
	if err != nil {
		return err
	}

	// PathMappings
	{
		pathMappingsQuery := "INSERT OR IGNORE INTO path_mapping VALUES "

		ID := 20000

		var pathMappingsThatWillBeInserted []string

		for i := range requiredPathMappings {

			path := requiredPathMappings[i].Path
			table := requiredPathMappings[i].Table

			pathMappingsThatWillBeInserted = append(pathMappingsThatWillBeInserted, fmt.Sprintf("(%d, '%s', '%s')", ID, path, table))

			ID += 1
		}

		pathMappingsQuery += strings.Join(pathMappingsThatWillBeInserted, ",")

		_, err = LocalConn.Query(pathMappingsQuery)
		if err != nil {
			return err
		}
	}

	// KeyMappings
	{
		keyMappingsQuery := "INSERT OR IGNORE INTO key_mapping VALUES "

		ID := 30000

		var keyMappingsThatWillBeInserted []string

		for i := range requiredKeyMappings {

			key := requiredKeyMappings[i].Key
			column := requiredKeyMappings[i].Column

			keyMappingsThatWillBeInserted = append(keyMappingsThatWillBeInserted, fmt.Sprintf("(%d, '%s', '%s')", ID, key, column))
			ID += 1
		}

		keyMappingsQuery += strings.Join(keyMappingsThatWillBeInserted, ", ")

		_, err = LocalConn.Query(keyMappingsQuery)
		if err != nil {
			return err
		}
	}

	// Behaviors
	{
		behaviorsQuery := "INSERT INTO behavior VALUES "

		ID := 10000

		var behaviorsThatWillBeInserted []string

		for i := range requiredBehaviors {

			pathMapping := requiredBehaviors[i].PathMapping
			keyMappings := requiredBehaviors[i].KeyMappings

			var pathMappingsID int
			var keyMappingIDs []int

			rows, err := LocalConn.Query("SELECT path_mapping_id FROM path_mapping WHERE path = ? AND _table = ?", pathMapping.Path, pathMapping.Table)

			if err != nil {

				panic(err.Error())

			}

			for rows.Next() {

				err := rows.Scan(&pathMappingsID)

				if err != nil {

					panic(err.Error())

				}

			}

			for j := range keyMappings {

				key := keyMappings[j].Key
				column := keyMappings[j].Column

				var keyMappingID int

				rows, err := LocalConn.Query("SELECT key_mapping_id FROM key_mapping WHERE _key = ? AND _column = ?", key, column)

				if err != nil {

					panic(err.Error())

				}

				for rows.Next() {

					err := rows.Scan(&keyMappingID)

					if err != nil {

						panic(err.Error())

					}

				}

				keyMappingIDs = append(keyMappingIDs, keyMappingID)

				rows.Close()

			}

			for j := range keyMappingIDs {

				behaviorsThatWillBeInserted = append(behaviorsThatWillBeInserted, fmt.Sprintf("(%d, %d, %d)", ID, pathMappingsID, keyMappingIDs[j]))

				ID += 1

			}

			rows.Close()

		}

		behaviorsQuery += strings.Join(behaviorsThatWillBeInserted, ", ")

		_, err = LocalConn.Query(behaviorsQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
### MySQL
CREATE TABLE IF NOT EXISTS path_mappings (
	path_mapping_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	path VARCHAR(255) NOT NULL,
	`table` VARCHAR(255) NOT NULL
) AUTO_INCREMENT=20000;

CREATE TABLE IF NOT EXISTS key_mappings (
	key_mapping_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`key` VARCHAR(255) NOT NULL,
	`column` VARCHAR(255) NOT NULL
) AUTO_INCREMENT=30000;

CREATE TABLE IF NOT EXISTS behaviors (
	behavior_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	path_mapping_id INT NOT NULL,
	key_mapping_id INT NOT NULL
) AUTO_INCREMENT=10000;



### SQLite
CREATE TABLE IF NOT EXISTS path_mappings (
	path_mapping_id	INTEGER PRIMARY KEY AUTOINCREMENT,
	path			TEXT NOT NULL,
	`table`			TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS key_mappings (
	key_mapping_id	INTEGER PRIMARY KEY AUTOINCREMENT,
	`key`			TEXT NOT NULL,
	`column`		TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS behaviors (
	behavior_id		INTEGER PRIMARY KEY AUTOINCREMENT,
	path_mapping_id	INTEGER NOT NULL,
	key_mapping_id	INTEGER NOT NULL
);
*/
