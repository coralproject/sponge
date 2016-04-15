/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
	"github.com/stretchr/stew/objects"
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
	case "api":
		// Get API connection url
		u := connectionAPI()
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

//**** Utility functions used for the API and Mongo GetData func ****//

// it prepares the data to have the transformations in fiddler
// normalize converts into a map[string]string with the key a breadcrumb to the leaf, and the value being the leaf itself
func flattenData(fields []string, mongoData []map[string]interface{}) ([]map[string]interface{}, error) {
	var dat []map[string]interface{}

	for _, j := range mongoData { // this is a slice of maps
		d, e := flattDocument(fields, j)
		if e != nil {
			return nil, e
		}
		// add d to dat
		dat = append(dat, d)
	}

	return dat, nil
}

func flattDocument(fields []string, document map[string]interface{}) (map[string]interface{}, error) {

	var err error
	result := make(map[string]interface{})

	for _, field := range fields {
		result[field] = objects.Map(document).Get(field)
	}

	return result, err
}

// it prepares the data to have the transformations in fiddler
// normalize converts into a map[string]string with the key a breadcrumb to the leaf, and the value being the leaf itself
func normalizeData(mongoData []map[string]interface{}) ([]map[string]interface{}, error) {
	var dat []map[string]interface{}

	for _, j := range mongoData { // this is a slice of maps
		d, e := normalizeDocument(j)
		if e != nil {
			return nil, e
		}
		// add d to dat
		dat = append(dat, d)
	}

	return dat, nil
}

func normalizeDocument(document map[string]interface{}) (map[string]interface{}, error) {
	d := make(map[string]interface{})

	for i, k := range document { // this is the actual map
		//fmt.Printf("## Index %s, normalizing %v \n\n", i, k)
		m := normalize(i, k) // gets the id being the breadcrumb and val the leaf
		//fmt.Printf("- Got %v\n\n\n", m)
		for r, p := range m {
			d[r] = p
		}
	}

	return d, nil
}

// i is the first key. k could be a string or a map
// it returns key and value
func normalize(i string, k interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	switch v := k.(type) {
	case map[string]interface{}: // if k is a map then go deeper
		//fmt.Printf("** %v is map[string]interface{}\n\n", k)
		for p, e := range k.(map[string]interface{}) {
			switch v := e.(type) {
			case map[string]interface{}:
				//fmt.Printf("* %v is map[string]interface{}\n\n", e)
				newi := strings.Join([]string{i, p}, ".")
				newresult := normalize(newi, e)
				// copy newresult into result
				for v1, v2 := range newresult {
					result[v1] = v2
				}
			case []map[string]interface{}:
				//fmt.Printf("* %v is []map[string]interface{}\n\n", e)
				for u, c := range v {
					newi := strings.Join([]string{i, p, strconv.Itoa(u)}, ".")
					newresult := normalize(newi, c)
					for v1, v2 := range newresult {
						result[v1] = v2
					}
				}
			case map[string]string:
				//fmt.Printf("* %v is map[string]string\n\n", e)
				for a, b := range v {
					newi := strings.Join([]string{i, p, a}, ".")
					result[newi] = b
				}
			case []map[string]string:
				//fmt.Printf("* %v is []map[string]string\n\n", e)
				for u, c := range v {
					newi := strings.Join([]string{i, p, strconv.Itoa(u)}, ".")
					newresult := normalize(newi, c)
					for v1, v2 := range newresult {
						result[v1] = v2
					}
				}
			case []interface{}:
				//fmt.Printf("* %v is []interface{}\n\n", e)
				for d1, d2 := range v {
					newi := strings.Join([]string{i, p, strconv.Itoa(d1)}, ".")
					newresult := normalize(newi, d2)
					for v1, v2 := range newresult {
						result[v1] = v2
					}
				}
			default:
				// fmt.Printf("* %v is no idea\n\n", e)
				// fmt.Println(reflect.TypeOf(v))
				newi := strings.Join([]string{i, p}, ".")
				result[newi] = e
			}
		}
	case map[string]string:
		//fmt.Printf("** %v is map[string]string\n\n", k)
		for p, e := range v {
			newi := strings.Join([]string{i, p}, ".")
			result[newi] = e
		}
	case []map[string]string:
		//fmt.Printf("** %v is []map[string]string\n\n", k)
		for u, c := range v {
			newi := strings.Join([]string{i, strconv.Itoa(u)}, ".")
			newresult := normalize(newi, c)
			for v1, v2 := range newresult {
				result[v1] = v2
			}
		}
	case []map[string]interface{}:
		for u, c := range v {
			newi := strings.Join([]string{i, strconv.Itoa(u)}, ".")
			newresult := normalize(newi, c)
			for v1, v2 := range newresult {
				result[v1] = v2
			}
		}
	default: // if k is not a map then just return it as a string
		//fmt.Printf("** %v is no idea\n\n", k)
		result[i] = k
	}

	return result
}
