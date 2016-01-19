package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/tests/fixtures/"
}

// GetFixture retrieves a query record from the filesystem for testing.
func GetFixture(fileName string) (map[string]interface{}, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var qs map[string]interface{}
	err = json.Unmarshal(content, &qs)
	if err != nil {
		return nil, err
	}

	return qs, nil
}
