package server

import (
	"github.com/ojalmeida/GREST/src/db"
	"log"
	"net/http"
)

// Transform each behavior in a function that process requests
func setServerMuxes() {

	serverMux = http.NewServeMux()

	for _, behavior := range behaviors {

		serverMux.HandleFunc(behavior.PathMapping.Path, GetHandler(behavior))

	}

	server.Handler = serverMux

}

// ReloadServer reloads data handlers
func ReloadServer() {

	checkHealth()
	setServerMuxes()

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
