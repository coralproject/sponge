package source

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	uuidimported "github.com/pborman/uuid"
)

const (
	cfgLoggingLevel = "LOGGING_LEVEL"
)

var (
	m   Sourcer
	mm  MySQL
	mdb MongoDB
	mp  PostgreSQL
)

var oStrategy string

func init() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel, log.Ldefault)
}

func setupMysql() {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_nyt_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	u := uuidimported.New()

	s, e := Init(u)
	if e != nil {
		fmt.Printf("Error when initializing strategy, %v.\n", e)
	}
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

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_wapo_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	u := uuidimported.New()
	s, e := Init(u)
	if e != nil {
		fmt.Printf("Error when initializing strategy, %v.\n", e)
	}

	m, e := New(s) // function being tested
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mdb, ok = m.(MongoDB)
	if !ok {
		fmt.Println("It should return a type MySQL")
	}
}

func setupPostgreSQL() {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	oStrategy = os.Getenv("STRATEGY_CONF")

	// MOCK STRATEGY CONF
	strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/strategy_discourse_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	var ok bool

	u := uuidimported.New()

	s, e := Init(u)
	if e != nil {
		fmt.Printf("Error when initializing strategy, %v.\n", e)
	}
	m, e := New(s) // function being tested
	if e != nil {
		fmt.Printf("Error when calling the function, %v.\n", e)
	}

	mp, ok = m.(PostgreSQL)
	if !ok {
		fmt.Println("It should return a type PostgreSQL")
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
func TestPostgreSQLNewSource(t *testing.T) {

	m, e := New("postgresql") // function being tested
	if e != nil {
		t.Fatalf("error when calling the function, %v.", e)
	}

	mp, ok := m.(PostgreSQL)
	if !ok {
		t.Fatalf("it should return a type PostgreSQL")
	}

	// m should have a valid connection string
	if mp.Connection == "" {
		t.Fatalf("connection string should not be nil")
	}

	// m should not have a database connection
	if mp.Database != nil {
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
