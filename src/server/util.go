package server

import (
	"encoding/json"
	"github.com/ojalmeida/GREST/src/db"
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
