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
	pillarURL = os.Getenv("PILLAR_URL")
}

//* Strategy Structure *//

// Strategy is a structure with all the information for the specitic strategy (how to get the data, from which source)
type Strategy struct {
	Name        string
	Map         Map
	Credentials Credentials // map[string][]Credential // String is "Databases" or "APIs" indicating which kind of credentials are those
}

// Map explains which entities or data we are getting from the source.
type Map struct {
	Foreign        string            `json:"foreign"`
	DateTimeFormat string            `json:"datetimeformat"`
	Entities       map[string]Entity `json:"entities"`
}

// Entity holds the struct on what is the external source's entity name and fields
type Entity struct {
	Foreign        string                   `json:"foreign"`
	Local          string                   `json:"local"`
	Priority       int                      `json:"priority"`
	OrderBy        string                   `json:"orderby"`
	ID             string                   `json:"id"`
	Index          []mgo.Index              `json:"index"`  //map[string]interface{} `json:"index"`
	Fields         []map[string]interface{} `json:"fields"` // foreign (name in the foreign source), local (name in the local source), relation (relationship between each other), type (data type)
	Status         map[string]string        `json:"status"`
	PillarEndpoint string                   `json:"endpoint"`
}

///////////////////////////////////////////////////////////////////////////////

//** CREDENTIALS TO EXTERNAL SOURCES **//

// Credentials are all the credentials for external and internal data sources
type Credentials struct {
	Databases []CredentialDatabase `json:"databases"`
	APIs      []CredentialAPI      `json:"apis"`
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
	AppKey     string `json:"appkey"`
	Endpoint   string `json:"endpoint"`
	Adapter    string `json:"adapter"`
	Type       string `json:"type"`
	Records    string `json:"records"`
	Pagination string `json:"pagination"`
	UserAgent  string `json:"useragent"`
	Attributes string `json:"attributes"`
}

// GetAppKey gets the app key to access the api
func (c CredentialAPI) GetAppKey() string {
	return c.AppKey
}

// GetAdapter returns the adapter
func (c CredentialAPI) GetAdapter() string {
	return c.Adapter
}

// GetType returns the adapter
func (c CredentialAPI) GetType() string {
	return c.Type
}

// GetEndpoint returns all the endpoints
func (c CredentialAPI) GetEndpoint() string {
	return c.Endpoint
}

// GetRecordsFieldName returns the name of the field that holds the records
func (c CredentialAPI) GetRecordsFieldName() string {
	return c.Records
}

// GetPaginationFieldName returns the name of the field where to look for pagination
func (c CredentialAPI) GetPaginationFieldName() string {
	return c.Pagination
}

// GetUserAgent returns the name of the field that holds the user agent
func (c CredentialAPI) GetUserAgent() string {
	return c.UserAgent
}

// GetAttributes returns the attributes for the query
func (c CredentialAPI) GetAttributes() string {
	return c.Attributes
}

///////////////////////////////////////////////////////////////////////////////

/* Exported Functions */

// New creates a new strategy struct variable from the json file
func New() (Strategy, error) {

	var strategy Strategy
	var err error

	//read STRATEGY_CONF env variable
	strategyFile := os.Getenv("STRATEGY_CONF")
	if strategyFile == "" {
		log.Fatal(uuid, "strategy.new", "Enviromental variable STRATEGY_CONF not setup.")
		return strategy, err
	}

	strategy, err = read(strategyFile)
	if err != nil {
		log.Error(uuid, "strategy.new", err, "Reading strategy file %s.", strategyFile)
	}

	return strategy, err
}

// GetCredential returns the credentials for connection with the external source adapter a, type t
func (s Strategy) GetCredential(a string, t string) (Credential, error) {
	var cred Credential
	var err error

	creda := s.Credentials

	// look at the credentials related to local database (mongodb in our original example)
	for i := 0; i < len(creda.Databases); i++ {
		if creda.Databases[i].GetAdapter() == a && creda.Databases[i].GetType() == t {
			cred = creda.Databases[i]
			return cred, err
		}
	}

	// look at the credentials related to local database (mongodb in our original example)
	for i := 0; i < len(creda.APIs); i++ {
		if creda.APIs[i].GetAdapter() == a && creda.APIs[i].GetType() == t {
			cred = creda.APIs[i]
			return cred, err
		}
	}

	err = fmt.Errorf("Credential %s not found.", a)
	log.Error(uuid, "strategy.getCredentials", err, "Getting credential %s for strategy.", a)

	return cred, err
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
func (s Strategy) GetDateTimeFormat(entity string, field string) string {

	for _, f := range s.Map.Entities[entity].Fields {
		if f["local"] == field {
			val, exists := f["datetimeformat"]
			if exists {
				return val.(string)
			}
		}
	}
	return s.GetDefaultDateTimeFormat()
}

// GetEntities returns a list of the entities defined in the transformations file
func (s Strategy) GetEntities() map[string]Entity {
	return s.Map.Entities
}

// HasArrayField returns true if the entity has fields that are type array and need to be loop through
func (s Strategy) HasArrayField(t Entity) bool {
	//Fields   []map[string]interface{}

	for _, f := range t.Fields {
		if f["type"] == "Array" {
			return true
		}
	}
	return false
}

// GetFieldsForSubDocument get all the fields in the case of a subdocumetn
func (s Strategy) GetFieldsForSubDocument(model string, foreignfield string) []map[string]interface{} {
	var fields []map[string]interface{}

	for _, f := range s.Map.Entities[model].Fields { // search foreign field in []map[string]interface{}
		if f["foreign"] == foreignfield {
			fi := f["fields"].([]interface{})
			// Convert the []interface into []map[string]interface{}
			fields = make([]map[string]interface{}, len(fi))
			for i := range fields {
				fields[i] = fi[i].(map[string]interface{})
			}
			return fields
		}
	}
	return fields
}

// GetEntityForeignName returns the external source's entity mapped to the coral model
func (s Strategy) GetEntityForeignName(coralName string) string {
	return s.Map.Entities[coralName].Foreign
}

// GetEntityForeignFields returns the external source's entity fields mapped to the coral model
func (s Strategy) GetEntityForeignFields(coralName string) []map[string]interface{} {
	return s.Map.Entities[coralName].Fields
}

// GetOrderBy returns the order by field definied in the transformations file
func (s Strategy) GetOrderBy(coralName string) string {
	return s.Map.Entities[coralName].OrderBy
}

// GetIndexBy returns the structure to use to create indexes for the coral entity
func (s Strategy) GetIndexBy(coralName string) []mgo.Index {
	return s.Map.Entities[coralName].Index
}

// GetIDField returns the identifier for the entity coralname setup in the transformations file
func (s Strategy) GetIDField(coralName string) string {
	return s.Map.Entities[coralName].ID
}

// GetStatus returns the mapping of the external status into the coral one
func (s Strategy) GetStatus(coralName string, foreign string) string {
	return s.Map.Entities[coralName].Status[foreign]
}

// GetPillarEndpoints return the endpoints configured in the strategy
func (s Strategy) GetPillarEndpoints() map[string]string {
	endpoints := map[string]string{}

	entities := s.GetEntities()
	for _, entity := range entities {
		endpoints[entity.Local] = pillarURL + entity.PillarEndpoint
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
