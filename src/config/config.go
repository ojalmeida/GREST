package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Project string `yaml:"Project"`
	Version string `yaml:"Version"`
	API     struct {
		Production struct {
			Address string `yaml:"Address"`
			Port    string `yaml:"Port"`
		} `yaml:"Production"`
		Management struct {
			Address string `yaml:"Address"`
			Port    string `yaml:"Port"`
		} `yaml:"Management"`
	} `yaml:"API"`
	Database struct {
		DBMS     string `yaml:"DBMS"`
		Address  string `yaml:"Address"`
		Port     string `yaml:"Port"`
		Schema   string `yaml:"Schema"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	} `yaml:"Database"`
}


var Conf = Config{}


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
func init() {
	confFile := "config.yaml"
	log.Println("Opening configuration file")
	cf, err := os.Open(confFile)
	if err != nil {
		log.Println(confFile, "was not found!")
		log.Println("\t└──Trying to create", confFile)
		_, err = os.Create(confFile)
		if err != nil {
			log.Println("\t\t└──Fail!")
		} else {
			log.Println("\t\t├──Success")
			log.Println("\t\t└──Writing default configuration")
			file, _ := os.OpenFile(confFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			_, err = file.WriteString(getDefaultConfig())
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

	err = yaml.Unmarshal(body, &Conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func getDefaultConfig() string {
	return `
Project: ProjectName
Version: 1.0

API:
  Production:
    Address: 0.0.0.0
    Port: 8080
  Management:
    Address: 0.0.0.0
    Port: 9090

Database:
  DBMS: mysql
  Address: 127.0.0.1
  Port: 3306
  Schema: grestdb
  Username: root
  Password: root
`
}
