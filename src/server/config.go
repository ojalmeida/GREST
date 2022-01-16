package server

import (
	"io/ioutil"
	"log"
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Project  string `yaml:"Project"`
	Version  string `yaml:"Version"`
	Database struct {
		DBMS     string `yaml:"DBMS"`
		Address  string `yaml:"Address"`
		Port     string `yaml:"Port"`
		Database string `yaml:"database"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	} `yaml:"Database"`
	Listener struct {
		Production struct {
			Address string `yaml:"Address"`
			Port    string `yaml:"Port"`
		} `yaml:"Production"`
		Management struct {
			Address string `yaml:"Address"`
			Port    string `yaml:"Port"`
		} `yaml:"Management"`
	} `yaml:"Listener"`
}

/*
	GConfig opens and turns the configuration file into a struct
	that will be used as base to start the connections.
	Steps the func takes
	1. Open file
		if failed: create file and insert default configuration.
	2. Get content
	3. create a "Config" type struct
	4. Unmarshal yaml to struct
	5. return Config structure
 */
func GConfig() Config {
	confFile := "config.yaml"
	log.Println("Opening configuration file")
	cf,err := os.Open(confFile)
	if err != nil {
		log.Println(confFile, "was not found!")
		log.Println("\t└──Trying to create", confFile)
		_, err = os.Create(confFile)
		if err != nil {
			log.Println("\t\t└──Fail!")
		} else {
			log.Println("\t\t├──Success")
			log.Println("\t\t└──Writing default configuration")
			file,_ := os.OpenFile(confFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			_,err = file.WriteString(defualtConf())
			if err != nil {
				log.Println(err.Error())
				log.Println("\t\t\t└──Fail!")
			} else {
				log.Println("\t\t\t└──Success!")
			}
			file.Close()
		}
	} else {
		log.Println(confFile, "found!")
	}
	cf.Close()

	body, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Println(err.Error())
	}

	C := Config{}
	err = yaml.Unmarshal(body, &C)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return C
}


func defualtConf() string {
	return `Project: ProjectName
Version: 1.0

Database:
  DBMS: mysql
  Address: 127.0.0.1
  Port: 3306
  database: grestdb
  Username: Filipe
  Password: Tetris123

Listener:
  Production:
    Address: 0.0.0.0
    Port: 8080
  Management:
    Address: 0.0.0.0
    Port: 9090
`
}