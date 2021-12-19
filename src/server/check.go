package server

import "encoding/json"

// isValidPayload checks if a request payload matches with the defined models of payload
func isValidPayload(method string, payload string) bool {

	switch method {

	case "GET":

		var testPayload GetPayload

		err := json.Unmarshal([]byte(payload), &testPayload)

		return (err == nil && testPayload.Match != nil) || payload == ""

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
