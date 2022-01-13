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

	configServerMux = http.NewServeMux()

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

	log.Println("Checking health of config database")

	ok := db.CheckLocalDB()

	if !ok {

		log.Println("\t├──Not ok")
		log.Println("\t└──Trying to self-healing")

		err := db.CreateTables()

		if err != nil {
			log.Println("\t\t└──Fail!")
			panic(err.Error())
		} else {
			log.Println("\t\t└──Success!")
		}

	}

	log.Println("Health ok")

	var err error

	if err != nil {
		panic(err.Error())
	}

}
