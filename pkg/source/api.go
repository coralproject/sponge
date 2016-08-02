// Package source API is used to connect to a web service to GET data
package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"text/template"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
	"github.com/coralproject/sponge/pkg/webservice"
)

/* Implementing the Sources */

// API is the struct that has the connection string to the external mysql database
type API struct {
	Connection string
}

// QueryData is used to get all the data about the different parameters of the query
type QueryData struct {
	Basicurl   string
	Attributes string
	Appkey     string
	Next       string
}

// GetData does the request to the webservice once and get back the data based on the parameters
// This is empty as it is not applicable to web services but it needs to be here to implement the Source interface
func (a API) GetData(entity string, options *Options) ([]map[string]interface{}, error) { //offset int, limit int, orderby string, q string
	var data []map[string]interface{}
	var err error

	return data, err
}

// GetWebServiceData does the request to the webservice once and get back the data
func (a API) GetWebServiceData() ([]map[string]interface{}, bool, error) {

	notFinish := false
	var err error

	cred, err := strategy.GetCredential("service", "foreign")
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

	flattenData, err := flattenizeData(r)
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
	cred, err := strategy.GetCredential("service", "foreign")
	if err != nil {
		log.Error(uuid, "api.getFirehoseData", err, "Getting credentials with API")
	}

	// Assert Type into a credential API struct
	credA, ok := cred.(str.CredentialService)
	if !ok {
		log.Error(uuid, "api.getFirehoseData", err, "Asserting type.")
	}

	// TO DO: THIS IS VERY WAPO API HARCODED!
	url := connectionAPI(pageAfter)

	log.User(uuid, "api.getFirehoseData", "Querying URL %s", url)
	// Build the request
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Error(uuid, "api.getFirehoseData", err, "New request.")
		return nil, pageAfter, err
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
		return nil, pageAfter, err
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf("Bad Request %v", resp.Status)
		log.Error(uuid, "api.getFirehoseData", err, "Doing a call to API.")
		return nil, pageAfter, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	var d map[string]interface{}

	// Use json.Decode for reading streams of JSON data
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Error(uuid, "api.getFirehoseData", err, "Decoding data from API.")
		return nil, pageAfter, err
	}

	recordsField := credA.GetRecordsFieldName()

	// Emtpy records means there are no more data to send
	if d[recordsField] == nil {
		return nil, pageAfter, err
	}

	records, ok := d[recordsField].([]interface{}) //this are all the entries
	if !ok {
		log.Error(uuid, "api.getFirehoseData", err, "Asserting type.")
		return nil, pageAfter, err
	}

	var r []map[string]interface{}
	for _, i := range records { // all the entries in the type we need
		r = append(r, i.(map[string]interface{}))
	}

	flattenData, err = flattenizeData(r)
	if err != nil {
		log.Error(uuid, "api.getfirehosedata", err, "Normalizing data from api to fit into fiddler.")
		return nil, nextPageAfter, err
	}

	nextPageField := credA.GetNextPageField()

	if d[nextPageField] != nil {
		switch pf := d[nextPageField].(type) {
		case float64:
			nextPageAfter = strconv.FormatFloat(pf, 'f', 6, 64)
		case string:
			nextPageAfter = pf
		default:
			log.User(uuid, "api.getfirehosedata", "Do not know what is the type asserting for the Next Pagination value.")
		}
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
func connectionAPI(pageAfter string) *url.URL {

	var u *url.URL

	credA, ok := credential.(str.CredentialService)
	if !ok {
		err := fmt.Errorf("Error when asserting type CredentialService on credential. It is type %v", reflect.TypeOf(credential))
		log.Error(uuid, "api.connectionAPI", err, "Asserting type.")
		return u
	}

	data := QueryData{}
	// Get all the data from credentials we need
	data.Basicurl = credA.GetEndpoint() //"https://comments-api.ext.nile.works/v1/search"
	data.Appkey = credA.GetAppKey()
	data.Next = pageAfter                                        //credA.GetPageAfterField()                        // field that we are going to get the next value from
	data.Attributes = credA.GetAttributes()                      // Attributes for the query. Eg, for WaPo we have scope and sortOrder
	urltemplate, urltemplatepagination := credA.GetQueryFormat() //format for the query
	regexToEscape := credA.GetRegexToEscape()

	//url.QueryEscape(

	surl, err := formatURL(data, urltemplate, urltemplatepagination, regexToEscape)
	if err != nil {
		log.Error(uuid, "api.connectionAPI", err, "Parsing url %s", surl)
		return u
	}

	u, err = url.Parse(surl)
	if err != nil {
		log.Error(uuid, "api.connectionAPI", err, "Parsing url %s", surl)
	}

	return u
}

//func doPrintf(format string, a []interface{}) string {
func formatURL(data QueryData, urltemplate string, urltemplatepagination string, regextoescape string) (string, error) {

	//	"queryformat": "{{basicurl}}?q=(({{attributes}}))&appkey={{appkey}}&nextSince={{next}}",

	var surl string
	var tmpl *template.Template
	var err error

	if data.Next != "" {
		tmpl, err = template.New("url").Parse(urltemplatepagination)
		if err != nil {
			log.Error(uuid, "formatURL", err, "Failing when parsing the query template.")
			return surl, err
		}

	} else {
		// this is the first time when we still do not have a page to query
		tmpl, err = template.New("url").Parse(urltemplate)
		if err != nil {
			log.Error(uuid, "formatURL", err, "Failing when parsing the query template.")
			return surl, err
		}
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Error(uuid, "formatURL", err, "Failing when parsing the query template.")
		return surl, err
	}
	surl = buf.String()

	// I need to escape the attributes...
	// First needs to look for the parameters to escape...
	re := regexp.MustCompile(regextoescape) //"\\(\\([A-Za-z0-9-&:/_. ]+\\w+\\)\\)")

	beginCharacters := "((" // Move it to Strategy
	endCharacters := "))"   // Move it to Strategy
	whattoescape := strings.Trim(re.FindString(surl), beginCharacters+endCharacters)
	replaceWith := fmt.Sprintf("%s%s%s", beginCharacters, url.QueryEscape(whattoescape), endCharacters)

	escapedURL := string(re.ReplaceAll([]byte(surl), []byte(replaceWith)))

	return escapedURL, nil
}
