/* package source_test is doing unit tests for the source package */
package strategy

import (
	"os"
	"testing"

	uuidimported "github.com/pborman/uuid"
)

// Stubing the Config
func fakeStrategy() Strategy {

	var cdatabases []CredentialDatabase
	cdatabases = make([]CredentialDatabase, 2)
	cdatabases[0] = CredentialDatabase{
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
			Tables: map[string]Table{
				"Comment": Table{
					Foreign:  "crnr_comment",
					Local:    "comment",
					OrderBy:  "commentid",
					ID:       "commentid",
					Fields:   cfields,
					Status:   status,
					Endpoint: "/api/import/comment",
				},
				"Asset": Table{
					Foreign:  "crnr_asset",
					Local:    "asset",
					OrderBy:  "assetid",
					ID:       "assetid",
					Fields:   afields,
					Endpoint: "/api/import/asset",
				},
				"User": Table{
					Foreign:  "crnr_comment",
					Local:    "user",
					OrderBy:  "userid",
					ID:       "commentid",
					Fields:   ufields,
					Endpoint: "/api/import/user",
				},
			},
		},
		Credentials: Credentials{
			Databases: cdatabases,
		},
	}

	return fakeConf
}

// Signature GetIDField(coralName string) string {
func TestGetID(t *testing.T) {
	fakeConf := fakeStrategy()
	modelName := "Comment"

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

	// credential should have fields
	if credential.Database != "coral" {
		t.Error("Expected coral, got ", credential.Database)
	}

	if credential.Username != "user" {
		t.Error("Expected user, got ", credential.Username)
	}

	if credential.Password != "password" {
		t.Error("Expected password, got ", credential.Password)
	}

	if credential.Host != "host" {
		t.Error("Expected host, got ", credential.Host)
	}

	if credential.Port != "5432" {
		t.Error("Expected 5432, got ", credential.Port)
	}

	if credential.Adapter != "mysql" {
		t.Error("Expected mysql, got ", credential.Adapter)
	}

	if credential.Type != "source" {
		t.Error("Expected source, got ", credential.Type)
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

	var tables map[string]Table
	tables = fakeConf.GetTables()

	if tables["Comment"].Foreign != "crnr_comment" {
		t.Error("Expected crnr_comment, got ", tables["Comment"])
	}
}

// Signature func (s Strategy) HasArrayField(t Table) bool {
func TestHasArrayField(t *testing.T) {
	fakeConf := fakeStrategy()

	var tables map[string]Table
	tables = fakeConf.GetTables()

	if fakeConf.HasArrayField(tables["Comment"]) {
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
	if endpoints["comment"] != expectedCommentEndpoint {
		t.Errorf("Expected %s, got %s", expectedCommentEndpoint, endpoints["Comment"])
	}

	val, exists := endpoints["index"]
	if !exists {
		t.Errorf("Expected to have endpoint 'index'. Got %v.", val)
	}
}

func TestGetDatetimeFormat(t *testing.T) {
	fakeConf := fakeStrategy()

	table := "Comment"

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

	table := "Comment"
	value := "ModeratorApproved"
	expectedsfield := "1"

	sfield := fakeConf.GetStatus(table, value)
	if sfield != expectedsfield {
		t.Errorf("Expected %s, got %s", expectedsfield, sfield)
	}
}
