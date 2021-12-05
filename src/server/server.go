package server

import (
	"fmt"
	"grest/db"
	"io"
	"log"
	"net/http"
)

var apiDefaultHandler = func(writer http.ResponseWriter, request *http.Request) {

	_, err := io.WriteString(writer, "Testing...")

	if err != nil {
		panic(err.Error())
	}

}

var behaviors []db.Behavior

func init() {

	behaviors = db.GetBehaviors()

}

func StartServer() {

	http.HandleFunc("/", apiDefaultHandler)

	for _, behavior := range behaviors {

		setHandler(behavior)

	}

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func setHandler(behavior db.Behavior) {

	http.HandleFunc(behavior.PathMapping.Path, handlerFactory(behavior))
}

func handlerFactory(behavior db.Behavior) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		_, _ = io.WriteString(writer, fmt.Sprintf("You entered in the %s endpoint", behavior.PathMapping.Table))

	}
}
