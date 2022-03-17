package server

import (
	"context"
	"github.com/ojalmeida/GREST/src/config"
	"github.com/ojalmeida/GREST/src/db"
	log "github.com/ojalmeida/GREST/src/log"
	builtinLog "log"
	"net/http"
	"sync"
	"time"
)

var (
	reloadMainServerChannel         chan bool
	stopMainServerChannel           chan bool
	reloadConfigServerChannel       chan bool
	stopConfigServerChannel         chan bool
	stopDeclarationListeningChannel chan bool
	reloadAppChannel                chan bool
)

func Start() {

	builtinLog.Println("Starting GREST")

	wg := &sync.WaitGroup{}

	wg.Add(1)
	config.Start(&reloadAppChannel, wg)

	log.Start()

	initVariables()
	initDatabases()

	switch declarationMode {

	case "api":

		wg.Add(1)
		go startConfigServer(wg)

	case "file":

		wg.Add(1)
		go listenToDeclarationChanges(wg)

	}

	wg.Add(1)
	go startMainServer(wg)

	<-reloadAppChannel // blocks execution

	log.InfoLogger.Println("Application reload requested")

	log.InfoLogger.Println("Reloading application")

	stopMainServerChannel <- true

	switch declarationMode {

	case "api":

		stopConfigServerChannel <- true

	case "file":

		stopDeclarationListeningChannel <- true

	}

	closeChannels(
		&reloadAppChannel,
		&reloadMainServerChannel,
		&reloadConfigServerChannel,
		&stopMainServerChannel,
		&stopConfigServerChannel,
		&stopDeclarationListeningChannel,
	)

	wg.Wait()

}

func initVariables() {

	initChannels(
		&reloadAppChannel,
		&reloadMainServerChannel,
		&reloadConfigServerChannel,
		&stopMainServerChannel,
		&stopConfigServerChannel,
		&stopDeclarationListeningChannel,
	)

	time.Sleep(time.Millisecond * 500)

	autoReload = config.Conf.Configuration.AutoReload
	declarationMode = config.Conf.Configuration.DeclarationMode

}

func initDatabases() {

	db.ConnectToRemoteDB()
	if declarationMode == "api" {

		db.ConnectToLocalDB()

	}
}

func prepareMainServer() {

	log.InfoLogger.Println("Configuring server")

	if declarationMode == "api" {

		log.InfoLogger.Println("Getting endpoints defined via API")
		var err error

		// behaviors got by database query
		behaviors, err = db.GetBehaviors()

		if err != nil {

			log.ErrorLogger.Println("Fail on getting endpoints defined via API")
			panic(err.Error())

		} else {

			log.InfoLogger.Println("Success on getting endpoints defined via API")
		}

	} else if declarationMode == "file" {

		log.InfoLogger.Println("Getting endpoints declared via file")
		var errs []error

		// behaviors declared in file
		unverifiedBehaviors, errs := db.GetDeclaredBehaviors()

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

					log.WarningLogger.Println("Duplicated path \"" + pivotBehavior.PathMapping.Path + "\" detected, skipping")

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

				log.ErrorLogger.Println(errs[i].Error())

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

			log.WarningLogger.Println("Zombie behavior ignored")
			continue

		}

		serverMux.HandleFunc(behavior.PathMapping.Path, GetHandler(behavior))

	}

	server = http.Server{
		Addr:    config.Conf.API.Production.Address + ":" + config.Conf.API.Production.Port,
		Handler: serverMux}

	log.InfoLogger.Println("Success on configuring server")

}

func prepareConfigServer() {

	log.InfoLogger.Println("Configuring management server")

	checkHealth()

	configServerMux = http.NewServeMux()

	for i := range implementedFunctionalities {

		configServerMux.HandleFunc(implementedFunctionalities[i], GetConfigHandler(implementedFunctionalities[i]))

	}

	configServer = http.Server{
		Addr:    config.Conf.API.Management.Address + ":" + config.Conf.API.Management.Port,
		Handler: configServerMux,
	}

	log.InfoLogger.Println("Success on configuring management server")

}

