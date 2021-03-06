/* package source_test is doing unit tests for the source package */
package source

import "testing"

import "gopkg.in/mgo.v2/bson"

func TestMongoGetData(t *testing.T) {
	setupMongo()

	// Default Flags
	coralName := "comments"
	options := &Options{
		Offset:  0,
		Limit:   9999999999,
		Orderby: "",
		Query:   "",
	}

	// no error
	data, err := mdb.GetData(coralName, options)
	if err != nil {
		t.Fatalf("expected no error, got '%s'.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 150901
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("expected %d, got %d", expectedlen, len(data))
	}
}

func TestMongoQueryGetData(t *testing.T) {
	setupMongo()

	// Default Flags
	coralName := "comments"
	options := &Options{
		Offset:  0,
		Limit:   9999999999,
		Orderby: "",
		Query:   "{ \"updated\": { \"$gt\": \"2013-01-02 15:04:05\", \"$lt\": \"2016-01-02 15:04:05\" } }",
	}
	// no error
	data, err := mdb.GetData(coralName, options)
	if err != nil {
		t.Fatalf("expected no error, got '%s'.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 0
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("expected %d, got %d", expectedlen, len(data))
	}
}

// Signature func (m MongoDB) GetQueryData(coralTableName string, offset int, limit int, orderby string, ids []string) ([]map[string]interface{}, error) {
func TestMongoGetQueryData(t *testing.T) {
	setupMongo()

	// Default Flags
	coralName := "comments"
	options := &Options{
		Offset:  0,
		Limit:   9999999999,
		Orderby: ""}
	ids := []string{"56ac0c7010780b0a357bdec3", "56ac0c7010780b0a357bdec1", "56ac0c7010780b0a357bdec4"}

	// no error
	data, err := mdb.GetQueryData(coralName, options, ids)
	if err != nil {
		t.Fatalf("expected no error, got '%s'.", err)
	}

	// data should be []map[string]interface{}
	expectedlen := 0
	if len(data) != expectedlen { // this is a setup for the seed data
		t.Fatalf("expected %d, got %d", expectedlen, len(data))
	}
}

// Signature func (m MongoDB) GetEntities() ([]string, error) {
func TestMongoGetEntities(t *testing.T) {

	setupMongo()

	s, e := GetEntities()
	if e != nil {
		t.Fatalf("expected no error, got %s.", e)
	}

	expectedLen := 4
	if len(s) != expectedLen {
		t.Fatalf("got %d, it should be %d", len(s), expectedLen)
	}

	if s[0] != "users" {
		t.Fatalf("got %s, it should be asset", s[0])
	}

	if s[1] != "comments" {
		t.Fatalf("got %s, it should be comments", s[1])
	}

	teardown()
}

// Signature func flattenize(i string, k interface{}) (string, string) {
func Testflattenize(t *testing.T) {

	// when k is a map
	i := "firstkey"
	k1 := map[string]string{"secondkey1": "value1", "secondkey2": "value2"}

	m := flattenize(i, k1)

	if m["firstkey.secondkey1"] != "value1" {
		t.Error("Expected different map")
	}
	if m["firstkey.secondkey2"] != "value2" {
		t.Error("Expected different map")
	}

	// when k is a leaf
	i = "firstkey"
	k2 := "value"

	m = flattenize(i, k2)

	if m[i] != k2 {
		t.Error("Expected different map")
	}

	// when k is more complex
	i = "firstkey"
	n := map[string]string{
		"thirdkey1": "value21",
		"thirdkey2": "value22",
	}
	k3 := map[string]interface{}{
		"secondkey1": "value1",
		"secondkey2": n,
	}

	m = flattenize(i, k3)

	if m["firstkey.secondkey1"] != "value1" {
		t.Error("Expected different map")
	}

	if m["firstkey.secondkey2.thirdkey1"] != "value21" {
		t.Error("Expected different map")
	}

	if m["firstkey.secondkey2.thirdkey2"] != "value22" {
		t.Error("Expected different map")
	}
}

// Signature func (m MongoDB) flattenize(mongoData []map[string]interface{}) ([]map[string]interface{}, error) {
func TestflattenizeDocument(t *testing.T) {

	// Simple map[string]strings
	k1 := map[string]interface{}{"a": "1", "b": "2"}

	m, e := flattenizeDocument(k1)
	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if m["a"] != "1" {
		t.Error("Expected different map")
	}

	if len(m) != 2 {
		t.Errorf("Expected 2, got %d", len(m))
	}

	k2 := map[string]interface{}{"a": "1", "b": map[string]string{"c": "3", "d": "4"}}

	m, e = flattenizeDocument(k2)
	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if m["a"] != "1" {
		t.Error("Expected different map")
	}

	if m["b.c"] != "3" {
		t.Error("Expected different map")
	}

	if len(m) != 3 {
		t.Errorf("Expected 3, got %d", len(m))
	}

	k3 := map[string]interface{}{"a": map[string]string{"e": "5", "f": "6"}, "b": map[string]string{"c": "3", "d": "4"}}

	m, e = flattenizeDocument(k3)
	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if m["a.f"] != "6" {
		t.Error("Expected different map")
	}

	if len(m) != 4 {
		t.Errorf("Expected 4, got %d", len(m))
	}

	k4 := map[string]interface{}{"a": map[string]interface{}{"e": "5", "f": map[string]string{"g": "7", "h": "8"}}, "b": map[string]string{"c": "3", "d": "4"}}

	m, e = flattenizeDocument(k4)
	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if m["a.f.g"] != "7" {
		t.Error("Expected different map")
	}

	if m["b.c"] != "3" {
		t.Error("Expected different map")
	}

	if m["b.d"] != "4" {
		t.Error("Expected different map")
	}

	if len(m) != 5 {
		t.Errorf("Expected 5, got %d", len(m))
	}

	k5 := map[string]interface{}{
		"_id": bson.ObjectIdHex("556ba089d710290036ef099d"),
		"object": map[string]interface{}{
			"context": []map[string]interface{}{{
				"uri": "http://washingtonpost.com/opinions/reformers-want-to-erase-confuciuss-influence-in-asia-thats-a-mistake/2015/05/28/529c1d3a-042e-11e5-a428-c984eb077d4e_story.html",
			}},
			"published": "2015-06-01T00:00:09Z",
		},
	}

	m, e = flattenizeDocument(k5)
	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if m["_id"] != bson.ObjectIdHex("556ba089d710290036ef099d") {
		t.Errorf("Expected other map.")
	}

	if m["object.context.0.uri"] != "http://washingtonpost.com/opinions/reformers-want-to-erase-confuciuss-influence-in-asia-thats-a-mistake/2015/05/28/529c1d3a-042e-11e5-a428-c984eb077d4e_story.html" {
		t.Errorf("Expected other map.")
	}

	if m["object.published"] != "2015-06-01T00:00:09Z" {
		t.Errorf("Expected other map.")
	}

	k6 := map[string]interface{}{
		"_id": bson.ObjectIdHex("556ba08cd710290035cf6c74"),
		"object": map[string]interface{}{
			"context": []map[string]interface{}{{
				"uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			}},
			"published": "2015-06-01T00:00:12Z",
		},
	}

	m, e = flattenizeDocument(k6)
	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if m["_id"] != bson.ObjectIdHex("556ba08cd710290035cf6c74") {
		t.Errorf("Expected other map.")
	}
	if m["object.context.0.uri"] != "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/" {
		t.Errorf("Expected http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/. Got %v.", m["object.context.uri"])
	}

	if m["object.published"] != "2015-06-01T00:00:12Z" {
		t.Errorf("Expected other map.")
	}

}

// Signature: func (m MongoDB) flattenData(mongoData []map[string]interface{}) ([]map[string]interface{}, error)
func TestFlattenData(t *testing.T) {

	var k []map[string]interface{}

	k = []map[string]interface{}{
		map[string]interface{}{
			"_id": bson.ObjectIdHex("556ba089d710290036ef099d"),
			"object": map[string]interface{}{
				"context": []map[string]interface{}{{
					"uri": "http://washingtonpost.com/opinions/reformers-want-to-erase-confuciuss-influence-in-asia-thats-a-mistake/2015/05/28/529c1d3a-042e-11e5-a428-c984eb077d4e_story.html",
				}},
				"published": "2015-06-01T00:00:09Z",
			},
		},
		map[string]interface{}{
			"_id": bson.ObjectIdHex("556ba08cd710290035cf6c74"),
			"object": map[string]interface{}{
				"context": []map[string]interface{}{{
					"uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
				}},
				"published": "2015-06-01T00:00:12Z",
			},
		},
	}

	//fields := []string{"_id", "object.context.0.uri", "object.published"}
	//m, e := flattenData(fields, k)

	m, e := flattenizeData(k)

	if e != nil {
		t.Errorf("Expected not error, got %v", e)
	}

	if m[0]["_id"] != bson.ObjectIdHex("556ba089d710290036ef099d") {
		t.Errorf("Expected other map.")
	}

	if m[0]["object.context.0.uri"] != "http://washingtonpost.com/opinions/reformers-want-to-erase-confuciuss-influence-in-asia-thats-a-mistake/2015/05/28/529c1d3a-042e-11e5-a428-c984eb077d4e_story.html" {
		t.Errorf("Expected other map.")
	}

	if m[0]["object.published"] != "2015-06-01T00:00:09Z" {
		t.Errorf("Expected other map.")
	}

	if m[1]["_id"] != bson.ObjectIdHex("556ba08cd710290035cf6c74") {
		t.Errorf("Expected other map.")
	}
	if m[1]["object.context.0.uri"] != "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/" {
		t.Errorf("Expected http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/. Got %v.", m[1]["object.context.uri"])
	}

	if m[1]["object.published"] != "2015-06-01T00:00:12Z" {
		t.Errorf("Expected other map.")
	}

}
