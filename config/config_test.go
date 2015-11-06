/* package source_test is doing unit tests for the source package */
package config

import "testing"

// Stubing the Config

var fakeConf = Config{
	Name: "Testing",
	Strategy: Strategy{
		Typesource: "mysql",
		Tables:     map[string]string{"comments": "nyt_comments"}},
	Credentials: []Credential{Credential{
		Database: "coral",
		Username: "user",
		Password: "password",
		Host:     "host",
		Port:     "5432",
		Adapter:  "mysql",
		Type:     "source",
	}}}

// Not Testing New() *Config as it is only reading the file and unmarshalling it...

// GetCredential returns the credentials for connection with the external source
// Signature  GetCredential(adapter string) Credential
func TestGetCredential(t *testing.T) {
	adapter := "mysql"

	credential := fakeConf.GetCredential(adapter)

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
	strategy := fakeConf.GetStrategy()

	if strategy.Typesource != "mysql" {
		t.Error("Expected mysql, got ", strategy.Typesource)
	}
}

// GetTables returns a list of tables to be imported
// Signature GetTables() map[string]string {
func TestGetTables(t *testing.T) {
	tables := fakeConf.GetTables()

	if tables["comments"] != "nyt_comments" {
		t.Error("Expected nyt_comments, got ", tables["comments"])
	}
}
