package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ojalmeida/GREST/src/config"
	"github.com/ojalmeida/GREST/src/db"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func correctData(behavior db.Behavior, data map[string]string, direction string) (correctedData map[string]string, unknownKeys []string) {

	correctedData = map[string]string{}

	for k, v := range data {

		found := false

		for i := range behavior.KeyMappings {

			bKey := behavior.KeyMappings[i].Key
			bColumn := behavior.KeyMappings[i].Column

			if strings.ToLower(direction) == "inbound" {

				if k == bKey {

					correctedData[bColumn] = v
					found = true
					break

				}
			} else if strings.ToLower(direction) == "outbound" {

				if k == bColumn {

					correctedData[bKey] = v
					found = true
					break

				}

			}

		}

		if !found {

			unknownKeys = append(unknownKeys, k)

		}

	}

	return

}

func toGetPayload(queryString string) GetPayload {

	result := GetPayload{}

	if unescapedQueryString, err := url.QueryUnescape(queryString); err == nil {

		queryString = unescapedQueryString

	}

	params := strings.Split(queryString, "&")

	parametersMap := map[string]string{}

	if queryString == "" {
		return result
	}

	for _, param := range params {

		kv := strings.Split(param, "=")

		parametersMap[kv[0]] = kv[1]

	}

	result.Must = parametersMap

	return result

}

func isValidPayload(method string, payload string) bool {

	switch method {

	case "POST":

		var testPayload PostPayload

		err := json.Unmarshal([]byte(payload), &testPayload)

		return err == nil && testPayload.Set != nil

	case "PUT":

		var testPayload PutPayload

		err := json.Unmarshal([]byte(payload), &testPayload)

		return err == nil && testPayload.Set != nil && testPayload.Must != nil

	case "DELETE":

		var testPayload DeletePayload

		err := json.Unmarshal([]byte(payload), &testPayload)

		return err == nil && testPayload.Must != nil

	}

	return true
}

func isValidQueryString(url string) bool {

	urlFields := strings.Split(url, "?")

	if len(urlFields) == 1 {
		return true
	}

	queryString := urlFields[1]

	params := strings.Split(queryString, "&")

	for i := range params {

		kv := strings.Split(params[i], "=")

		if len(kv) != 2 {

			return false

		}

	}

	return true

}

func getDeclaredBehaviors() ([]db.Behavior, []error) {

	var (
		declaredBehaviors []db.Behavior
		behavior          db.Behavior
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

	readPathMappingDeclaration := func(line string) db.PathMapping {

		args := regexp.MustCompile(`\s+`).Split(line, -1)

		// privacy := args[1]

		table := strings.Trim(args[3], "'")

		path := strings.Trim(args[5], "'")

		return db.PathMapping{Path: path, Table: table}

	}

	readKeyMappingDeclaration := func(line string) db.KeyMapping {

		line = strings.Trim(line, "\t")

		args := regexp.MustCompile(`\s+`).Split(line, -1)

		// privacy := args[1]

		column := strings.Trim(args[3], "'")

		key := strings.Trim(args[5], "'")

		return db.KeyMapping{Key: key, Column: column}

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

		if db.ComparePathMappings(behavior.PathMapping, db.PathMapping{}) && pathMappingRegexp.MatchString(line) {

			if !db.CompareBehaviors(behavior, db.Behavior{}) {

				declaredBehaviors = append(declaredBehaviors, behavior)

				behavior = db.Behavior{}

			}

			behavior.PathMapping = readPathMappingDeclaration(line)

		} else if !db.ComparePathMappings(behavior.PathMapping, db.PathMapping{}) && pathMappingRegexp.MatchString(line) {

			declaredBehaviors = append(declaredBehaviors, behavior)

			behavior = db.Behavior{}

			behavior.PathMapping = readPathMappingDeclaration(line)

		} else if db.ComparePathMappings(behavior.PathMapping, db.PathMapping{}) && keyMappingRegexp.MatchString(line) {

			errs = append(errs, errors.New(fmt.Sprintf("expecting PathMapping declaration in line %d, got \"%s\"", lineNumber, line)))

		} else if !db.ComparePathMappings(behavior.PathMapping, db.PathMapping{}) && keyMappingRegexp.MatchString(line) {

			behavior.KeyMappings = append(behavior.KeyMappings, readKeyMappingDeclaration(line))
		}

	}

	if !db.CompareBehaviors(behavior, db.Behavior{}) {
		declaredBehaviors = append(declaredBehaviors, behavior)
	}

	return declaredBehaviors, errs

}
