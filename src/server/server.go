package server

import (
	"github.com/ojalmeida/GREST/src/db"
	"log"
	"net/http"
	"sync"
)

var behaviors []db.Behavior

var server http.Server
var serverMux *http.ServeMux

var configServerMux *http.ServeMux
var configServer http.Server

// Checks the health of the application and loads all the data necessary to start the server
func init() {

	checkHealth()

	log.Println("Getting user-defined endpoints")

	var err error
	behaviors, err = db.GetBehaviors()

	if err != nil {

		log.Println("\t└──Fail!")
		panic(err.Error())

	} else {

		log.Println("\t└──Success!")
	}

}

// StartServers applies all behaviors and starts to listen for requests
func StartServers() {

	log.Println("Configuring servers")

	setServerMux()
	setConfigServerMux()

	log.Println("\t└──Success!")

	wg := new(sync.WaitGroup)

	wg.Add(2)

	server = http.Server{Addr: ":8080", Handler: serverMux}
	configServer = http.Server{Addr: ":9090", Handler: configServerMux}

	log.Println("Starting servers")

	go func() {

		log.Println("Listen requests to configuration endpoints in port 9090")
		log.Fatal(configServer.ListenAndServe())
		wg.Done()

	}()

	go func() {

		log.Println("Listen requests to user-defined endpoints in port 8080")
		log.Fatal(server.ListenAndServe())
		wg.Done()

	}()

	wg.Wait()

}
