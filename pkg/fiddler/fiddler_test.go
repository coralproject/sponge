package fiddler

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
)

var (
	oStrategy  string
	oPillarURL string
)

func setup() {

	// Save original enviroment variables
	oStrategy = os.Getenv("STRATEGY_CONF")
	oPillarURL = os.Getenv("PILLAR_URL")

	// Initialize log
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	// Mock strategy configuration
	strategyConf := "../../tests/strategy_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	// Initialize fiddler
	Init()
}

func teardown() {

	// recover the environment variables

	os.Setenv("STRATEGY_CONF", oStrategy)
	os.Setenv("PILLAR_URL", oPillarURL)
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

	id, result, err := TransformRow(row, modelName)
	if err != nil {
		t.Fatalf("error should be nil. Error is %v", err)
	}

	expectedResult := map[string]interface{}{"date_updated": "2014-12-04T00:01:11Z", "source": map[string]string{"id": "3416344"}, "url": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html"}

	if len(result) != len(expectedResult) {
		t.Fatalf("got %d , expected %d", len(result), len(expectedResult))
	}

	if result["date_updated"] != expectedResult["date_updated"] {
		t.Fatalf("got %s , expected %s", result["date_updated"], expectedResult["date_updated"])
	}

	if result["url"] != expectedResult["url"] {
		t.Fatalf("got %s , expected %s", result["url"], expectedResult["url"])
	}

	expectedID := "3416344"
	if id != expectedID {
		t.Fatalf("got %s, expected %s", id, expectedID)
	}
}

// Test there is an error when model does not exist.
// Signature: TransformRow(row map[string]interface{}, modelName string) ([]byte, error)
func TestTransformRowNoModel(t *testing.T) {
	row := map[string]interface{}{}
	modelName := "papafrita"

	_, result, err := TransformRow(row, modelName)
	if err == nil {
		t.Fatalf("It should give an error")
	}

	if result != nil {
		t.Fatalf("It should give back no transformed row")
	}
}

// Signature:  GetID(modelName string) string
func TestGetID(t *testing.T) {
	modelName := "asset"
	expectedID := "assetid"

	id := GetID(modelName)
	if id != expectedID {
		t.Fatalf("got %s , expected %s", id, expectedID)
	}
}

// Signature:  GetCollections() []string {
func TestGetCollections(t *testing.T) {
	expectedCollections := []string{
		"asset",
		"user",
		"comment",
	}

	collections := GetCollections()

	if !equal(collections, expectedCollections) {
		t.Fatalf("got %s , expected %s", collections, expectedCollections)
	}
}

func equal(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if vnotinb(v, b) {
			return false
		}
	}

	return true
}

func vnotinb(v string, b []string) bool {

	for _, k := range b {
		if v == k {
			return false
		}
	}
	return true
}

// Signature: appendField(source []map[string]interface{}, item interface{}) []map[string]interface{}
func TestappendField(t *testing.T) {

	var source []map[string]interface{}

	source[0] = make(map[string]interface{})
	source[0]["asset_id"] = 1
	source[1] = make(map[string]interface{})
	source[1]["comment_id"] = 2

	var item map[string]int

	item = make(map[string]int)
	item["user_id"] = 3

	result := appendField(source, item)

	if result[3]["user_id"] == 3 {
		t.Fatalf("got  %v, expected 3", result[3]["user_id"])
	}
}
