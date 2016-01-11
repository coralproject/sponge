package source

import (
	"testing"
)

// TO DO: IT NEEDS TO MOCK STRATEGY AND CREDENTIAL!!!

// NewSource returns a new MySQL struct
// Signature: New(d string) (Sourcer, error) {
// It depends on the credentials to get the connection string
func TestMySQLNewSource(t *testing.T) {

	m, e := New("mysql") // function being tested
	if e != nil {
		t.Fatalf("error when calling the function, %v.", e)
	}

	mm, ok := m.(MySQL)
	// it returns type MySQL
	if !ok {
		t.Fatalf("it should return a type MySQL")
	}

	// m should have a valid connection string
	if mm.Connection == "" {
		t.Fatalf("connection string should not be nil")
	}

	// m should not have a database connection
	if mm.Database != nil {
		t.Error("database should be nil.")
	}
}
