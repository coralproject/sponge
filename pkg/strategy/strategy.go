/*
Package strategy handles the loading and distribution of configuration related with external sources.
*/
package strategy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/ardanlabs/kit/log"
)

var (
	pillarURL string
	uuid      string
)

// Init initialize log and get pillar url env variable
func Init(u string) {

	uuid = u

	// logLevel := func() int {
	// 	ll, err := cfg.Int("LOGGING_LEVEL")
	// 	if err != nil {
	// 		return log.USER
	// 	}
	// 	return ll
	// }
	//
	// log.Init(os.Stderr, logLevel)

	pillarURL = os.Getenv("PILLAR_URL")
}

//* Strategy Structure *//

// Strategy is a structure with all the information for the specitic strategy (how to get the data, from which source)
type Strategy struct {
	Name        string
	Map         Map
	Credentials Credentials // map[string][]Credential // String is "Databases" or "APIs" indicating which kind of credentials are those
}

// Map explains which tables or data we are getting from the source.
type Map struct {
	Foreign        string           `json:"foreign"`
	DateTimeFormat string           `json:"datetimeformat"`
	Tables         map[string]Table `json:"tables"`
}

// Table holds the struct on what is the external source's table name and fields
type Table struct {
	Foreign  string              `json:"foreign"`
	Local    string              `json:"local"`
	Priority int                 `json:"priority"`
	OrderBy  string              `json:"orderby"`
	ID       string              `json:"id"`
	Index    []mgo.Index         `json:"index"`  //map[string]interface{} `json:"index"`
	Fields   []map[string]string `json:"fields"` // foreign (name in the foreign source), local (name in the local source), relation (relationship between each other), type (data type)
	Endpoint string              `json:"endpoint"`
}

///////////////////////////////////////////////////////////////////////////////

//** CREDENTIALS TO EXTERNAL SOURCES **//

// Credentials are all the credentials for external and internal data sources
type Credentials struct {
	Databases []CredentialDatabase
	APIs      []CredentialAPI
}

// Credential is the interface that CredentialDatabase and CredentialAPI will implement
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

///////////////////////////////////////////////////////////////////////////////

/* Exported Functions */

// New creates a new strategy struct variable from the json file
func New() Strategy {

	var strategy Strategy
	var err error

	//read STRATEGY_CONF env variable
	strategyFile := os.Getenv("STRATEGY_CONF")
	if strategyFile == "" {
		log.Fatal(uuid, "strategy.new", "Enviromental variable STRATEGY_CONF not setup.")
		return strategy
	}

	strategy, err = read(strategyFile)
	if err != nil {
		log.Error(uuid, "strategy.new", err, "Reading strategy file %s.", strategyFile)
	}

	return strategy
}

// Validate checks that the strategy file is correct
// the tables are part of the coral system
// it has at least one external source credential
func (s Strategy) Validate() error {
	return nil
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

	log.Error(uuid, "strategy.getCredentials", fmt.Errorf("Credential %s not found.", a), "Getting credential %s for strategy.", a)

	return cred
}

// GetMap returns the strategy
func (s Strategy) GetMap() Map {
	return s.Map
}

// GetDefaultDateTimeFormat gets the default datetime format
func (s Strategy) GetDefaultDateTimeFormat() string {
	return s.Map.DateTimeFormat
}

// GetDateTimeFormat returns the datetime format for this strategy
func (s Strategy) GetDateTimeFormat(table string, field string) string {

	for _, f := range s.Map.Tables[table].Fields {
		if f["local"] == field {
			val, exists := f["datetimeformat"]
			if exists {
				return val
			}
		}
	}
	return s.GetDefaultDateTimeFormat()
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

// GetOrderBy returns the order by field definied in the strategy
func (s Strategy) GetOrderBy(coralName string) string {
	return s.Map.Tables[coralName].OrderBy
}

// GetIndexBy returns the structure to use to create indexes for the coral table
func (s Strategy) GetIndexBy(coralName string) []mgo.Index { //map[string]interface{} {
	return s.Map.Tables[coralName].Index
}

// GetIDField returns the identifier for the table coralname setup in the strategy file
func (s Strategy) GetIDField(coralName string) string {
	return s.Map.Tables[coralName].ID
}

// GetPillarEndpoints return the endpoints configured in the strategy
func (s Strategy) GetPillarEndpoints() map[string]string {
	endpoints := map[string]string{}

	tables := s.GetTables()
	for _, table := range tables {
		endpoints[table.Local] = pillarURL + table.Endpoint
	}

	// adds CREATE_INDEX endpoints
	endpoints["index"] = pillarURL + "/api/import/index"

	return endpoints
}

/* Not Exported Functions */

// Read the strategy file and do the validation into the Strategy struct
func read(f string) (Strategy, error) {

	var strategy Strategy

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Error(uuid, "strategy.read", err, "Reading strategy file %s.", f)
		return strategy, err
	}

	err = json.Unmarshal(content, &strategy)
	if err != nil {
		log.Error(uuid, "strategy.read", err, "Unmarshal strategy %s.", f)
	}

	return strategy, err
}