func checkHealth() {

	log.InfoLogger.Println("Checking health of config database")

	ok := db.CheckLocalDB()

	if !ok {

		log.WarningLogger.Println("Config database health not ok")
		log.WarningLogger.Println("Trying to heal config database")

		err := db.CreateTables()

		if err != nil {
			log.ErrorLogger.Println("Fail on heal config database")
			panic(err.Error())
		} else {
			log.WarningLogger.Println("Success on heal config database")
		}

	}

	log.InfoLogger.Println("Config database health ok")

	var err error

	if err != nil {
		panic(err.Error())
	}

}

func startMainServer(wg *sync.WaitGroup) {

	prepareMainServer()

	log.InfoLogger.Println("Starting server")

	go func() {
		log.InfoLogger.Println("Listening requests to user-defined endpoints in port " + config.Conf.API.Production.Port)

		err := server.ListenAndServe()

		if err != http.ErrServerClosed {

			log.ErrorLogger.Println(err.Error())

		} // blocks execution
	}()

	select {

	case <-reloadMainServerChannel:

		if autoReload {

			log.InfoLogger.Println("Server reload requested")

			log.InfoLogger.Println("Reloading server")
			err := server.Shutdown(context.Background())

			if err != nil {
				log.ErrorLogger.Println("Fail on stopping server: ", err.Error())
			} else {

				log.InfoLogger.Println("Success on stopping server")
			}

			startMainServer(wg) // recursive function

		}

	case <-stopMainServerChannel:

		log.InfoLogger.Println("Server stop requested")

		log.InfoLogger.Println("Stopping server")
		err := server.Shutdown(context.Background())

		if err != nil {
			log.ErrorLogger.Println("Fail on stopping server: ", err.Error())
		} else {

			log.InfoLogger.Println("Success on stopping server")
		}

		wg.Done()

	}

}

func startConfigServer(wg *sync.WaitGroup) {

	prepareConfigServer()

	log.InfoLogger.Println("Starting config server")

	go func() {

		log.InfoLogger.Println("Listen requests to configuration endpoints in port " + config.Conf.API.Management.Port)

		err := configServer.ListenAndServe()

		if err != http.ErrServerClosed {

			log.ErrorLogger.Println(err.Error())
		}

	}()

	select {

	case <-reloadConfigServerChannel:

		log.InfoLogger.Println("Config server reload requested")

		log.InfoLogger.Println("Reloading config server")

		err := configServer.Shutdown(context.Background())

		if err != nil {
			log.ErrorLogger.Println("Fail on stopping server: ", err.Error())
		} else {

			log.InfoLogger.Println("Success on stopping server")
		}

		startConfigServer(wg) // recursive function

	case <-stopConfigServerChannel:

		log.InfoLogger.Println("Config server stop requested")

		log.InfoLogger.Println("Stopping config server")

		err := configServer.Shutdown(context.Background())

		if err != nil {
			log.ErrorLogger.Println("Fail on stopping config server: ", err.Error())
		} else {

			log.InfoLogger.Println("Success on stopping config server")
		}

		wg.Done()

	}

}

func listenToDeclarationChanges(wg *sync.WaitGroup) {

	var previousDeclaredBehaviors []db.Behavior
	firstRun := true

loop:
	for {

		select {

		case <-stopDeclarationListeningChannel:
			break loop

		default:

			time.Sleep(time.Second)

			func() {

				if declarationMode != "file" || !autoReload {

					return

				}

				declaredBehaviors, _ := db.GetDeclaredBehaviors()

				declaredBehaviorsChange := false

				if firstRun {

					previousDeclaredBehaviors = declaredBehaviors
					firstRun = false
					return

				}

				// Compare content of behaviors on slices
				for i := range declaredBehaviors {

					alreadyExists := false

					toBeInsertedBehavior := declaredBehaviors[i]

					for j := range previousDeclaredBehaviors {

						alreadyInsertedBehavior := previousDeclaredBehaviors[j]

						if db.CompareBehaviors(alreadyInsertedBehavior, toBeInsertedBehavior) {

							alreadyExists = true

							break
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

				if declaredBehaviorsChange {

					previousDeclaredBehaviors = declaredBehaviors
					reloadMainServerChannel <- true

				} else if previousDeclaredBehaviors == nil {

					previousDeclaredBehaviors = declaredBehaviors

				}

			}() // check changes

		}

	}

	wg.Done()

}

func initChannels(channels ...*chan bool) {

	for _, channel := range channels {

		*channel = make(chan bool)

	}

}

func closeChannels(channels ...*chan bool) {

	for _, channel := range channels {

		close(*channel)

	}
}
