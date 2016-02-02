package source

import (
	"fmt"
	"os"
	"testing"
)

var (
	m   Sourcer
	mm  MySQL
	mdb MongoDB
)

var oStrategy string

func setupMysql() {

	//setup environment variable
	// strategy File
	// create mysql db

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_nyt_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	s := Init()
	m, e := New(s) // function being tested
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mm, ok = m.(MySQL)
	if !ok {
		fmt.Println("It should return a type MySQL")
	}
}

func setupMongo() {

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_wapo_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	s := Init()

	m, e := New(s) // function being tested
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mdb, ok = m.(MongoDB)
	if !ok {
		fmt.Println("It should return a type MySQL")
	}
}

func teardown() {
	e := os.Setenv("STRATEGY_CONF", oStrategy)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}
}

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

// NewSource returns a new mongodb struct
// Signature: New(d string) (Sourcer, error) {
// It depends on the credentials to get the connection string
func TestMongoDBNewSource(t *testing.T) {

	m, e := New("mongodb") // function being tested
	if e != nil {
		t.Fatalf("error when calling the function, %v.", e)
	}

	mdb, ok := m.(MongoDB)
	// it returns type MySQL
	if !ok {
		t.Fatalf("it should return a type MySQL")
	}

	// m should have a valid connection string
	if mdb.Connection == "" {
		t.Fatalf("connection string should not be nil")
	}

	// m should not have a database connection
	if mdb.Database != nil {
		t.Error("database should be nil.")
	}
}
