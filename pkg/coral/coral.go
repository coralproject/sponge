package coral

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
)

const methodGet string = "GET"
const methodPost string = "POST"

type restResponse struct {
	status  string
	header  http.Header
	payload string
}

// Config has the environment variables values
var Config struct {
	urlUser    string
	urlAsset   string
	urlComment string
}

func init() {
	var err error
	// look for configuration urls
	Config.urlUser, err = cfg.String("USER_URL")
	if err != nil {
		Config.urlUser = "http://localhost:8080/api/import/user"
		//log.Error("coral", "init", err, "Getting USER_URL env variable")
	}

	Config.urlAsset, err = cfg.String("ASSET_URL")
	if err != nil {
		Config.urlAsset = "http://localhost:8080/api/import/asset"
		//log.Error("coral", "init", err, "Getting ASSET_URL env variable")
	}

	Config.urlComment, err = cfg.String("COMMENT_URL")
	if err != nil {
		Config.urlComment = "http://localhost:8080/api/import/comment"
		//log.Error("coral", "init", err, "Getting COMMENT_URL env variable")
	}

}

// AddRow send the row to pillar based on which collection is
func AddRow(data []byte, modelName string) error {
	var err error
	switch modelName {
	case "user":
		err = addUser(data)
	case "comment":
		err = addComment(data)
	case "asset":
		err = addAsset(data)
	}

	return err
}

// AddComment  calls the API to add the comment
func addComment(jcomment []byte) error {

	err := doRequest(methodPost, Config.urlComment, bytes.NewBuffer(jcomment))
	if err != nil {
		log.Error("coral", "addComment", err, "Error on sending request with comment.")
		return err
	}

	return nil
}

// AddAssets calls the API to add the asset
func addAsset(jasset []byte) error {

	err := doRequest(methodPost, Config.urlAsset, bytes.NewBuffer(jasset))
	if err != nil {
		return err
	}

	return nil
}

// AddUsers calls the API to add the user
func addUser(juser []byte) error {

	err := doRequest(methodPost, Config.urlUser, bytes.NewBuffer(juser))
	if err != nil {
		return err
	}

	return nil
}

func doRequest(method string, urlStr string, payload io.Reader) error {

	request, err := http.NewRequest(method, urlStr, payload)
	if err != nil {
		log.Error("coral", "doRequest", err, "New request")
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error("coral", "doRequest", err, "Processing request")
	}
	defer response.Body.Close()

	return err
}
