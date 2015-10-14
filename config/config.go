package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/coralproject/core/log"
	"github.com/fatih/structs"
)

/*
Package config handles the loading and distribution of configuration related with external sources.

{

	"Name": "New York Times",
	Strategy: {
		"type": "mysql", //To Do: look for a better Name
		"tables": ["comments"],
		"anonymize": false
	},

	"Credentials": {
		"database":  "",
		"username":  "",
		"password":  "",
		"host":      "",
		"port":     (int),
		"adapter":  ""
	}

}

*/

// Strategy explains which tables or data we are getting from the source.
type Strategy struct {
	tables []string
}

// Credentials has the information to connect to the external source.
type Credentials struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`    //= '5432'
	Adapter  string `json:"adapter"` //= 'mysql'
}

// Config is a structure with all the information for the specitic strategy (how to get the data, from which source)
type Config struct {
	Name        string
	Strategy    Strategy
	Credentials Credentials
}

// Pointer to the master config record
var config *Config

// Read the configuration file and load it into the Config
func unmarshal(content []byte) (*Config, error) {

	c := new(Config)

	err := json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func readFile(f string) *Config {

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("Unable to read config file ", f, err)
	}

	c, err := unmarshal(content)
	if err != nil {
		log.Fatal("Unable to parse JSON in config file ", f, err)
	}

	return c
}

// Get returns the config
func Get() *Config {
	return config
}

// GetCredentials returns the credentials for connection with the external source
func GetCredentials() (Credentials, error) {
	dict := structs.Map(config)
	_, ok := dict["Credentials"]
	if ok {
		err := errors.New("No Credentials option in the Configuration file.")
		return Credentials{}, err
	}
	return config.Credentials, nil
}

// GetStrategy returns the strategy
// To Do: Needs to manage errors
func GetStrategy() (Strategy, error) {
	return config.Strategy, nil
}
