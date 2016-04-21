package source

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
	"github.com/coralproject/sponge/pkg/webservice"
)

/* Implementing the Sources */

// API is the struct that has the connection string to the external mysql database
type API struct {
	Connection string
}

// GetData does the request to the webservice once and get back the data based on the parameters
func (a API) GetData(entity string, options *Options) ([]map[string]interface{}, error) { //offset int, limit int, orderby string, q string
	return nil, nil
}

// GetWebServiceData does the request to the webservice once and get back the data
func (a API) GetWebServiceData() ([]map[string]interface{}, bool, error) {

	notFinish := false
	var err error

	cred, err := strategy.GetCredential("api", "foreign")
	if err != nil {
		log.Error(uuid, "api.getwebservicedata", err, "Getting credentials with API")
	}

	credA, ok := cred.(str.CredentialService)
	if !ok {
		log.Error(uuid, "api.getwebservicedata", err, "Asserting type.")
	}

	userAgent := credA.GetUserAgent()
	method := "GET"

	response, err := webservice.DoRequest(uuid, userAgent, method, a.Connection, nil)
	if err != nil {
		log.Error(uuid, "api.getwebservicedata", err, "Calling Web Service %s.", a.Connection)
		return nil, notFinish, err
	}

	var d map[string]interface{}

	// Use json.Decode for reading streams of JSON data
	if err = json.NewDecoder(strings.NewReader(string(response.Body))).Decode(&d); err != nil {
		log.Error(uuid, "api.getwebservicedata", err, "Decoding data from API.")
	}

	recordsField := credA.GetRecordsFieldName()

	records, ok := d[recordsField].([]interface{}) //this are all the entries
	if !ok {
		log.Error(uuid, "api.getwebservicedata", err, "Asserting type.")
	}

	r := make([]map[string]interface{}, len(records))
	for _, i := range records { // all the entries in the type we need
		r = append(r, i.(map[string]interface{}))
	}

	flattenData, err := normalizeData(r)
	if err != nil {
		log.Error(uuid, "api.getwebservicedata", err, "Normalizing data from api to fit into fiddler.")
		return nil, notFinish, err
	}

	return flattenData, notFinish, err
}

// GetFireHoseData use the firehose to constantly GET data from the web service
func (a API) GetFireHoseData(pageAfter string) ([]map[string]interface{}, string, error) {
	var (
		flattenData   []map[string]interface{}
		nextPageAfter string
		err           error
	)

	// Get the credentials to connect to the API
	cred, err := strategy.GetCredential("api", "foreign")
	if err != nil {
		log.Error(uuid, "api.getFirehoseData", err, "Getting credentials with API")
	}

	// Assert Type into a credential API struct
	credA, ok := cred.(str.CredentialService)
	if !ok {
		log.Error(uuid, "api.getFirehoseData", err, "Asserting type.")
	}

	// TO DO: THIS IS VERY WAPO API HARCODED!
	// We need to insert pageAfter inside the q

	// Build the request
	req, err := http.NewRequest("GET", a.Connection, nil)
	if err != nil {
		log.Error(uuid, "api.getFirehoseData", err, "New request.")
		return nil, nextPageAfter, err
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
		log.Error(uuid, "api.getFirehoseData", err, "Doing a call to API.")
		return nil, nextPageAfter, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	var d map[string]interface{}

	// Use json.Decode for reading streams of JSON data
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Error(uuid, "api.getFirehoseData", err, "Decoding data from API.")
		return nil, nextPageAfter, err
	}

	recordsField := credA.GetRecordsFieldName()

	// Emtpy records means there are no more data to send
	if d[recordsField] == nil {
		return nil, pageAfter, err
	}

	records, ok := d[recordsField].([]interface{}) //this are all the entries
	if !ok {
		log.Error(uuid, "api.getFirehoseData", err, "Asserting type.")
		return nil, nextPageAfter, err
	}

	var r []map[string]interface{}
	for _, i := range records { // all the entries in the type we need
		r = append(r, i.(map[string]interface{}))
	}

	flattenData, err = normalizeData(r)
	if err != nil {
		log.Error(uuid, "api.getdata", err, "Normalizing data from api to fit into fiddler.")
		return nil, nextPageAfter, err
	}

	paginationField := credA.GetPaginationFieldName()

	nextPageAfter, ok = d[paginationField].(string) //strconv.ParseFloat(d["nextPageAfter"].(string), 64)
	if !ok {
		err = fmt.Errorf("Error when asserting type string.")
		log.Error(uuid, "api.getfirehosedata", err, "Type assigment to string")
		return nil, nextPageAfter, err
	}
	return flattenData, nextPageAfter, err
}

// GetQueryData will return all the data based on a specific list of IDs
func (a API) GetQueryData(entityname string, options *Options, ids []string) ([]map[string]interface{}, error) {
	var d []map[string]interface{}
	var err error

	return d, err
}

// IsWebService is a func from the Sourcer interface. It tell us if the external source is a database or API
func (a API) IsWebService() bool {
	return true
}

//////* Not exported functions *//////

// ConnectionMySQL returns the connection string
func connectionAPI() *url.URL {

	credA, ok := credential.(str.CredentialService)
	if !ok {
		err := fmt.Errorf("Error when asserting type CredentialService on credential.")
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

	return u
}
