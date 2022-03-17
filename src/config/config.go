package config

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Project string  `yaml:"Project"`
	Version float64 `yaml:"Version"`
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
	Configuration struct {
		AutoReload      bool   `yaml:"AutoReload"`
		DeclarationMode string `yaml:"DeclarationMode"`
	} `yaml:"Configuration"`
	Logging struct {
		Info         string `yaml:"Info"`
		Warning      string `yaml:"Warning"`
		Error        string `yaml:"Error"`
		Server       string `yaml:"Server"`
		ConfigServer string `yaml:"ConfigServer"`
	} `yaml:"Logging"`
}

var Conf = Config{}
var MainFolder string

var REFRESH_TIME = time.Second

func init() {

	var home string
	var err error

	if home, err = os.UserHomeDir(); err != nil {

		log.Panicln("Impossible to get user home directory")

	}

	MainFolder = home + "/.grest"
}

func Start(reloadAppChannel *chan bool, wg *sync.WaitGroup) {

	LoadConfigs()
	go ListenToChanges(reloadAppChannel, wg)

}

func getDefaultConfig() string {
	return fmt.Sprintf(`
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
  Schema: grest
  Username: root
  Password: root

Configuration:
  AutoReload: false
  DeclarationMode: file

Logging:
  Info: %s/info.log
  Warning: %s/warn.log
  Error: %s/error.log
  Server: %s/server.log
  ConfigServer: %s/config-server.log

`, MainFolder, MainFolder, MainFolder, MainFolder, MainFolder)
}

func ListenToChanges(reloadAppChannel *chan bool, wg *sync.WaitGroup) {

	firstRun := true
	var oldConf = Conf
	var newConf Config

	defer func() {
		if err := recover(); err != nil {

		}
	}() // recover of closed channel reading's panic

	for {

		time.Sleep(REFRESH_TIME)

		configFilePath := MainFolder + "/config.yaml"

		data, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			log.Panicln(err.Error())
		}

		err = yaml.Unmarshal(data, &newConf)
		if err != nil {
			log.Println("Error unmarshalling configuration file: ", err)
		}

		if !firstRun {
			if !reflect.DeepEqual(oldConf, newConf) {

				*reloadAppChannel <- true

				wg.Done()

				return

			}
		}

		firstRun = false

	}

}

func LoadConfigs() {

	log.Println("Opening configuration file")

	if _, err := os.Stat(MainFolder); os.IsNotExist(err) {

		log.Println(MainFolder + " does not exists, trying to create")

		if err = os.Mkdir(MainFolder, 0770); err != nil {

			log.Panicln(err.Error())

		}

		log.Println("Success on creating GREST folder")

	}


	confFilePath := MainFolder + "/config.yaml"
	_, err := os.Open(confFilePath)
	if err != nil {

		log.Println(confFilePath, "was not found")
		log.Println("Trying to create", confFilePath)

		_, err = os.Create(confFilePath)
		if err != nil {
			log.Panicln("Fail on configuration file creation")
		} else {
			log.Println("Success on configuration file creation")
			log.Println("Writing default configuration")
			configFile, _ := os.OpenFile(confFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

			_, err = configFile.WriteString(getDefaultConfig())

			if err != nil {
				log.Panicln("Fail on writing default configuration", err.Error())
			} else {
				log.Println("Success on writing default configuration")
			}

			err := configFile.Close()

			if err != nil {
				log.Println("Error on file closing: ", err.Error())
			}

		}
	} else {

		log.Println(confFilePath, "found!")

	}

	data, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		log.Panicln(err.Error())
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Println("Error unmarshalling configuration file: ", err)
	}

}


func ChangeConfig(key, newValue string) error {

	confFilePath := MainFolder + "/config.yaml"

	configFile, err := os.Open(confFilePath)

	if err != nil {

		return err
	}

	reader := bufio.NewScanner(configFile)
	oldLineToBeChanged := ""
	var oldFile []string

	for reader.Scan() {

		line := reader.Text()
		oldFile = append(oldFile, line)

		re := regexp.MustCompile(fmt.Sprintf(`\s*%s\s*:\s*\S*\s*`, key))

		if re.MatchString(line) {

			oldLineToBeChanged = line

		}

	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	var newFile []string

	for _, line := range oldFile {

		if line == oldLineToBeChanged {

			spaces := regexp.MustCompile(`^\s*`).FindString(line)

			newFile = append(newFile, fmt.Sprintf("%s%s: %s", spaces, key, newValue))

		} else {

			newFile = append(newFile, line)
		}

	}

	err = os.WriteFile(confFilePath, []byte(strings.Join(newFile, "\n")), 0660)

	if err != nil {
		return err
	}

	return nil

}
