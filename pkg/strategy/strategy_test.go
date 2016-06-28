/* package source_test is doing unit tests for the source package */
package strategy

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	uuidimported "github.com/pborman/uuid"
)

// Stubing the Config
func fakeStrategy() Strategy {

	var cdatabase CredentialDatabase
	cdatabase = CredentialDatabase{
		Database: "coral",
		Username: "user",
		Password: "password",
		Host:     "host",
		Port:     "5432",
		Adapter:  "mysql",
		Type:     "source",
	}

	cfields := make([]map[string]interface{}, 8)

	cfields[0] = map[string]interface{}{
		"foreign":  "commentid",
		"local":    "CommentID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[1] = map[string]interface{}{
		"foreign":  "commentbody",
		"local":    "Body",
		"relation": "Identity",
		"type":     "[]byte",
	}
	cfields[2] = map[string]interface{}{
		"foreign":  "parentid",
		"local":    "ParentID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[3] = map[string]interface{}{
		"foreign":  "assetid",
		"local":    "AssetID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[4] = map[string]interface{}{
		"foreign":  "statusid",
		"local":    "status",
		"relation": "Status",
		"type":     "string",
	}
	cfields[5] = map[string]interface{}{
		"foreign":        "createdate",
		"local":          "DateCreated",
		"relation":       "Parse",
		"type":           "timedate",
		"datetimeformat": "February 1st, 2006",
	}
	cfields[6] = map[string]interface{}{
		"foreign":  "updatedate",
		"local":    "DateUpdated",
		"relation": "Parse",
		"type":     "timedate",
	}
	cfields[7] = map[string]interface{}{
		"foreign":  "approvedate",
		"local":    "DateApproved",
		"relation": "Parse",
		"type":     "timedate",
	}

	afields := make([]map[string]interface{}, 3)
	afields[0] = map[string]interface{}{
		"foreign":  "assetid",
		"local":    "AssetID",
		"relation": "identity",
		"type":     "int",
	}
	afields[1] = map[string]interface{}{
		"foreign":  "sourceid",
		"local":    "SourceID",
		"relation": "identity",
		"type":     "int",
	}
	afields[2] = map[string]interface{}{
		"foreign":  "asseturl",
		"local":    "URL",
		"relation": "identity",
		"type":     "[]byte",
	}

	ufields := make([]map[string]interface{}, 6)
	ufields[0] = map[string]interface{}{
		"foreign":  "userid",
		"local":    "UserID",
		"relation": "identity",
		"type":     "int",
	}
	ufields[1] = map[string]interface{}{
		"foreign":  "userdisplayname",
		"local":    "UserName",
		"relation": "identity",
		"type":     "[]byte",
	}

	var status = map[string]string{
		"ModeratorApproved": "1",
		"Untouched":         "2",
	}
	var fakeConf = Strategy{
		Name: "New York Times",
		Map: Map{
			Foreign:        "mysql",
			DateTimeFormat: "2006-01-02 15:04:05",
			Entities: map[string]Entity{
				"comments": Entity{
					Foreign:        "crnr_comment",
					Local:          "comments",
					OrderBy:        "commentid",
					ID:             "commentid",
					Fields:         cfields,
					Status:         status,
					PillarEndpoint: "/api/import/comment",
				},
				"assets": Entity{
					Foreign:        "crnr_asset",
					Local:          "assets",
					OrderBy:        "assetid",
					ID:             "assetid",
					Fields:         afields,
					PillarEndpoint: "/api/import/asset",
				},
				"users": Entity{
					Foreign:        "crnr_comment",
					Local:          "users",
					OrderBy:        "userid",
					ID:             "commentid",
					Fields:         ufields,
					PillarEndpoint: "/api/import/user",
				},
			},
		},
		Credentials: Credentials{
			Database: cdatabase,
		},
	}

	return fakeConf
}

// Signature GetIDField(coralName string) string {
func TestGetID(t *testing.T) {
	fakeConf := fakeStrategy()
	modelName := "comments"

	id := fakeConf.GetIDField(modelName)

	if id != "commentid" {
		t.Fatalf("Expected commentid, got %v", id)
	}

}

