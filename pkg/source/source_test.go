package source

import (
	"fmt"
	"os"
	"testing"
)

var m Sourcer
var mmysql MySQL
var mmongo MongoDB

var oStrategy string

func setupMysq() {
	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_mysql_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}
}

func setupMongo() {
	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_mongo_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}
}

func setup(source string) {

	oStrategy = os.Getenv("STRATEGY_CONF")

	if source == "mysql" {
		setupMysq()
	}

	if source == "mongo" {
		setupMongo()
	}

	var ok bool

	Init()

	m, e := New("mysql") // function being tested
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mmysql, ok = m.(MySQL)
	if !ok {
		fmt.Println("It should return a type MySQL")
	}

	m, e = New("mongodb")
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mmongo, ok = m.(MongoDB)
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

// func TestMain(m *testing.M) {
// 	setup()
// 	code := m.Run()
// 	teardown()
//
// 	os.Exit(code)
// }

// TO DO: IT NEEDS TO MOCK STRATEGY AND CREDENTIAL!!!

// NewSource returns a new MySQL struct
// Signature: New(d string) (Sourcer, error) {
// It depends on the credentials to get the connection string
func TestMySQLNewSource(t *testing.T) {

	setup("mysql")

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

	teardown()
}

//Signature: (m MySQL) GetTables() ([]string, error)
func TestGetTables(t *testing.T) {

	setup("mysql")

	s, e := GetTables()
	if e != nil {
		t.Fatalf("expected no error, got %s.", e)
	}

	expectedLen := 3
	if len(s) != expectedLen {
		t.Fatalf("got %d, it should be %d", len(s), expectedLen)
	}

	teardown()
}
