package server

import (
	"context"
	"github.com/ojalmeida/GREST/src/config"
	"github.com/ojalmeida/GREST/src/db"
	"log"
	"net/http"
	"time"
)

var implementedFunctionalities []string

var reloadChannel = make(chan bool)

var behaviorMode = "api"

func init() {

	implementedFunctionalities = append(implementedFunctionalities,
		"/config/behaviors",
		"/config/path-mappings",
		"/config/key-mappings")

}

// Transform each behavior in a function that process requests
func prepareServer() {

	checkHealth()

	if behaviorMode == "api" {

		log.Println("Getting endpoints defined via API")
		var err error

		// behaviors got by database query
		behaviors, err = db.GetBehaviors()

		if err != nil {

			log.Println("\t└── Fail!")
			panic(err.Error())

		} else {

			log.Println("\t└── Success!")
		}

	} else if behaviorMode == "file" {

		log.Println("Getting declared endpoints")
		var errs []error

		// behaviors declared in file
		unverifiedBehaviors, errs := getDeclaredBehaviors()

		if len(unverifiedBehaviors) > 1 {

			for i, pivotBehavior := range unverifiedBehaviors {

				duplicatedBehavior := false
				duplicatedPath := false
				var duplicatedPathSlice = []db.Behavior{pivotBehavior}

				for j, behavior := range unverifiedBehaviors {

					if i == j {
						continue
					}

					if pivotBehavior.PathMapping.Path == behavior.PathMapping.Path {

						duplicatedPathSlice = append(duplicatedPathSlice, behavior)
						duplicatedPath = true

					}

					if db.CompareBehaviors(pivotBehavior, behavior) {

						duplicatedBehavior = true
					}

				}

				if len(duplicatedPathSlice) != 1 {

					log.Println("\t└── duplicated path \"" + pivotBehavior.PathMapping.Path + "\" detected, skipping")

				}

				if !duplicatedPath || !duplicatedBehavior {

					behaviors = append(behaviors, pivotBehavior)

				}

			}

		} else {

			behaviors = unverifiedBehaviors

		}

		if errs != nil {

			for i := range errs {

				log.Println("\t└── " + errs[i].Error())

			}

		}

	}

	serverMux = http.NewServeMux()

	zombieBehavior := db.Behavior{
		PathMapping: db.PathMapping{},
		KeyMappings: nil,
	}

	for _, behavior := range behaviors {

		// checks zombie behaviors
		if db.ComparePathMappings(behavior.PathMapping, zombieBehavior.PathMapping) || behavior.KeyMappings == nil {

			log.Println("\t└── Zombie behavior ignored")
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

	if behaviorMode == "file" {

		time.Sleep(time.Second * 5)
		go listenToDeclarationChanges()
	}

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

// Serves the configured API
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

// Serves the config API
func listenConfig() {
	prepareConfigServer()
	startConfigServer()

}

// Detects changes in behavior declaration files and triggers server reload
func listenToDeclarationChanges() {

	var previousDeclaredBehaviors []db.Behavior

	for {

		time.Sleep(time.Second)

		declaredBehaviors, _ := getDeclaredBehaviors()

		declaredBehaviorsChange := false

		// Compare content of behaviors on slices
		for i := range declaredBehaviors {

			alreadyExists := false

			toBeInsertedBehavior := declaredBehaviors[i]

			for j := range previousDeclaredBehaviors {

				if behaviors != nil {

					alreadyInsertedBehavior := behaviors[j]

					if db.CompareBehaviors(alreadyInsertedBehavior, toBeInsertedBehavior) {

						alreadyExists = true

						break
					}

				}

			}

			if !alreadyExists {

				declaredBehaviorsChange = true

				break

			}

		}

		// Compare number of behaviors on slices
		if !declaredBehaviorsChange {

			declaredBehaviorsChange = len(declaredBehaviors) != len(previousDeclaredBehaviors)

		}

		if declaredBehaviorsChange && previousDeclaredBehaviors != nil {

			previousDeclaredBehaviors = declaredBehaviors
			reloadChannel <- true

		} else if previousDeclaredBehaviors == nil {

			previousDeclaredBehaviors = declaredBehaviors

		}
	}

}
