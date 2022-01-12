package server

import (
	"github.com/ojalmeida/GREST/src/db"
	"log"
	"net/http"
)

var implementedFunctionalities []string

func init() {
	implementedFunctionalities = append(implementedFunctionalities,
		"/config/behaviors",
		"/config/path-mappings",
		"/config/key-mappings")
}

// Transform each behavior in a function that process requests
func setServerMux() {

	serverMux = http.NewServeMux()

	for _, behavior := range behaviors {

		serverMux.HandleFunc(behavior.PathMapping.Path, GetHandler(behavior))

	}

}

// Assigns a function to each endpoint of implemented configuration functionalities
func setConfigServerMux() {

	for i := range implementedFunctionalities {

		configServerMux.HandleFunc(implementedFunctionalities[i], GetConfigHandler(implementedFunctionalities[i]))

	}

}

// ReloadServer reloads data handlers
func ReloadServer() {

	checkHealth()
	setServerMux()

}

// ReloadConfigServer reloads handlers of implemented configuration endpoints
func ReloadConfigServer() {

	checkHealth()
	setConfigServerMux()

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

	if err != nil {
		panic(err.Error())
	}

}
