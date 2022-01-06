package server

import (
	"encoding/json"
	"fmt"
	"grest/db"
	"io/ioutil"
	"log"
	"net/http"
)

var behaviors []db.Behavior

var mainBehavior db.Behavior

var server http.Server
var serverMux *http.ServeMux

type response struct {
	Status   int         `json:"status"`
	Response interface{} `json:"response"`
	Errors   []string    `json:"errors"`
}

type GetPayload struct {
	Must map[string]string `json:"match"`
}

type PostPayload struct {
	Set map[string]string
}

type PutPayload struct {
	Must map[string]string
	Set  map[string]string
}

type DeletePayload struct {
	Must map[string]string
}

// Checks the health of the application and loads all the data necessary to start the server
func init() {

	checkHealth()

	mainKeyMappings := []db.KeyMapping{
		{Key: "id", Column: "behavior_id"},
		{Key: "path_mapping_id", Column: "path_mapping_id"},
		{Key: "key_mapping_id", Column: "key_mapping_id"},
	}
	mainPathMapping := db.PathMapping{Path: "/config/behaviors", Table: "behaviors"}
	mainBehavior.PathMapping = mainPathMapping
	mainBehavior.KeyMappings = mainKeyMappings

}

// Reload request handlers
func reloadServer() {

	checkHealth()
	defineHandlers()

}

// Checks pre-requisites to start server
func checkHealth() {

	ok, missingPathMappings, missingKeyMappings, missingBehaviors := db.CheckConfigs()

	if !ok {

		log.Printf("Missing following PathMappings: %s", missingPathMappings)
		log.Printf("Missing following KeyMappings: %s", missingKeyMappings)
		log.Printf("Missing following Behaviors: %s", missingBehaviors)
		log.Println("Trying to insert default configs")

		err := db.PopulateConfigs()

		if err != nil {
			log.Println("Fail!")
			panic(err.Error())
		} else {
			log.Println("Success!")
		}

	}

	var err error

	behaviors, err = db.GetBehaviors()

	if err != nil {
		panic(err.Error())
	}

}

// StartServer applies all behaviors and starts to listen for requests
func StartServer() {

	defineHandlers()

	server = http.Server{Addr: ":8080", Handler: serverMux}

	log.Fatal(server.ListenAndServe())

}

// Transform each behavior in a function that process requests
func defineHandlers() {

	serverMux = http.NewServeMux()

	for _, behavior := range behaviors {

		setHandler(behavior, serverMux)

	}

	server.Handler = serverMux

}

// setHandler applies a request handling based on given behavior object.
func setHandler(behavior db.Behavior, serverMux *http.ServeMux) {

	serverMux.HandleFunc(behavior.PathMapping.Path, handlerFactory(behavior))
}

// handlerFactory returns a request handler based on given behavior. This handler determines how a request will be handled
// for the different http methods.
func handlerFactory(behavior db.Behavior) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		var needReload = float32(0)

		log.Println(fmt.Sprintf("Request received from %s with method %s", request.RemoteAddr, request.Method))

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")

		if db.CompareBehaviors(behavior, mainBehavior) {
			needReload += float32(0.5)
		}

		switch request.Method {

		case http.MethodGet:

			var res response
			var responseData []map[string]string
			var requestPayload GetPayload
			var responseStatus = http.StatusOK
			var errors []string

			if isValidQueryString(request.URL.String()) {

				requestPayload = convertQueryStringToGETPayload(request.URL.RawQuery)

				fixedFilters, unknownFilters := correctData(behavior, requestPayload.Must, "inbound")

				if unknownFilters == nil {

					var err error

					responseData, err = db.Read(behavior.PathMapping.Table, fixedFilters)

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

			var res response
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

						err = db.Create(behavior.PathMapping.Table, fixedData)

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

			var res response
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

						err = db.Update(behavior.PathMapping.Table, fixedMustData, fixedSetData)

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
			var res response
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

						err = db.Delete(behavior.PathMapping.Table, fixedMustData)

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

			var res response
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

						_, err = db.Read(behavior.PathMapping.Table, fixedFilters)

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

			var res response
			var responseStatus = http.StatusOK

			res.Status = responseStatus

			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
			writer.WriteHeader(res.Status)

			log.Println(fmt.Sprintf("Response status: %d", res.Status))

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}

		if needReload == 1 {
			defer reloadServer()
		}

	}
}
