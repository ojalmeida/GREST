package server

import (
	"grest/db"
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

func convertQueryStringToGETPayload(queryString string) GetPayload {

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