// GetCredential returns the credentials for connection with the external source
// Signature  GetCredential(adapter string) Credential
func TestGetCredential(t *testing.T) {
	a := "mysql"
	ty := "source"
	fakeConf := fakeStrategy()

	credential, err := fakeConf.GetCredential(a, ty)
	if err != nil {
		t.Error("Expected not error, got ", err)
	}

	dcredential, ok := credential.(CredentialDatabase)
	if !ok {
		t.Error("Expected type credentialdtabase.")
	}

	// credential should have fields
	if dcredential.Database != "coral" {
		t.Error("Expected coral, got ", dcredential.Database)
	}

	if dcredential.Username != "user" {
		t.Error("Expected user, got ", dcredential.Username)
	}

	if dcredential.Password != "password" {
		t.Error("Expected password, got ", dcredential.Password)
	}

	if dcredential.Host != "host" {
		t.Error("Expected host, got ", dcredential.Host)
	}

	if dcredential.Port != "5432" {
		t.Error("Expected 5432, got ", dcredential.Port)
	}

	if dcredential.Adapter != "mysql" {
		t.Error("Expected mysql, got ", dcredential.Adapter)
	}

	if dcredential.Type != "source" {
		t.Error("Expected source, got ", dcredential.Type)
	}
}

// GetStrategy returns the strategy
// Signature GetStrategy() (Strategy, error)
func TestGetStrategy(t *testing.T) {
	fakeConf := fakeStrategy()

	strategy := fakeConf.GetMap()

	if strategy.Foreign != "mysql" {
		t.Error("Expected mysql, got ", strategy.Foreign)
	}
}

// GetTables returns a list of tables to be imported
// Signature GetTables() map[string]string {
func TestGetTables(t *testing.T) {
	fakeConf := fakeStrategy()

	var tables map[string]Entity
	tables = fakeConf.GetEntities()

	if tables["comments"].Foreign != "crnr_comment" {
		t.Error("Expected crnr_comment, got ", tables["comments"])
	}
}

// Signature func (s Strategy) HasArrayField(t Table) bool {
func TestHasArrayField(t *testing.T) {
	fakeConf := fakeStrategy()

	var tables map[string]Entity
	tables = fakeConf.GetEntities()

	if fakeConf.HasArrayField(tables["comments"]) {
		t.Error("Expected not to have an array field.")
	}
}

// Signature GetPillarEndpoints() map[string]string {
func TestGetPillarEndpoints(t *testing.T) {

	os.Setenv("PILLAR_URL", "http://localhost:8080")
	fakeConf := fakeStrategy()
	u := uuidimported.New()
	Init(u)

	endpoints := fakeConf.GetPillarEndpoints()

	expectedEndpoints := 4
	if len(endpoints) != expectedEndpoints { // the 3 on strategy plus create index
		t.Errorf("Expected %d endpoints, got %v", expectedEndpoints, len(endpoints))
	}

	expectedCommentEndpoint := "http://localhost:8080/api/import/comment"
	if endpoints["comments"] != expectedCommentEndpoint {
		t.Errorf("Expected %s, got %s", expectedCommentEndpoint, endpoints["comments"])
	}

	val, exists := endpoints["index"]
	if !exists {
		t.Errorf("Expected to have endpoint 'index'. Got %v.", val)
	}
}

func TestGetDatetimeFormat(t *testing.T) {
	fakeConf := fakeStrategy()

	table := "comments"

	// It should get the format for the strategy
	field := "DateUpdated"
	expectedDTformat := "2006-01-02 15:04:05"
	dtformat := fakeConf.GetDateTimeFormat(table, field)

	if dtformat != expectedDTformat {
		t.Errorf("Expected %s, got %s", expectedDTformat, dtformat)
	}

	// // It should get the format for the field
	field = "DateCreated"
	expectedDTformat = "February 1st, 2006"
	dtformat = fakeConf.GetDateTimeFormat(table, field)

	if dtformat != expectedDTformat {
		t.Errorf("Expected %s, got %s", expectedDTformat, dtformat)
	}

}

// Signature func (s Strategy) GetStatus(coralName string, foreign string) string {
func TestGetStatus(t *testing.T) {
	fakeConf := fakeStrategy()

	table := "comments"
	value := "ModeratorApproved"
	expectedsfield := "1"

	sfield := fakeConf.GetStatus(table, value)
	if sfield != expectedsfield {
		t.Errorf("Expected %s, got %s", expectedsfield, sfield)
	}
}

func TestGetEntityForeignName(t *testing.T) {
	fakeConf := fakeStrategy()
	collectionName := "comments"
	expectedForeigName := "crnr_comment"

	foreigName := fakeConf.GetEntityForeignName(collectionName)
	if foreigName != expectedForeigName {
		t.Errorf("Expected %s, got %s", expectedForeigName, foreigName)
	}
}

func TestValidation(t *testing.T) {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.NONE
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	u := uuidimported.New()
	Init(u)

	validstrategyfile := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/test/strategy_api_test.json"
	notvalidstrategyfile := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/test/not_valid_Strategy_test.json"

	if !Validate(validstrategyfile) {
		t.Errorf("Expected validation, did not validated")
	}

	if Validate(notvalidstrategyfile) {
		t.Errorf("Expected not validation, it validated")
	}

}
