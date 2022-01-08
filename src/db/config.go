package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type databaseAttributes struct {
	username string
	password string
	ip       string
	protocol string
	port     uint16
}

var attributes = databaseAttributes{"root", "root", "127.0.0.1", "tcp", 3306}

var connection *sqlx.DB

var requiredPathMappings = []PathMapping{
	{Path: "/config/path-mappings", Table: "path_mappings"},
	{Path: "/config/key-mappings", Table: "key_mappings"},
	{Path: "/config/behaviors", Table: "behaviors"},
}
var requiredKeyMappings = []KeyMapping{
	{Key: "id", Column: "path_mapping_id"},
	{Key: "id", Column: "key_mapping_id"},
	{Key: "id", Column: "behavior_id"},
	{Key: "key", Column: "key"},
	{Key: "column", Column: "column"},
	{Key: "path", Column: "path"},
	{Key: "table", Column: "table"},
	{Key: "path_mapping_id", Column: "path_mapping_id"},
	{Key: "key_mapping_id", Column: "key_mapping_id"},
}
var requiredBehaviors = []Behavior{
	{PathMapping: PathMapping{Path: "/config/path-mappings", Table: "path_mappings"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "path_mapping_id"},
		{Key: "path", Column: "path"},
		{Key: "table", Column: "table"},
	},
	},
	{PathMapping: PathMapping{Path: "/config/key-mappings", Table: "key_mappings"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "key_mapping_id"},
		{Key: "key", Column: "key"},
		{Key: "column", Column: "column"},
	},
	},
	{PathMapping: PathMapping{Path: "/config/behaviors", Table: "behaviors"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "behavior_id"},
		{Key: "path_mapping_id", Column: "path_mapping_id"},
		{Key: "key_mapping_id", Column: "key_mapping_id"},
	},
	},
}

func init() {

	openConnection()

}

type Behavior struct {
	PathMapping PathMapping
	KeyMappings []KeyMapping
}

type KeyMapping struct {
	Key    string
	Column string
}

type PathMapping struct {
	Path  string
	Table string
}

// getPathMapping returns a PathMapping object present in the database based on given id
func getPathMapping(id int) (PathMapping, error) {

	convertToPathMapping := func(data map[string]string) PathMapping {

		return PathMapping{Path: data["path"], Table: data["table"]}
	}

	var pathMapping PathMapping

	result, err := Read("path_mappings", map[string]string{"path_mapping_id": strconv.Itoa(id)})

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

	result, err := Read("key_mappings", map[string]string{"key_mapping_id": strconv.Itoa(id)})

	if err != nil {
		return KeyMapping{}, err
	}

	keyMapping = convertToKeyMapping(result[0])

	return keyMapping, nil

}

// GetBehaviors return an array of behaviors present in the database
func GetBehaviors() ([]Behavior, error) {

	var behaviors []Behavior

	rows, _ := connection.Query("SELECT path_mapping_id FROM behaviors group by path_mapping_id;")

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

			rows, err := connection.Queryx("SELECT b.key_mapping_id FROM behaviors b INNER JOIN key_mappings km ON b.key_mapping_id = km.key_mapping_id WHERE b.path_mapping_id = ?", pathMappingIds[i])

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

	rows, err := connection.Query("SELECT path, `table` FROM path_mappings;")

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

	rows, err := connection.Query("SELECT `key`, `column`  FROM key_mappings;")

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
		pathMappingsQuery := "INSERT INTO path_mappings VALUES "
		ID := 20000
		var pathMappingsThatWillBeInserted []string
		for i := range requiredPathMappings {

			path := requiredPathMappings[i].Path
			table := requiredPathMappings[i].Table

			pathMappingsThatWillBeInserted = append(pathMappingsThatWillBeInserted, fmt.Sprintf("(%d, '%s', '%s')", ID, path, table))

			ID += 1
		}
		pathMappingsQuery += strings.Join(pathMappingsThatWillBeInserted, ",")
		_, err = connection.Query(pathMappingsQuery)
		if err != nil {
			return err
		}
	}

	// KeyMappings
	{
		keyMappingsQuery := "INSERT INTO key_mappings VALUES "
		ID := 30000
		var keyMappingsThatWillBeInserted []string
		for i := range requiredKeyMappings {

			key := requiredKeyMappings[i].Key
			column := requiredKeyMappings[i].Column

			keyMappingsThatWillBeInserted = append(keyMappingsThatWillBeInserted, fmt.Sprintf("(%d, '%s', '%s')", ID, key, column))
			ID += 1
		}
		keyMappingsQuery += strings.Join(keyMappingsThatWillBeInserted, ", ")
		_, err = connection.Query(keyMappingsQuery)
		if err != nil {
			return err
		}
	}

	// Behaviors
	{
		behaviorsQuery := "INSERT INTO behaviors VALUES "
		ID := 10000
		for i := range requiredBehaviors {

			pathMapping := requiredBehaviors[i].PathMapping
			keyMappings := requiredBehaviors[i].KeyMappings
			var behaviorsThatWillBeInserted []string

			var pathMappingsID int
			var keyMappingIDs []int

			rows, err := connection.Query("SELECT path_mapping_id FROM path_mappings WHERE path = ? AND `table` = ?", pathMapping.Path, pathMapping.Table)

			if err != nil {

				panic(err.Error())

			}

			defer rows.Close()

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

				rows, err := connection.Query("SELECT key_mapping_id FROM key_mappings WHERE `key` = ? AND `column` = ?", key, column)

				if err != nil {

					panic(err.Error())

				}

				defer rows.Close()

				for rows.Next() {

					err := rows.Scan(&keyMappingID)

					if err != nil {

						panic(err.Error())

					}

				}

				keyMappingIDs = append(keyMappingIDs, keyMappingID)

			}

			for j := range keyMappingIDs {

				behaviorsThatWillBeInserted = append(behaviorsThatWillBeInserted, fmt.Sprintf("(%d, %d, %d)", ID, pathMappingsID, keyMappingIDs[j]))

				ID += 1

			}

			behaviorsQuery += strings.Join(behaviorsThatWillBeInserted, ", ")

		}

		_, err = connection.Query(behaviorsQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

func openConnection() {

	_connection, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/grest", attributes.username, attributes.password, attributes.protocol, attributes.ip, attributes.port))

	if err != nil {
		panic(err.Error())
	}

	connection = _connection

}
