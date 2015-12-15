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

	var fakeConf = Strategy{
		Name: "New York Times",
		Map: Map{
			Typesource: "mysql",
			Tables: map[string]Table{
				"Comment": Table{
					Name: "crnr_comment",
					Fields: map[string]string{
						"CommentID":    "commentid",
						"Body":         "commentbody",
						"ParentID":     "parentid",
						"AssetID":      "assetid",
						"Status":       "statusid",
						"DateCreated":  "createdate",
						"DateUpdated":  "updatedate",
						"DateApproved": "approvedate",
					},
				},
				"Asset": Table{
					Name: "crnr_asset",
					Fields: map[string]string{
						"AssetID":  "assetid",
						"SourceID": "sourceID",
						"URL":      "asseturl",
					},
				},
				"User": Table{
					Name: "crnr_comment",
					Fields: map[string]string{
						"UserID":      "userid",
						"UserName":    "userdisplayname",
						"Avatar":      "",
						"LastLogin":   "",
						"MemberSince": "",
						"TrustScore":  "",
					},
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

	if strategy.Typesource != "mysql" {
		t.Error("Expected mysql, got ", strategy.Typesource)
	}
}

// GetTables returns a list of tables to be imported
// Signature GetTables() map[string]string {
func TestGetTables(t *testing.T) {
	fakeConf := fakeStrategy()

	var tables map[string]Table
	tables = fakeConf.GetTables()

	if tables["Comment"].Name != "crnr_comment" {
		t.Error("Expected crnr_comment, got ", tables["Comment"])
	}
}
