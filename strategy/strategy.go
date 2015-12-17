/*
Package strategy handles the loading and distribution of configuration related with external sources.
*/
package strategy

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ardanlabs/kit/log"
)

// func init() {
//  // validate CONFIGURATION
// }

//* Strategy Structure *//

// Strategy is a structure with all the information for the specitic strategy (how to get the data, from which source)
type Strategy struct {
	Name        string
	Map         Map
	Credentials Credentials // map[string][]Credential // String is "Databases" or "APIs" indicating which kind of credentials are those
}

// Map explains which tables or data we are getting from the source.
type Map struct {
	Foreign string           `json:"foreign"`
	Tables  map[string]Table `json:"tables"`
}

// Table holds the struct on what is the external source's table name and fields
type Table struct {
	Foreign string              `json:"foreign"`
	Local   string              `json:"local"`
	Fields  []map[string]string `json:"fields"` // foreign (name in the foreign source), local (name in the local source), relation (relationship between each other), type (data type)
}

// ^
//Fields has maps in the style
// {
// 	"foreign": "parentid",
// 	"local": "ParentID",
// 	"relation": "Identity",
// 	"type": "int"
// }

// Credentials are all the credentials for external and internal data sources
type Credentials struct {
	Databases []CredentialDatabase
	APIs      []CredentialAPI
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

// GetType returns the adapter
func (c CredentialDatabase) GetType() string {
	return c.Type
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

	return "", endpointError{key: modelName}
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

/* Exported Functions */

// New creates a new strategy from configuration file
func New() Strategy {

	var strategy Strategy
	var err error

	f := "strategy.json"

	strategy, err = readConfigFile(f)

	if err != nil {
		log.Error("setting", "new", err, "Getting strategy file")
	}

	// err = validate(strategy)
	// if err != nil {
	// 	log.Error("setting", "new", err, "Validating strategy file")
	// }

	return strategy
}

// GetCredential returns the credentials for connection with the external source adapter a, type t
func (s Strategy) GetCredential(a string, t string) CredentialDatabase {
	var cred CredentialDatabase

	creda := s.Credentials.Databases

	// look at the credentials related to local database (mongodb in our original example)
	for i := 0; i < len(creda); i++ {
		if creda[i].GetAdapter() == a && creda[i].GetType() == t {
			cred = creda[i]
			return cred
		}
	}

	//	log.Error("startup", "getCredentials", errors.New("Credential not found."), "Getting strategy")

	return cred
}

// GetMap returns the strategy
func (s Strategy) GetMap() Map {
	return s.Map
}

// GetTables returns a list of tables to be imported
func (s Strategy) GetTables() map[string]Table {
	// To Do: catch the error when no Tables
	return s.Map.Tables
}

// GetTableForeignName returns the external source's table mapped to the coral model
func (s Strategy) GetTableForeignName(coralName string) string {
	return s.Map.Tables[coralName].Foreign
}

// GetTableForeignFields returns the external source's table fields mapped to the coral model
func (s Strategy) GetTableForeignFields(coralName string) []map[string]string {
	return s.Map.Tables[coralName].Fields
}

/* Not Exported Functions */

// Read the configuration file and load it into the Config
func unmarshal(content []byte) (Strategy, error) {

	s := Strategy{}

	err := json.Unmarshal(content, &s)
	if err != nil {
		log.Error("startup", "unmarshal", err, "Getting strategy")
		return Strategy{}, err
	}

	return s, nil
}

func readConfigFile(f string) (Strategy, error) {

	//log.Dev("startup", "readConfigFile", "Reading Config File : %s", f)

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Error("startup", "readConfigFile", err, "Getting strategy")

		return Strategy{}, err
	}

	strategy, err := unmarshal(content)
	if err != nil {
		log.Error("startup", "readConfigFile", err, "Getting strategy")
		return Strategy{}, err
	}

	return strategy, err
}
