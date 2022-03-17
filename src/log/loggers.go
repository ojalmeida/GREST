package log

import (
	"github.com/ojalmeida/GREST/src/config"
	"log"
	"os"
)

var (
	WarningLogger          *log.Logger
	warningLogPath         *string
	InfoLogger             *log.Logger
	infoLogPath            *string
	ErrorLogger            *log.Logger
	errorLogPath           *string
	ServerLogger           *log.Logger
	serverLoggerPath       *string
	ConfigServerLogger     *log.Logger
	configServerLoggerPath *string
)

func Start() {

	warningLogPath = &config.Conf.Logging.Warning
	infoLogPath = &config.Conf.Logging.Info
	errorLogPath = &config.Conf.Logging.Error
	serverLoggerPath = &config.Conf.Logging.Server
	configServerLoggerPath = &config.Conf.Logging.ConfigServer

	infoLogFile, err := os.OpenFile(*infoLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	warnLogFile, err := os.OpenFile(*warningLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	errLogFile, err := os.OpenFile(*errorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	serverLogFile, err := os.OpenFile(*serverLoggerPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	configServerLogFile, err := os.OpenFile(*configServerLoggerPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(warnLogFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	ServerLogger = log.New(serverLogFile, "SERVER: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	ConfigServerLogger = log.New(configServerLogFile, "CONFIG_SERVER: ", log.Ldate|log.Ltime|log.Lmicroseconds)
}
