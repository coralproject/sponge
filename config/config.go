/*
Package config handles the loading and distribution of configuration related with external sources.

{

	"Name": "New York Times",
	Strategy: {
		"type": "mysql", //To Do: look for a better Name
		"tables": {"comments": "nyt_comments"},
	},

	"Credentials": [{
		"database":  "",
		"username":  "",
		"password":  "",
		"host":      "",
		"port":     (int),
		"adapter":  "".
		"type": "source"
	},
	{
		"database":  "",
		"username":  "",
		"password":  "",
		"host":      "",
		"port":     (int),
		"adapter":  "".
		"type": "local"
	}
	]

}

*/
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/coralproject/core/log"
)

// Strategy explains which tables or data we are getting from the source.
type Strategy struct {
	Typesource string            `json:"typesource"`
	Tables     map[string]string `json:"tables"`
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

	config, err := readConfigFile("config/config.json")
	if err != nil {
		log.Fatal("Error when getting the configuration file. ", err)
	}

	return config
}

// GetCredential returns the credentials for connection with the external source
func (conf Config) GetCredential(adapter string) Credential {
	var cred Credential

	credentials := conf.Credentials

	// look at the credentials related to local database (mongodb in our original example)
	for i := 0; i < len(credentials); i++ {
		if credentials[i].Adapter == adapter {
			cred = credentials[i]
			return cred
		}
	}

	log.Fatal("Error when trying to get the credential for ", adapter)

	return cred
}

// GetStrategy returns the strategy
// To Do: Needs to manage errors
func (conf Config) GetStrategy() Strategy {
	strategy := conf.Strategy

	// To Do: catch the error on getting credentials and return it
	return strategy
}

// GetTables returns a list of tables to be imported
func (conf Config) GetTables() map[string]string {
	// To Do: catch the error when no Tables
	return conf.Strategy.Tables
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

func readConfigFile(f string) (*Config, error) {

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("Unable to read config file ", f, err)
		return nil, err
	}

	config, err := unmarshal(content)
	if err != nil {
		log.Fatal("Unable to parse JSON in config file ", f, err)
		return nil, err
	}

	return config, err
}
