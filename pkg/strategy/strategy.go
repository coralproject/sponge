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
	Database CredentialDatabase `json:"database"`
	Service  CredentialService  `json:"service"`
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

// CredentialService has the information to connect to an external web service source.
type CredentialService struct {
	AppKey          string `json:"appkey"`
	Endpoint        string `json:"endpoint"`
	Adapter         string `json:"adapter"`
	Type            string `json:"type"`
	Records         string `json:"records"`
	PaginationField string `json:"paginationfield"`
	Pagination      string `json:"pagination"`
	UserAgent       string `json:"useragent"`
	Attributes      string `json:"attributes"`
}

// GetAppKey gets the app key to access the api
func (c CredentialService) GetAppKey() string {
	return c.AppKey
}

// GetPageAfterField returns the field that I need to send to the API to get the next page
func (c CredentialService) GetPageAfterField() string {
	return c.PaginationField
}

// GetAdapter returns the adapter
func (c CredentialService) GetAdapter() string {
	return c.Adapter
}

// GetType returns the adapter
func (c CredentialService) GetType() string {
	return c.Type
}

// GetEndpoint returns all the endpoints
func (c CredentialService) GetEndpoint() string {
	return c.Endpoint
}

// GetRecordsFieldName returns the name of the field that holds the records
func (c CredentialService) GetRecordsFieldName() string {
	return c.Records
}

// GetPaginationFieldName returns the name of the field where to look for pagination
func (c CredentialService) GetPaginationFieldName() string {
	return c.Pagination
}

// GetUserAgent returns the name of the field that holds the user agent
func (c CredentialService) GetUserAgent() string {
	return c.UserAgent
}

// GetAttributes returns the attributes for the query
func (c CredentialService) GetAttributes() string {
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
		return Strategy{}, err
	}

	err = strategy.setCredential()
	if err != nil {
		log.Error(uuid, "strategy.new", err, "Setting credentials.")
		return Strategy{}, err
	}

	return strategy, err
}

func (s Strategy) setCredential() error {
	var err error
	// DATABASE
	DBdatabase := os.Getenv("DB_DATABASE")
	if DBdatabase != "" {
		DBusername := os.Getenv("DB_USERNAME")
		DBpassword := os.Getenv("DB_PASSWORD")
		DBhost := os.Getenv("DB_HOST")
		DBport := os.Getenv("DB_PORT")

		s.Credentials.Database = CredentialDatabase{
			Database: DBdatabase,
			Username: DBusername,
			Password: DBpassword,
			Host:     DBhost,
			Port:     DBport,
		}
	}

	// WEB SERVICE
	WSappkey := os.Getenv("WS_APPKEY")
	if WSappkey != "" {

		WSendpoint := os.Getenv("WS_ENDPOINT")
		WSrecords := os.Getenv("WS_RECORDS")
		WSpagination := os.Getenv("WS_PAGINATION")
		WSuseragent := os.Getenv("WS_USERAGENT")
		WSattributes := os.Getenv("WS_ATTRIBUTES")

		s.Credentials.Service = CredentialService{
			AppKey:     WSappkey,
			Endpoint:   WSendpoint,
			Records:    WSrecords,
			Pagination: WSpagination,
			UserAgent:  WSuseragent,
			Attributes: WSattributes,
		}
	}
	return err
}

// GetCredential returns the credential for connection with the external source adapter a, type t
func (s Strategy) GetCredential(a string, t string) (Credential, error) {
	var cred Credential
	var err error

	creda := s.Credentials

	if creda.Database.GetAdapter() == a && creda.Database.GetType() == t {
		return creda.Database, err
	}

	if creda.Service.GetAdapter() == a && creda.Service.GetType() == t {
		return creda.Service, err
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
	for name, entity := range entities {
		endpoints[name] = pillarURL + entity.PillarEndpoint //entity.Local
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
