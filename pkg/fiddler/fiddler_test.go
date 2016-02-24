package fiddler

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	uuidimported "github.com/pborman/uuid"
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

	u := uuidimported.New()

	// Initialize fiddler
	Init(u)
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

// Signature: TransformRow(row map[string]interface{}, modelName string) (interface{}, []map[string]interface{}, error) { // id row, transformation, error
func TestTransformRow(t *testing.T) {
	row := map[string]interface{}{"assetid": "3416344", "asseturl": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html", "updatedate": "2014-12-04 00:01:11", "createdate": "2014-12-04 00:01:11"}
	modelName := "assets"

	// interface{}, []map[string]interface{}, error)
	id, result, err := TransformRow(row, modelName)
	if err != nil {
		t.Fatalf("error should be nil. Error is %v", err)
	}

	expectedResult := make([]map[string]interface{}, 1)
	expectedResult[0] = map[string]interface{}{"date_updated": "2014-12-04T00:01:11Z", "date_created": "2014-12-04T00:01:11Z", "source": map[string]string{"id": "3416344"}, "url": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html"}

	if len(result) != len(expectedResult) {
		t.Fatalf("got %d , expected %d", len(result), len(expectedResult))
	}

	if result[0]["date_updated"] != expectedResult[0]["date_updated"] {
		t.Fatalf("got %s , expected %s", result[0]["date_updated"], expectedResult[0]["date_updated"])
	}

	if result[0]["url"] != expectedResult[0]["url"] {
		t.Fatalf("got %s , expected %s", result[0]["url"], expectedResult[0]["url"])
	}

	expectedID := "3416344"
	if id != expectedID {
		t.Fatalf("got %s, expected %s", id, expectedID)
	}
}

// Test array documents
// Signature: TransformRow(row map[string]interface{}, modelName string) (interface{}, []map[string]interface{}, error) { // id row, transformation, error
func TestTransformRowArrayTypesNotDuplicating(t *testing.T) {

	row := map[string]interface{}{
		"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
		"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
		"object.permalink": "",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
		"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
		"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
		"targets": []map[string]string{
			map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
				"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
		},
		"object.tags.0":        "replyto_Tropicat",
		"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
		"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
		"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
		"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
		"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
		"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
		"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
		"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
		"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
		"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
		"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
		"object.published": "2015-06-01T00:00:12Z",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
		"object.context.0.title": "", "actor.title": "Zeus Mom",
		"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
		"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
		"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
		"object.likes.0.published":           "2015-06-03T16:50:15Z",
	}
	modelName := "actions"

	// interface{}, []map[string]interface{}, error)
	id, result, err := TransformRow(row, modelName)
	if err != nil {
		t.Fatalf("error should be nil. Error is %v", err)
	}

	expectedResult := make([]map[string]interface{}, 1)
	expectedResult[0] = map[string]interface{}{
		"type": "likes", "target": "comments",
		"date": "2015-06-03T16:50:15Z",
		"source": map[string]interface{}{
			"user_id":   "http://washingtonpost.com/user0/",
			"target_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
		},
	}

	if len(result) != len(expectedResult) {
		t.Fatalf("got %d , expected %d", len(result), len(expectedResult))
	}

	expectedID := "ObjectIdHex(\"556ba08cd710290035cf6c74\")"
	if id != expectedID {
		t.Fatalf("got %s, expected %s", id, expectedID)
	}
}

// Test array documents
// Signature: TransformRow(row map[string]interface{}, modelName string) (interface{}, []map[string]interface{}, error) { // id row, transformation, error
func TestTransformRowArrayTypes(t *testing.T) {

	row := map[string]interface{}{
		"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
		"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
		"object.permalink": "",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
		"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
		"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
		"targets": []map[string]string{
			map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
				"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
		},
		"object.tags.0":        "replyto_Tropicat",
		"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
		"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
		"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
		"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
		"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
		"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
		"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
		"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
		"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
		"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
		"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
		"object.published": "2015-06-01T00:00:12Z",
		"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
		"object.context.0.title": "", "actor.title": "Zeus Mom",
		"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
		"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
		"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
		"object.likes.0.published":           "2015-06-03T16:50:15Z",
		"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
		"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
		"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
		"object.likes.1.published":           "2015-06-03T16:50:15Z",
	}
	modelName := "actions"

	// interface{}, []map[string]interface{}, error)
	id, result, err := TransformRow(row, modelName)
	if err != nil {
		t.Fatalf("error should be nil. Error is %v", err)
	}

	expectedResult := make([]map[string]interface{}, 2)
	expectedResult[0] = map[string]interface{}{
		"type": "likes", "target": "comments",
		"date": "2015-06-03T16:50:15Z",
		"source": map[string]interface{}{
			"user_id":   "http://washingtonpost.com/user0/",
			"target_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
		},
	}
	expectedResult[1] = map[string]interface{}{
		"type": "likes", "target": "comments",
		"date": "2015-06-03T16:50:15Z",
		"source": map[string]interface{}{
			"user_id":   "http://washingtonpost.com/user1/",
			"target_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
		},
	}

	if len(result) != len(expectedResult) {
		t.Fatalf("got %d , expected %d", len(result), len(expectedResult))
	}

	expectedID := "ObjectIdHex(\"556ba08cd710290035cf6c74\")"
	if id != expectedID {
		t.Fatalf("got %s, expected %s", id, expectedID)
	}

	// all the rows have different ids
	if result[0]["source"].(map[string]interface{})["user_id"] == result[1]["source"].(map[string]interface{})["user_id"] {
		t.Fatalf("It is duplicating documents")
	}
}

// Test there is an error when model does not exist.
// Signature: TransformRow(row map[string]interface{}, modelName string) (interface{}, []map[string]interface{}, error) { // id row, transformation, error
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
	modelName := "assets"
	expectedID := "assetid"

	id := GetID(modelName)
	if id != expectedID {
		t.Fatalf("got %s , expected %s", id, expectedID)
	}
}

// Signature:  GetCollections() []string {
func TestGetCollections(t *testing.T) {
	expectedCollections := []string{
		"assets",
		"users",
		"comments",
		"actions",
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
