package server

import (
	"encoding/json"
	"strings"
)

// isValidPayload checks if a request payload matches with the defined models of payload
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
