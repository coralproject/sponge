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
	"fmt"
	"io/ioutil"

	"github.com/coralproject/core/log"
)

//* Errors used in this package *//

// When trying to connect to the database with the connection string
type endpointError struct {
	key string
}

func (e endpointError) Error() string {
	return fmt.Sprintf("Error when trying to get endpoint %s.", e.key)
}

// Table holds the struct on what is the external source's table name and fields
type Table struct {
	Name   string            `json:"name"`
	Fields map[string]string `json:"fields"`
}

// Strategy explains which tables or data we are getting from the source.
type Strategy struct {
	Typesource string           `json:"typesource"`
	Tables     map[string]Table `json:"tables"`
	Action     []Table
}

// Credential is the interface for APIs or Database Sources
type Credential interface {
	GetAdapter() string
}

// CredentialDatabase has the information to connect to the external databaase source.
type CredentialDatabase struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`    //= '5432'
	Adapter  string `json:"adapter"` //= 'mysql'
	Type     string `json:"type"`    //= 'source' or 'local'
}

// GetAdapter returns the adapter
func (c CredentialDatabase) GetAdapter() string {
	return c.Adapter
}

// CredentialAPI has the information to connect to an external API source.
type CredentialAPI struct {
	Username  string            `json:"username"` // BasicAuth
	Password  string            `json:"password"` // BasicAuth
	Adapter   string            `json:"adapter"`
	Endpoints map[string]string `json:"endpoints"`
}

// GetAdapter returns the adapter
func (c CredentialAPI) GetAdapter() string {
	return c.Adapter
}

// GetEndpoints returns all the endpoints
func (c CredentialAPI) GetEndpoints() map[string]string {
	return c.Endpoints
}

// GetEndpoint gives the endpoint for that modelName
func (c CredentialAPI) GetEndpoint(modelName string) (string, error) {
	endpoints := c.GetEndpoints()
	for k, e := range endpoints {
		if k == modelName {
			return e, nil
		}
	}
	err := fmt.Errorf("Endpoint %s not found.", modelName)

	return "", err
}

// GetAuthenticationEndpoint returns the authentication url
func (c CredentialAPI) GetAuthenticationEndpoint() (string, error) {
	for key, val := range c.Endpoints {
		if key == "authentication" {
			return val, nil
		}
	}

	return "", endpointError{key: "authentication"}
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
		if credentials[i].GetAdapter() == adapter {
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
func (conf Config) GetTables() map[string]Table {
	// To Do: catch the error when no Tables
	return conf.Strategy.Tables
}

// GetTableName returns the external source's table mapped to the coral model
func (conf Config) GetTableName(modelName string) string {
	return conf.Strategy.Tables[modelName].Name
}

// GetTableFields returns the external source's table fields mapped to the coral model
func (conf Config) GetTableFields(modelName string) map[string]string {
	return conf.Strategy.Tables[modelName].Fields
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
