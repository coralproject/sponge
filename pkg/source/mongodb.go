/*
Package source implements a way to get data from external MongoDB sources.
*/
package source

import (
	"strconv"
	"strings"

	"github.com/ardanlabs/kit/log"
	"gopkg.in/mgo.v2"
)

/* Implementing the Sources */

// MongoDB is the struct that has the connection string to the external mysql database
type MongoDB struct {
	Connection string
	Database   *mgo.Session
}

/* Exported Functions */

// GetTables gets all the collections names from this data source
func (m MongoDB) GetTables() ([]string, error) {
	keys := []string{}
	for k := range strategy.Map.Tables {
		keys = append(keys, k)
	}
	return keys, nil
}

// GetData returns the raw data from the tableName
func (m MongoDB) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error) { //(*sql.Rows, error) {

	var dat []map[string]interface{}

	// Get the corresponding table to the modelName
	collectionName := strategy.GetTableForeignName(coralTableName)
	fields := strategy.GetTableForeignFields(coralTableName) //[]]map[string]string

	// open a connection
	session, err := m.initSession()
	if err != nil {
		log.Error("Importing", "GetData", err, "Init mongo session")
		return nil, err
	}
	defer m.closeSession(session)

	cred := mgo.Credential{
		Username: credential.Username,
		Password: credential.Password,
	}

	err = session.Login(&cred)
	if err != nil {
		log.Error("Importing", "GetData", err, "Login mongo session")
		return nil, err
	}

	db := session.DB(credential.Database)
	col := db.C(collectionName)

	//Get all the fields that we are going to get from the document { field: 1}
	query := make(map[string]bool)
	for _, f := range fields {
		query[f["foreign"]] = true
	}

	//.Select(query)
	err = col.Find(nil).Limit(limit).All(&dat) // TO DO: CHECK TO SEE IF THERE IS ANY OTHER WAY TO GET THIS
	if err != nil {
		log.Error("Importing", "GetData", err, "Get collection")
		return nil, err
	}

	datsm, err := normalizeData(dat)
	if err != nil {
		log.Error("Importing", "GetData", err, "Converting into simple structure")
		return nil, err
	}

	return datsm, nil
}

// GetQueryData needs to be implemented for mongodb to implement the sourcer interface
func (m MongoDB) GetQueryData(string, int, int, string, []string) ([]map[string]interface{}, error) {
	return nil, nil
}

//////* Not exported functions *//////

// ConnectionMongoDB returns the connection string
func connectionMongoDB() string {
	return credential.Username + ":" + credential.Password + "@" + "/" + credential.Database
}

// Open gives back a pointer to the DB
func (m *MongoDB) initSession() (*mgo.Session, error) {

	database, err := mgo.Dial(m.Connection)
	if err != nil {
		log.Error("Importing", "initSession", err, "Dial into session")
		return nil, err
	}

	m.Database = database

	return database, nil
}

// Close closes the db
func (m MongoDB) closeSession(session *mgo.Session) error {
	session.Close()
	return nil
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
