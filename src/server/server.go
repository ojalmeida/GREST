package server

import (
	"github.com/ojalmeida/GREST/src/db"
	"log"
	"net/http"
)

var behaviors []db.Behavior

var server http.Server
var serverMux *http.ServeMux

var configServerMux *http.ServeMux
var configServer http.Server

// Checks the health of the application and loads all the data necessary to start the server
func init() {

	checkHealth()

	var err error
	behaviors, err = db.GetBehaviors()

	if err != nil {

		panic(err.Error())

	}

}

// StartServers applies all behaviors and starts to listen for requests
func StartServers() {

	setServerMux()
	setConfigServerMux()

	server = http.Server{Addr: ":8080", Handler: serverMux}
	configServer = http.Server{Addr: ":9090", Handler: configServerMux}

	go log.Fatal(server.ListenAndServe())

	go log.Fatal(configServer.ListenAndServe())

}
