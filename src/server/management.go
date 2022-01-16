package server

import (
	"context"
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

	for _, behavior := range behaviors {

		serverMux.HandleFunc(behavior.PathMapping.Path, GetHandler(behavior))

	}

	server = http.Server{
		Addr: Conf.Listener.Production.Address+":"+Conf.Listener.Production.Port,
		Handler: serverMux}

}

// Assigns a function to each endpoint of implemented configuration functionalities
func prepareConfigServer() {

	configServerMux = http.NewServeMux()

	for i := range implementedFunctionalities {

		configServerMux.HandleFunc(implementedFunctionalities[i], GetConfigHandler(implementedFunctionalities[i], reloadChannel))

	}

	configServer = http.Server{
		Addr: Conf.Listener.Management.Address+":"+Conf.Listener.Management.Port,
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

	log.Println("Listen requests to user-defined endpoints in port 80")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
	}

}

func startConfigServer() {

	log.Println("Configuring config server")

	prepareConfigServer()

	log.Println("\t└──Success!")

	go func() {

		log.Println("Listen requests to configuration endpoints in port 9090")
		log.Fatal(configServer.ListenAndServe())

	}()

}

func listen() {

	log.Println("Starting server")

	prepareServer()
	go startServer()

	needReload := <-reloadChannel

	if needReload {

		log.Println("Stopping server...")

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
