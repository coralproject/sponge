package source

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
)

/* Implementing the Sources */

// API is the struct that has the connection string to the external mysql database
type API struct {
	Connection string
}

// while there is data --> GetData (nextPage)
// send data to channel

// process workers <--- get data from channel

// GetData does the request into the API and get back the data
// We are not using it because it is specific for a collection
func (a API) GetData(coralTableName string, offset int, limit int, orderby string, q string) ([]map[string]interface{}, bool, error) {

	notFinish := false
	var err error

	cred, err := strategy.GetCredential("api", "foreign")
	if err != nil {
		log.Error(uuid, "api.getdata", err, "Getting credentials with API")
	}

	credA, ok := cred.(str.CredentialAPI)
	if !ok {
		log.Error(uuid, "api.getdata", err, "Asserting type.")
	}

	// Build the request
	req, err := http.NewRequest("GET", a.Connection, nil)
	if err != nil {
		log.Error(uuid, "api.getdata", err, "New request.")
		return nil, notFinish, err
	}
	req.Header.Add("User-Agent", credA.GetUserAgent())

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Error(uuid, "api.getdata", err, "Doing a call to API.")
		return nil, notFinish, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	var d map[string]interface{}

	// Use json.Decode for reading streams of JSON data
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Error(uuid, "api.getdata", err, "Decoding data from API.")
	}

	recordsField := credA.GetRecordsFieldName()

	records, ok := d[recordsField].([]interface{}) //this are all the entries
	if !ok {
		log.Error(uuid, "api.getdata", err, "Asserting type.")
	}

	r := make([]map[string]interface{}, len(records))
	for _, i := range records { // all the entries in the type we need
		r = append(r, i.(map[string]interface{}))
	}

	flattenData, err := normalizeData(r)
	if err != nil {
		log.Error(uuid, "api.getdata", err, "Normalizing data from api to fit into fiddler.")
		return nil, notFinish, err
	}

	return flattenData, notFinish, err
}

// GetAPIData use the fire host to constantly give data from the API
func (a API) GetAPIData(pageAfter string) ([]map[string]interface{}, bool, string, error) {
	var (
		flattenData   []map[string]interface{}
		finish        bool
		nextPageAfter string
		err           error
	)

	finish = false

	// Get the credentials to connect to the API
	cred, err := strategy.GetCredential("api", "foreign")
	if err != nil {
		log.Error(uuid, "api.getdata", err, "Getting credentials with API")
	}

	// Assert Type into a credential API struct
	credA, ok := cred.(str.CredentialAPI)
	if !ok {
		log.Error(uuid, "api.getdata", err, "Asserting type.")
	}

	// TO DO: THIS IS VERY WAPO API HARCODED!
	// We need to insert pageAfter inside the q

	// Build the request
	req, err := http.NewRequest("GET", a.Connection, nil)
	if err != nil {
		log.Error(uuid, "api.getAPIData", err, "New request.")
		return nil, finish, nextPageAfter, err
	}
	req.Header.Add("User-Agent", credA.GetUserAgent())

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Error(uuid, "api.getdata", err, "Doing a call to API.")
		return nil, finish, nextPageAfter, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	var d map[string]interface{}

	// Use json.Decode for reading streams of JSON data
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Error(uuid, "api.getdata", err, "Decoding data from API.")
		return nil, finish, nextPageAfter, err
	}

	// // continue trying to get the data
	// if d["errorCode"] == "waiting" {
	// 	time.Sleep(10 * time.Second)
	// 	return nil, true, pageAfter, nil
	// }

	recordsField := credA.GetRecordsFieldName()

	// Emtpy records means there are no more data to send
	if d[recordsField] == nil {
		finish = true
		return nil, finish, pageAfter, err
	}

	records, ok := d[recordsField].([]interface{}) //this are all the entries
	if !ok {
		log.Error(uuid, "api.getdata", err, "Asserting type.")
		return nil, finish, nextPageAfter, err
	}

	//fmt.Printf("DEBUG entries: %s, length %v\n\n\n", recordsField, len(records))

	r := make([]map[string]interface{}, len(records))
	for _, i := range records { // all the entries in the type we need
		r = append(r, i.(map[string]interface{}))
	}

	flattenData, err = normalizeData(r)
	if err != nil {
		log.Error(uuid, "api.getdata", err, "Normalizing data from api to fit into fiddler.")
		return nil, finish, nextPageAfter, err
	}

	// TO DO: NEEDS TO MOVE THIS ATTRIBUTES TO THE STRATEGY FILE!!
	nextPageAfter, ok = d["nextPageAfter"].(string) //strconv.ParseFloat(d["nextPageAfter"].(string), 64)
	if !ok {
		finish = true
		//fmt.Println("DEBUG: no nextPageAfter: ", d)
	}

	return flattenData, finish, nextPageAfter, err
}

// GetQueryData will return all the data based on a specific list of IDs
func (a API) GetQueryData(coralTableName string, offset int, limit int, orderby string, ids []string) ([]map[string]interface{}, error) {
	var d []map[string]interface{}
	var err error

	return d, err
}

// IsAPI is a func from the Sourcer interface. It tell us if the external source is a database or API
func (a API) IsAPI() bool {
	return true
}

//////* Not exported functions *//////

// ConnectionMySQL returns the connection string
func connectionAPI() *url.URL {

	credA, ok := credential.(str.CredentialAPI)
	if !ok {
		err := fmt.Errorf("Error when asserting type credentialAPI on credential.")
		log.Error(uuid, "api.connectionAPI", err, "Asserting type.")
	}

	basicurl := credA.GetEndpoint()                      //"https://comments-api.ext.nile.works/v1/search"
	appkey := credA.GetAppKey()                          //"prod.washpost.com"
	attributes := url.QueryEscape(credA.GetAttributes()) // Attributes for the query. Eg, for WaPo we have scope and sortOrder

	surl := fmt.Sprintf("%s?q=((%s))&appkey=%s", basicurl, attributes, appkey)

	u, err := url.Parse(surl)
	if err != nil {
		log.Error(uuid, "api.connectionAPI", err, "Parsing url %s", surl)
	}

	// u, err := url.Parse(basicurl)
	// if err != nil {
	// 	log.Error(uuid, "api.connectionAPI", err, "Parsing basic url %s", basicurl)
	// }
	//
	// q := u.Query()
	// q.Set("q", attributes)
	// q.Set("appkey", appkey)
	//
	// u.RawQuery = q.Encode()

	return u
}
