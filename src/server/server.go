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

var pathMappings []db.PathMapping

func init() {

	pathMappings = db.GetPathMappings()

}

func StartServer() {

	http.HandleFunc("/", apiDefaultHandler)

	for _, pathMapping := range pathMappings {

		http.HandleFunc(pathMapping.Path, handlerFactory(pathMapping))

	}

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handlerFactory(pathMapping db.PathMapping) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		_, _ = io.WriteString(writer, fmt.Sprintf("You entered in the %s endpoint", pathMapping.Table))

	}
}
