package fiddler

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
)

func setup() {

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	// MOCK STRATEGY CONF
	strategyConf := "../../tests/strategy_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	// Initialize fiddler
	Init()
}

func teardown() {
}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

// Signature: TransformRow(row map[string]interface{}, modelName string) ([]byte, error)
func TestTransformRow(t *testing.T) {
	row := map[string]interface{}{"assetid": "3416344", "asseturl": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html", "updatedate": "2014-12-04 00:01:11"}
	modelName := "asset"

	result, err := TransformRow(row, modelName)
	if err != nil {
		t.Fatalf("error should be nil. Error is %v", err)
	}
	r := string(result)

	expectedResult := "{\"date_updated\":\"2014-12-04T00:01:11Z\",\"source\":{\"id\":\"3416344\"},\"url\":\"http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html\"}"

	if r != expectedResult {
		t.Fatalf("got %s , expected %s", r, expectedResult)
	}
}

// Test there is an error when model does not exist.
// Signature: TransformRow(row map[string]interface{}, modelName string) ([]byte, error)
func TestTransformRowNoModel(t *testing.T) {
	row := map[string]interface{}{}
	modelName := "papafrita"

	result, err := TransformRow(row, modelName)
	if err == nil {
		t.Fatalf("It should give an error")
	}

	if result != nil {
		t.Fatalf("It should give back no transformed row")
	}
}
