package server

import (
	"encoding/json"
	"fmt"
	"github.com/ojalmeida/GREST/src/db"
	"io/ioutil"
	"log"
	"net/http"
)

// GetHandler returns a data handle based on given behavior. This handle determines how a data will be handled
// for the different http methods.
func GetHandler(behavior db.Behavior) func(writer http.ResponseWriter, request *http.Request) {

	var dbName = "mysql"

	return func(writer http.ResponseWriter, request *http.Request) {

		log.Println(fmt.Sprintf("Request received from %s with method %s", request.RemoteAddr, request.Method))

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")

		switch request.Method {

		case http.MethodGet:

			var res Response
			var responseData []map[string]string
			var requestPayload GetPayload
			var responseStatus = http.StatusOK
			var errors []string

			if isValidQueryString(request.URL.String()) {

				requestPayload = toGetPayload(request.URL.RawQuery)

				fixedFilters, unknownFilters := correctData(behavior, requestPayload.Must, "inbound")

				if unknownFilters == nil {

					var err error

					responseData, err = db.Read(behavior.PathMapping.Table, fixedFilters, dbName)

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

			var res Response
			var requestPayload PostPayload
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

						err = db.Create(behavior.PathMapping.Table, fixedData, dbName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
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

			var res Response
			var requestPayload PutPayload
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

						err = db.Update(behavior.PathMapping.Table, fixedMustData, fixedSetData, dbName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
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
			var res Response
			var requestPayload DeletePayload
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

						err = db.Delete(behavior.PathMapping.Table, fixedMustData, dbName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
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

			var res Response
			var requestPayload GetPayload
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

						_, err = db.Read(behavior.PathMapping.Table, fixedFilters, dbName)

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

			var res Response
			var responseStatus = http.StatusOK

			res.Status = responseStatus

			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
			writer.WriteHeader(res.Status)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
}

func GetConfigHandler(endpoint string, reload chan bool) func(writer http.ResponseWriter, request *http.Request) {

	driverName := "sqlite3-config"

	switch endpoint {

	case "/config/behaviors":

		return func(writer http.ResponseWriter, request *http.Request) {

			var needReload float32 = 0
			tableName := "behavior"

			switch request.Method {

			case http.MethodGet:

				var res Response
				var responseData []map[string]string
				var requestPayload GetPayload
				var responseStatus = http.StatusOK
				var errors []string

				if isValidQueryString(request.URL.String()) {

					requestPayload = toGetPayload(request.URL.RawQuery)

					var err error

					responseData, err = db.Read(tableName, requestPayload.Must, driverName)

					if err != nil {
						errors = append(errors, err.Error())
					}

				} else {

					errors = append(errors, "Invalid payload")
					responseStatus = http.StatusBadRequest
				}

				res.Response = responseData
				res.Status = responseStatus
				res.Errors = errors

				response, _ := json.Marshal(res)

				writer.WriteHeader(responseStatus)
				_, _ = writer.Write(response)

				log.Println(fmt.Sprintf("Response status: %d", res.Status))
				log.Println(fmt.Sprintf("Errors: %s", res.Errors))

			case http.MethodPost: // If Method is Post.

				var res Response
				var requestPayload PostPayload
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

						err = db.Create(tableName, requestPayload.Set, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						}

						needReload += float32(1)

					} else {

						errors = append(errors, "Invalid payload")
						responseStatus = http.StatusBadRequest
					}

				}

				res.Response = nil
				res.Status = responseStatus
				res.Errors = errors

				response, err := json.Marshal(res)

				writer.WriteHeader(res.Status)
				_, _ = writer.Write(response)

				log.Println(fmt.Sprintf("Response status: %d", res.Status))
				log.Println(fmt.Sprintf("Errors: %s", res.Errors))

			case http.MethodPut: // If Method is Put

				var res Response
				var requestPayload PutPayload
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

						err = db.Update(tableName, requestPayload.Must, requestPayload.Set, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						}

						needReload += float32(1)

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
				var res Response
				var requestPayload DeletePayload
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

						err = db.Delete(tableName, requestPayload.Must, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						}

						needReload += float32(1)

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

				var res Response
				var requestPayload GetPayload
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

						_, err = db.Read(tableName, requestPayload.Must, driverName)

						if err != nil {
							errors = append(errors, err.Error())
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

			}

			if needReload == 1 {
				defer func() { reload <- true }()
			}

		}

	case "/config/path-mappings":

		tableName := "path_mapping"

		return func(writer http.ResponseWriter, request *http.Request) {

			// In the end of processing, if equals 1, reload serverMux
			var needReload = float32(0)

			switch request.Method {

			case http.MethodGet:

				var res Response
				var responseData []map[string]string
				var requestPayload GetPayload
				var responseStatus = http.StatusOK
				var errors []string

				if isValidQueryString(request.URL.String()) {

					requestPayload = toGetPayload(request.URL.RawQuery)

					var err error

					responseData, err = db.Read(tableName, requestPayload.Must, "sqlite3-config")

					if err != nil {
						errors = append(errors, err.Error())
					}

				} else {

					errors = append(errors, "Invalid payload")
					responseStatus = http.StatusBadRequest
				}

				res.Response = responseData
				res.Status = responseStatus
				res.Errors = errors

				response, _ := json.Marshal(res)

				writer.WriteHeader(responseStatus)
				_, _ = writer.Write(response)

				log.Println(fmt.Sprintf("Response status: %d", res.Status))
				log.Println(fmt.Sprintf("Errors: %s", res.Errors))

			case http.MethodPost:

				var res Response
				var requestPayload PostPayload
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

						err = db.Create(tableName, requestPayload.Set, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						}

					} else {

						errors = append(errors, "Invalid payload")
						responseStatus = http.StatusBadRequest
					}

				}

				res.Response = nil
				res.Status = responseStatus
				res.Errors = errors

				response, err := json.Marshal(res)

				writer.WriteHeader(res.Status)
				_, _ = writer.Write(response)

				log.Println(fmt.Sprintf("Response status: %d", res.Status))
				log.Println(fmt.Sprintf("Errors: %s", res.Errors))

			case http.MethodPut:

				var res Response
				var requestPayload PutPayload
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

						err = db.Update(tableName, requestPayload.Must, requestPayload.Set, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
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
				var res Response
				var requestPayload DeletePayload
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

						err = db.Delete(tableName, requestPayload.Must, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
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

				var res Response
				var requestPayload GetPayload
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

						_, err = db.Read(tableName, requestPayload.Must, driverName)

						if err != nil {
							errors = append(errors, err.Error())
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

			}

			if needReload == 1 {
				defer func() { reload <- true }()
			}

		}

	case "/config/key-mappings":

		tableName := "key_mapping"

		return func(writer http.ResponseWriter, request *http.Request) {

			// In the end of processing, if equals 1, reload serverMux
			var needReload = float32(0)

			switch request.Method {

			case http.MethodGet:

				var res Response
				var responseData []map[string]string
				var requestPayload GetPayload
				var responseStatus = http.StatusOK
				var errors []string

				if isValidQueryString(request.URL.String()) {

					requestPayload = toGetPayload(request.URL.RawQuery)

					var err error

					responseData, err = db.Read(tableName, requestPayload.Must, "sqlite3-config")

					if err != nil {
						errors = append(errors, err.Error())
					}

				} else {

					errors = append(errors, "Invalid payload")
					responseStatus = http.StatusBadRequest
				}

				res.Response = responseData
				res.Status = responseStatus
				res.Errors = errors

				response, _ := json.Marshal(res)

				writer.WriteHeader(responseStatus)
				_, _ = writer.Write(response)

				log.Println(fmt.Sprintf("Response status: %d", res.Status))
				log.Println(fmt.Sprintf("Errors: %s", res.Errors))

			case http.MethodPost:

				var res Response
				var requestPayload PostPayload
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

						err = db.Create(tableName, requestPayload.Set, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
						}

					} else {

						errors = append(errors, "Invalid payload")
						responseStatus = http.StatusBadRequest
					}

				}

				res.Response = nil
				res.Status = responseStatus
				res.Errors = errors

				response, err := json.Marshal(res)

				writer.WriteHeader(res.Status)
				_, _ = writer.Write(response)

				log.Println(fmt.Sprintf("Response status: %d", res.Status))
				log.Println(fmt.Sprintf("Errors: %s", res.Errors))

			case http.MethodPut:

				var res Response
				var requestPayload PutPayload
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

						err = db.Update(tableName, requestPayload.Must, requestPayload.Set, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
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
				var res Response
				var requestPayload DeletePayload
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

						err = db.Delete(tableName, requestPayload.Must, driverName)

						if err != nil {
							errors = append(errors, err.Error())
							responseStatus = http.StatusInternalServerError
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

				var res Response
				var requestPayload GetPayload
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

						_, err = db.Read(tableName, requestPayload.Must, driverName)

						if err != nil {
							errors = append(errors, err.Error())
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

			}

			if needReload == 1 {
				defer func() { reload <- true }()
			}

		}

	case "/config/auth":

		return func(writer http.ResponseWriter, request *http.Request) {

			switch request.Method {

			case http.MethodPost:

				// get token

			case http.MethodPut:

				// change token

			case http.MethodDelete:

				// expire token

			case http.MethodOptions:

				// methods allowed

			}

		}

	case "/config/rate-limit":

		return func(writer http.ResponseWriter, request *http.Request) {

			// In the end of processing, if equals 1, reload serverMux
			var needReload = float32(0)

			switch request.Method {

			case http.MethodGet:

			case http.MethodPost:

			case http.MethodPut:

			case http.MethodDelete:

			case http.MethodOptions:

			case http.MethodHead:

			}

			if needReload == 1 {
			}

		}

	default:

		// Default, if not matched with any endpoint
		return func(writer http.ResponseWriter, request *http.Request) {

			writer.WriteHeader(http.StatusNotFound)

		}

	}

}
