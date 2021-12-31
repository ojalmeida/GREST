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

type response struct {
	Status   int         `json:"status"`
	Response interface{} `json:"response"`
	Errors   []string    `json:"errors"`
}

type GetPayload struct {
	Match map[string]string `json:"match"`
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

	for _, behavior := range behaviors {

		setHandler(behavior)

	}

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// setHandler applies a request handling based on given behavior object.
func setHandler(behavior db.Behavior) {

	http.HandleFunc(behavior.PathMapping.Path, handlerFactory(behavior))
}

// handlerFactory returns a request handler based on given behavior. This handler determines how a request will be handled
// for the different http methods.
func handlerFactory(behavior db.Behavior) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		switch request.Method {

		case http.MethodGet:

			var res response
			var responseData []map[string]string
			var requestPayload GetPayload
			var responseStatus = http.StatusOK
			var errors []string

			writer.Header().Set("Content-Type", "application/json")
			writer.Header().Set("Access-Control-Allow-Origin", "*")

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("GET", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedFilters, unknownFilters := correctData(behavior, requestPayload.Match)

					if unknownFilters == nil {

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

			}

			var correctedResponse []map[string]string

			for i := range responseData {

				correctedDatum, _ := correctData(behavior, responseData[i])

				correctedResponse = append(correctedResponse, correctedDatum)

			}

			res.Response = correctedResponse
			res.Status = responseStatus
			res.Errors = errors

			response, err := json.Marshal(res)

			writer.WriteHeader(responseStatus)
			_, _ = writer.Write(response)

		case http.MethodPost:

			var res response
			var requestPayload PostPayload
			var responseStatus = http.StatusOK
			var errors []string

			writer.Header().Set("Content-Type", "application/json")

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("POST", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedData, unknownKeys := correctData(behavior, requestPayload.Set)

					if unknownKeys == nil {

						err = db.Create(behavior.PathMapping.Table, fixedData)

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

		case http.MethodPut:

			var res response
			var requestPayload PutPayload
			var responseStatus = http.StatusOK
			var errors []string

			writer.Header().Set("Content-Type", "application/json")

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("PUT", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedMustData, unknownKeys := correctData(behavior, requestPayload.Must)
					fixedSetData, unknownKeys := correctData(behavior, requestPayload.Set)

					if unknownKeys == nil {

						err = db.Update(behavior.PathMapping.Table, fixedMustData, fixedSetData)

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

		case http.MethodDelete:
			var res response
			var requestPayload DeletePayload
			var responseStatus = http.StatusOK
			var errors []string

			writer.Header().Set("Content-Type", "application/json")

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("DELETE", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedMustData, unknownKeys := correctData(behavior, requestPayload.Must)

					if unknownKeys == nil {

						err = db.Delete(behavior.PathMapping.Table, fixedMustData)

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

		case http.MethodHead:

			var res response
			var requestPayload GetPayload
			var responseStatus = http.StatusOK
			var errors []string

			writer.Header().Set("Content-Type", "application/json")

			rawReqPayload, err := ioutil.ReadAll(request.Body)

			if err != nil {
				errors = append(errors, "Impossible to read body")
				responseStatus = http.StatusInternalServerError
			} else {

				if isValidPayload("GET", string(rawReqPayload)) {

					_ = json.Unmarshal(rawReqPayload, &requestPayload)

					fixedFilters, unknownFilters := correctData(behavior, requestPayload.Match)

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

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
}
