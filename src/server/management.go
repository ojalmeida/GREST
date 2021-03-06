package server

import (
	"context"
	"github.com/ojalmeida/GREST/src/config"
	"github.com/ojalmeida/GREST/src/db"
	"log"
	"net/http"
)

var implementedFunctionalities []string

var reloadChannel = make(chan bool)

func init() {

	implementedFunctionalities = append(implementedFunctionalities,
		"/config/behaviors",
		"/config/path-mappings",
		"/config/key-mappings")

}

// Transform each behavior in a function that process requests
func prepareServer() {

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

	serverMux = http.NewServeMux()

	zombieBehavior := db.Behavior{
		PathMapping: db.PathMapping{},
		KeyMappings: nil,
	}

	for _, behavior := range behaviors {

		if db.ComparePathMappings(behavior.PathMapping, zombieBehavior.PathMapping) || behavior.KeyMappings == nil {

			continue

		}

		serverMux.HandleFunc(behavior.PathMapping.Path, GetHandler(behavior))

	}

	server = http.Server{
		Addr:    config.Conf.API.Production.Address + ":" + config.Conf.API.Production.Port,
		Handler: serverMux}

}

// Assigns a function to each endpoint of implemented configuration functionalities
func prepareConfigServer() {

	configServerMux = http.NewServeMux()

	for i := range implementedFunctionalities {

		configServerMux.HandleFunc(implementedFunctionalities[i], GetConfigHandler(implementedFunctionalities[i], reloadChannel))

	}

	configServer = http.Server{
		Addr:    config.Conf.API.Management.Address + ":" + config.Conf.API.Management.Port,
		Handler: configServerMux}

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

// StartServers applies all behaviors and starts to listen for requests
func StartServers() {
	log.Println("Starting servers")
	go listen()

	go listenConfig()

}

func startServer() {

	log.Println("Listen requests to user-defined endpoints in port " + config.Conf.API.Production.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
	}

}

func startConfigServer() {

	log.Println("Configuring management server")

	prepareConfigServer()

	log.Println("\t└──Success!")

	go func() {

		log.Println("Listen requests to configuration endpoints in port " + config.Conf.API.Management.Port)
		log.Fatal(configServer.ListenAndServe())

	}()

}

func listen() {

	log.Println("Starting server")

	prepareServer()
	go startServer()

	needReload := <-reloadChannel

	if needReload {

		log.Println("Behaviors change detected, stopping server...")

		err := server.Shutdown(context.Background())

		if err != nil {
			log.Println("\t└──Fail!")
		} else {

			log.Println("\t└──Success!")
			listen()
		}

	}

}

func listenConfig() {
	prepareConfigServer()
	startConfigServer()

}
