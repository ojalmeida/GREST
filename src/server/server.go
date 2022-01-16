package server

import (
	"github.com/ojalmeida/GREST/src/db"
	"net/http"
)

var behaviors []db.Behavior

var server http.Server
var serverMux *http.ServeMux

var configServerMux *http.ServeMux
var configServer http.Server
