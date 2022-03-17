package db

import (
	"bufio"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ojalmeida/GREST/src/config"
	log "github.com/ojalmeida/GREST/src/log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
		log.ErrorLogger.Panicln("Fail on trying to create tables :", err.Error())
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

			rows, err := LocalConn.Query("SELECT b.key_mapping_id FROM behavior b INNER JOIN key_mapping km ON b.key_mapping_id = km.key_mapping_id WHERE b.path_mapping_id = ?", pathMappingIds[i])

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

func GetDeclaredBehaviors() ([]Behavior, []error) {

	var (
		declaredBehaviors []Behavior
		behavior          Behavior
		errs              []error
	)

	var pathMappingPattern = `expose\s*?(public|protected)\s*?table\s*?'\w+'\s*?as\s*?'/\w+'`
	var keyMappingPattern = `\t+expose\s*?(public|protected)\s*?column\s*?'\w+'\s*?as\s*?'\w+'`

	var pathMappingRegexp = regexp.MustCompile(pathMappingPattern)
	var keyMappingRegexp = regexp.MustCompile(keyMappingPattern)

	filePath := config.MainFolder + "/mappings.conf"

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {

		errs = append(errs, err)

		return nil, errs
	}

	readPathMappingDeclaration := func(line string) PathMapping {

		args := regexp.MustCompile(`\s+`).Split(line, -1)

		// privacy := args[1]

		table := strings.Trim(args[3], "'")

		path := strings.Trim(args[5], "'")

		return PathMapping{Path: path, Table: table}

	}

	readKeyMappingDeclaration := func(line string) KeyMapping {

		line = strings.Trim(line, "\t")

		args := regexp.MustCompile(`\s+`).Split(line, -1)

		// privacy := args[1]

		column := strings.Trim(args[3], "'")

		key := strings.Trim(args[5], "'")

		return KeyMapping{Key: key, Column: column}

	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {

		line := scanner.Text()
		lineNumber += 1

		// skip empty lines
		if regexp.MustCompile(`^\s*$`).MatchString(line) {

			continue
		}

		if ComparePathMappings(behavior.PathMapping, PathMapping{}) && pathMappingRegexp.MatchString(line) {

			if !CompareBehaviors(behavior, Behavior{}) {

				declaredBehaviors = append(declaredBehaviors, behavior)

				behavior = Behavior{}

			}

			behavior.PathMapping = readPathMappingDeclaration(line)

		} else if !ComparePathMappings(behavior.PathMapping, PathMapping{}) && pathMappingRegexp.MatchString(line) {

			declaredBehaviors = append(declaredBehaviors, behavior)

			behavior = Behavior{}

			behavior.PathMapping = readPathMappingDeclaration(line)

		} else if ComparePathMappings(behavior.PathMapping, PathMapping{}) && keyMappingRegexp.MatchString(line) {

			errs = append(errs, errors.New(fmt.Sprintf("expecting PathMapping declaration in line %d, got \"%s\"", lineNumber, line)))

		} else if !ComparePathMappings(behavior.PathMapping, PathMapping{}) && keyMappingRegexp.MatchString(line) {

			behavior.KeyMappings = append(behavior.KeyMappings, readKeyMappingDeclaration(line))
		}

	}

	if !CompareBehaviors(behavior, Behavior{}) {
		declaredBehaviors = append(declaredBehaviors, behavior)
	}

	return declaredBehaviors, errs

}
