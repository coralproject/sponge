/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* PostgreSQL
* MongoDB
* Webservice

*/
package source

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
)

// Options will hold all the options that came from flags
type Options struct {
	Limit                 int
	Offset                int
	Orderby               string
	Query                 string
	Types                 string
	Importonlyfailed      bool
	ReportOnFailedRecords bool
	Reportdbfile          string
	TimeWaiting           int
}

// global variables related to strategy
var (
	strategy str.Strategy
	uuid     string
)

// Global configuration variables that holds the credentials for the foreign data source connection
var credential str.Credential

// Init initialize the needed variables
func Init(u string) (string, error) {

	var err error

	uuid = u
	str.Init(uuid)

	strategy, err = str.New()
	if err != nil {
		log.Error(uuid, "source.init", err, "Get Strategy Configuration")
		return "", err
	}

	credential, err = strategy.GetCredential(strategy.Map.Foreign, "foreign")
	if err != nil {
		log.Error(uuid, "source.init", err, "Get Credentials for external source")
		return "", err
	}

	return strategy.Map.Foreign, err
}

// Sourcer is where the data is coming from (webservice or database)
type Sourcer interface {
	// GetData returns data for a specific entity filter by all the options
	GetData(string, *Options) ([]map[string]interface{}, error)                //int, int, string, string // args ...interface{}) ([]map[string]interface{}, error) //
	GetQueryData(string, *Options, []string) ([]map[string]interface{}, error) //int, int, string
	IsWebService() bool
}

// New returns a new Source struct with the connection string in it
func New(d string) (Sourcer, error) {

	switch d {
	case "mysql":
		// Get MySQL connection string
		return MySQL{Connection: connectionMySQL(), Database: nil}, nil
	case "mongodb":
		// Get MongoDB connection string
		return MongoDB{Connection: connectionMongoDB(), Database: nil}, nil
	case "postgresql":
		// Get MySQL connection string
		return PostgreSQL{Connection: connectionPostgreSQL(), Database: nil}, nil
	case "service":
		// Get API connection url
		u := connectionAPI("")
		return API{Connection: u.String()}, nil
	}

	return nil, fmt.Errorf("Configuration not found for source database %s.", d)
}

// GetEntities gets all the entities names from this data source
func GetEntities() ([]string, error) {

	collections := strategy.GetEntities()

	keys := make([]string, len(collections))

	for k, val := range collections {
		keys[val.Priority] = k
	}

	return keys, nil
}

// GetForeignEntity returns the foreign's source of the entity
func GetForeignEntity(name string) string {
	return strategy.GetEntityForeignName(name)
}

//**** Utility functions used for the API and Mongo GetData func ****//

// it prepares the data to have the transformations in fiddler
// flattenizeData converts data into a map[string]string with the key a breadcrumb to the leaf, and the value being the leaf itself
func flattenizeData(originalData []map[string]interface{}) ([]map[string]interface{}, error) {
	var dat []map[string]interface{}

	for _, j := range originalData { // this is a slice of maps
		d, e := flattenizeDocument(j)
		if e != nil {
			return nil, e
		}
		// add d to dat
		dat = append(dat, d)
	}

	return dat, nil
}

//
func flattenizeDocument(document map[string]interface{}) (map[string]interface{}, error) {
	d := make(map[string]interface{})

	for i, k := range document { // this is the actual map
		m := flattenize(i, k) // gets the id being the breadcrumb and val the leaf
		for r, p := range m {
			d[r] = p
		}
	}

	return d, nil
}

// flattenize is a recursive function
// that ends up returning a map which key is the breadcrumb and
// the value is the leaf
// i is the first key. k could be a string or a map
func flattenize(i string, k interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	switch v := k.(type) {
	case map[string]interface{}: // if k is a map then go deeper
		for p, e := range k.(map[string]interface{}) {
			switch v := e.(type) {
			case map[string]interface{}:
				newi := strings.Join([]string{i, p}, ".")
				newresult := flattenize(newi, e)
				for v1, v2 := range newresult {
					result[strings.ToLower(v1)] = v2
				}
			case []map[string]interface{}:
				for u, c := range v {
					newi := strings.Join([]string{i, p, strconv.Itoa(u)}, ".")
					newresult := flattenize(newi, c)
					for v1, v2 := range newresult {
						result[strings.ToLower(v1)] = v2
					}
				}
			case map[string]string:
				//fmt.Printf("* %v is map[string]string\n\n", e)
				for a, b := range v {
					newi := strings.Join([]string{i, p, a}, ".")
					result[strings.ToLower(newi)] = b
				}
			case []map[string]string:
				for u, c := range v {
					newi := strings.Join([]string{i, p, strconv.Itoa(u)}, ".")
					newresult := flattenize(newi, c)
					for v1, v2 := range newresult {
						result[strings.ToLower(v1)] = v2
					}
				}
			case []interface{}:
				// Ugly work-around to check if this is a []string and we need to send it as it is
				if len(v) != 0 {
					switch v[0].(type) {
					case string:
						newi := strings.Join([]string{i, p}, ".")
						result[strings.ToLower(newi)] = e
						continue
					}
				}
				for d1, d2 := range v {
					newi := strings.Join([]string{i, p, strconv.Itoa(d1)}, ".")
					newresult := flattenize(newi, d2)
					for v1, v2 := range newresult {
						result[strings.ToLower(v1)] = v2
					}
				}
			default:
				newi := strings.Join([]string{i, p}, ".")
				result[strings.ToLower(newi)] = e
			}
		}
	case map[string]string:
		for p, e := range v {
			newi := strings.Join([]string{i, p}, ".")
			result[strings.ToLower(newi)] = e
		}
	case []map[string]string:
		for u, c := range v {
			newi := strings.Join([]string{i, strconv.Itoa(u)}, ".")
			newresult := flattenize(newi, c)
			for v1, v2 := range newresult {
				result[strings.ToLower(v1)] = v2
			}
		}
	case []map[string]interface{}:
		for u, c := range v {
			newi := strings.Join([]string{i, strconv.Itoa(u)}, ".")
			newresult := flattenize(newi, c)
			for v1, v2 := range newresult {
				result[strings.ToLower(v1)] = v2
			}
		}
	case []interface{}:
		for u, c := range v {
			newi := strings.Join([]string{i, strconv.Itoa(u)}, ".")
			newresult := flattenize(newi, c)
			for v1, v2 := range newresult {
				result[strings.ToLower(v1)] = v2
			}
		}
	default: // if k is not a map then just return it as a string
		result[strings.ToLower(i)] = k
	}

	return result
}
