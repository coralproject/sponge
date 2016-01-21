/* package source_test is doing unit tests for the source package */
package strategy

import "testing"

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

	cfields := make([]map[string]string, 8)

	cfields[0] = map[string]string{
		"foreign":  "commentid",
		"local":    "CommentID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[1] = map[string]string{
		"foreign":  "commentbody",
		"local":    "Body",
		"relation": "Identity",
		"type":     "[]byte",
	}
	cfields[2] = map[string]string{
		"foreign":  "parentid",
		"local":    "ParentID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[3] = map[string]string{
		"foreign":  "assetid",
		"local":    "AssetID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[4] = map[string]string{
		"foreign":  "statusid",
		"local":    "Status",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[5] = map[string]string{
		"foreign":  "createdate",
		"local":    "DateCreated",
		"relation": "Parse",
		"type":     "timedate",
	}
	cfields[6] = map[string]string{
		"foreign":  "updatedate",
		"local":    "DateUpdated",
		"relation": "Parse",
		"type":     "timedate",
	}
	cfields[7] = map[string]string{
		"foreign":  "approvedate",
		"local":    "DateApproved",
		"relation": "Parse",
		"type":     "timedate",
	}

	afields := make([]map[string]string, 3)
	afields[0] = map[string]string{
		"foreign":  "assetid",
		"local":    "AssetID",
		"relation": "identity",
		"type":     "int",
	}
	afields[1] = map[string]string{
		"foreign":  "sourceid",
		"local":    "SourceID",
		"relation": "identity",
		"type":     "int",
	}
	afields[2] = map[string]string{
		"foreign":  "asseturl",
		"local":    "URL",
		"relation": "identity",
		"type":     "[]byte",
	}

	ufields := make([]map[string]string, 6)
	ufields[0] = map[string]string{
		"foreign":  "userid",
		"local":    "UserID",
		"relation": "identity",
		"type":     "int",
	}
	ufields[1] = map[string]string{
		"foreign":  "userdisplayname",
		"local":    "UserName",
		"relation": "identity",
		"type":     "[]byte",
	}

	var fakeConf = Strategy{
		Name: "New York Times",
		Map: Map{
			Foreign: "mysql",
			Tables: map[string]Table{
				"Comment": Table{
					Foreign:  "crnr_comment",
					Local:    "comment",
					OrderBy:  "commentid",
					ID:       "commentid",
					Fields:   cfields,
					Endpoint: "http://localhost:8080/api/import/comment",
				},
				"Asset": Table{
					Foreign:  "crnr_asset",
					Local:    "asset",
					OrderBy:  "assetid",
					ID:       "assetid",
					Fields:   afields,
					Endpoint: "http://localhost:8080/api/import/asset",
				},
				"User": Table{
					Foreign:  "crnr_comment",
					Local:    "user",
					OrderBy:  "userid",
					ID:       "commentid",
					Fields:   ufields,
					Endpoint: "http://localhost:8080/api/import/user",
				},
			},
		},
		Credentials: Credentials{
			Databases: cdatabases,
		},
	}

	return fakeConf
}

// Not Testing New() *Config as it is only reading the file and unmarshalling it...

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

	credential := fakeConf.GetCredential(a, ty)

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

// Signature GetPillarEndpoints() map[string]string {
func TestGetPillarEndpoints(t *testing.T) {
	fakeConf := fakeStrategy()
	endpoints := fakeConf.GetPillarEndpoints()

	if len(endpoints) != 3 {
		t.Errorf("Expected 3 endpoints, got %v", len(endpoints))
	}
	if endpoints["comment"] != "http://localhost:8080/api/import/comment" {
		t.Errorf("Expected http://localhost:8080/api/import/comment, got %s", endpoints["Comment"])
	}
}
