package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/coralproject/core/log"
)

// Accesible has the methods for Config struct
type Accesible interface {
	New() *Accesible
	GetCredentials() (Credential, error)
	GetStrategy() (Strategy, error)
}

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

// Credential has the information to connect to the external source.
type Credential struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`    //= '5432'
	Adapter  string `json:"adapter"` //= 'mysql'
	Type     string `json:"type"`    //= 'source' or 'local'
}

// Config is a structure with all the information for the specitic strategy (how to get the data, from which source)
type Config struct {
	Name        string
	Strategy    Strategy
	Credentials []Credential
}

// Pointer to the master config record
//var config *Config

/* Exported Functions */

// New creates a new config
func New() *Config {
	config, err := readFile("config/config.json")
	if err != nil {
		log.Fatal("It couldn't read the configuration file. ", err)
	}

	return config
}

// GetCredentials returns the credentials for connection with the external source
func (conf Config) GetCredentials(typec string) Credential {
	var cred Credential

	credentials := conf.Credentials

	// look at the credentials related to local database (mongodb in our original example)
	for i := 0; i < len(credentials); i++ {
		if credentials[i].Type == typec {
			cred = credentials[i]
			return cred
		}
	}

	log.Fatal("Error when trying to get the connection string for mongodb.")

	return cred
}

// GetStrategy returns the strategy
// To Do: Needs to manage errors
func (conf Config) GetStrategy() (Strategy, error) {
	strategy := conf.Strategy

	// To Do: catch the error on getting credentials and return it
	return strategy, nil
}

/* Not Exported Functions */

// Read the configuration file and load it into the Config
func unmarshal(content []byte) (*Config, error) {

	c := new(Config)

	err := json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func readFile(f string) (*Config, error) {

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("Unable to read config file ", f, err)
		return nil, err
	}

	c, err := unmarshal(content)
	if err != nil {
		log.Fatal("Unable to parse JSON in config file ", f, err)
		return nil, err
	}

	return c, err
}
