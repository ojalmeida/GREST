package handle

import (
	"encoding/json"
	"fmt"
	"github.com/ojalmeida/GREST/src/db"
	"github.com/ojalmeida/GREST/src/db/operation"
	"github.com/ojalmeida/GREST/src/server"
	"github.com/ojalmeida/GREST/src/server/data"
	"io/ioutil"
	"log"
	"net/http"
)

// GetHandler returns a data handle based on given behavior. This handle determines how a data will be handled
// for the different http methods.
func GetHandler(behavior db.Behavior) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		var needReload = float32(0)

		log.Println(fmt.Sprintf("Request received from %s with method %s", request.RemoteAddr, request.Method))

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")

		if db.CompareBehaviors(behavior, server.MainBehavior) {
			needReload += float32(0.5)
		}

		switch request.Method {

		case http.MethodGet:

			var res data.Response
			var responseData []map[string]string
			var requestPayload data.GetPayload
			var responseStatus = http.StatusOK
			var errors []string

			if isValidQueryString(request.URL.String()) {

				requestPayload = toGetPayload(request.URL.RawQuery)

				fixedFilters, unknownFilters := correctData(behavior, requestPayload.Must, "inbound")

				if unknownFilters == nil {

					var err error

					responseData, err = operation.Read(behavior.PathMapping.Table, fixedFilters)

					if err != nil {
						errors = append(errors, err.Error())
					}

				} else {

					for i := range unknownFilters {

						errors = append(errors, fmt.Sprintf("Invalid criteria: %s", unknownFilters[i]))

					}

				}

			} else {

				errors = append(errors, "Invalid payload")
				responseStatus = http.StatusBadRequest
			}

			var correctedResponse []map[string]string

			for i := range responseData {

				correctedDatum, _ := correctData(behavior, responseData[i], "outbound")

				correctedResponse = append(correctedResponse, correctedDatum)

			}

			res.Response = correctedResponse
			res.Status = responseStatus
			res.Errors = errors

			response, _ := json.Marshal(res)

			writer.WriteHeader(responseStatus)
			_, _ = writer.Write(response)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))
			log.Println(fmt.Sprintf("Errors: %s", res.Errors))

		case http.MethodPost:

			var res data.Response
			var requestPayload data.PostPayload
			var responseStatus = http.StatusOK
			var errors []string

			writer.Header().Set("Content-Type", "application/json")

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			log.Println(fmt.Sprintf("Payload: %s", rawReqPayload))

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("POST", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedData, unknownKeys := correctData(behavior, requestPayload.Set, "inbound")

					if unknownKeys == nil {

						err = operation.Create(behavior.PathMapping.Table, fixedData)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						} else {
							needReload += float32(0.5)
						}

					} else {

						for i := range unknownKeys {

							errors = append(errors, fmt.Sprintf("Invalid criteria: %s", unknownKeys[i]))

						}

					}

				} else {

					errors = append(errors, "Invalid payload")
					responseStatus = http.StatusBadRequest
				}

			}

			var correctedResponse []map[string]string

			res.Response = correctedResponse
			res.Status = responseStatus
			res.Errors = errors

			response, err := json.Marshal(res)

			writer.WriteHeader(res.Status)
			_, _ = writer.Write(response)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))
			log.Println(fmt.Sprintf("Errors: %s", res.Errors))

		case http.MethodPut:

			var res data.Response
			var requestPayload data.PutPayload
			var responseStatus = http.StatusOK
			var errors []string

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			log.Println(fmt.Sprintf("Payload: %s", rawReqPayload))

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("PUT", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedMustData, unknownKeys := correctData(behavior, requestPayload.Must, "inbound")
					fixedSetData, unknownKeys := correctData(behavior, requestPayload.Set, "inbound")

					if unknownKeys == nil {

						err = operation.Update(behavior.PathMapping.Table, fixedMustData, fixedSetData)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						} else {
							needReload += float32(0.5)
						}

					} else {

						for i := range unknownKeys {

							errors = append(errors, fmt.Sprintf("Invalid criteria: %s", unknownKeys[i]))

						}

					}

				} else {

					errors = append(errors, "Invalid payload")
					responseStatus = http.StatusBadRequest
				}

			}

			var correctedResponse []map[string]string

			res.Response = correctedResponse
			res.Status = responseStatus
			res.Errors = errors

			response, err := json.Marshal(res)

			writer.WriteHeader(res.Status)
			_, _ = writer.Write(response)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))
			log.Println(fmt.Sprintf("Errors: %s", res.Errors))

		case http.MethodDelete:
			var res data.Response
			var requestPayload data.DeletePayload
			var responseStatus = http.StatusOK
			var errors []string

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			log.Println(fmt.Sprintf("Payload: %s", rawReqPayload))

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("DELETE", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedMustData, unknownKeys := correctData(behavior, requestPayload.Must, "inbound")

					if unknownKeys == nil {

						err = operation.Delete(behavior.PathMapping.Table, fixedMustData)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						} else {
							needReload += float32(0.5)
						}

					} else {

						for i := range unknownKeys {

							errors = append(errors, fmt.Sprintf("Invalid criteria: %s", unknownKeys[i]))

						}

						responseStatus = http.StatusBadRequest

					}

				} else {

					errors = append(errors, "Invalid payload")
					responseStatus = http.StatusBadRequest
				}

			}

			var correctedResponse []map[string]string

			res.Response = correctedResponse
			res.Status = responseStatus
			res.Errors = errors

			response, err := json.Marshal(res)
			writer.WriteHeader(res.Status)
			_, _ = writer.Write(response)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))
			log.Println(fmt.Sprintf("Errors: %s", res.Errors))

		case http.MethodHead:

			var res data.Response
			var requestPayload data.GetPayload
			var responseStatus = http.StatusOK
			var errors []string

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			log.Println(fmt.Sprintf("Payload: %s", rawReqPayload))

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("GET", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedFilters, unknownFilters := correctData(behavior, requestPayload.Must, "inbound")

					if unknownFilters == nil {

						_, err = operation.Read(behavior.PathMapping.Table, fixedFilters)

						if err != nil {
							errors = append(errors, err.Error())
						}

					} else {

						for i := range unknownFilters {

							errors = append(errors, fmt.Sprintf("Invalid criteria: %s", unknownFilters[i]))

						}

					}

				} else {

					errors = append(errors, "Invalid payload")
					responseStatus = http.StatusBadRequest
				}

			}

			res.Status = responseStatus

			writer.WriteHeader(res.Status)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))
			log.Println(fmt.Sprintf("Errors: %s", res.Errors))

		case http.MethodOptions:

			var res data.Response
			var responseStatus = http.StatusOK

			res.Status = responseStatus

			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
			writer.WriteHeader(res.Status)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}

		if needReload == 1 {
			defer server.ReloadServer()
		}

	}
}
