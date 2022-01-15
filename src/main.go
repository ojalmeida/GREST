package main

import (
	"fmt"
	"github.com/ojalmeida/GREST/src/server"
	"os"
	"os/signal"
)

func main() {

	fmt.Println("PID:", os.Getpid())

	go server.StartServers()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

}
