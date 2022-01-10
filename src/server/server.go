package server

import (
	"github.com/ojalmeida/GREST/src/db"
	"log"
	"net/http"
)

var behaviors []db.Behavior

var MainBehavior db.Behavior

var server http.Server
var serverMux *http.ServeMux

// Checks the health of the application and loads all the data necessary to start the server
func init() {

	checkHealth()

	mainKeyMappings := []db.KeyMapping{
		{Key: "id", Column: "behavior_id"},
		{Key: "path_mapping_id", Column: "path_mapping_id"},
		{Key: "key_mapping_id", Column: "key_mapping_id"},
	}
	mainPathMapping := db.PathMapping{Path: "/config/behaviors", Table: "behaviors"}
	MainBehavior.PathMapping = mainPathMapping
	MainBehavior.KeyMappings = mainKeyMappings

}

// StartServer applies all behaviors and starts to listen for requests
func StartServer() {

	setServerMuxes()

	server = http.Server{Addr: ":8080", Handler: serverMux}

	log.Fatal(server.ListenAndServe())

}
