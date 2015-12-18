package coral

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/server/model"
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

// AddData is the exposed func to import data into the coral db
func AddData(modelName string, data []byte) error {

	var err error
	//
	// fmt.Println("### DATA: ", string(data))

	switch modelName {
	case "user":
		err = addUsers(data)
		if err != nil {
			//log.Error("import", "AddUsers", err, "Send data to local database")
			fmt.Println(err)
		}
	case "comment":
		err = addComments(data)
		if err != nil {
			log.Error("import", "AddComments", err, "Send data to local database")
		}
	case "asset":
		err = addAssets(data)
		if err != nil {
			log.Error("import", "AddAssets", err, "Send data to local database")
		}
	}

	return err
}

// AddAssets calls the API to add assets
func addAssets(jassets []byte) error {

	objects := []model.Asset{}

	err := json.Unmarshal(jassets, &objects)
	if err != nil {
		log.Error("client", "addAssets", err, "Parsing asset data")
	}

	for _, asset := range objects {
		data, _ := json.Marshal(asset)

		err = doRequest(methodPost, Config.urlAsset, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}

	return nil
}

// AddUsers calls the API to add users
func addUsers(jusers []byte) error {

	objects := []model.User{}

	err := json.Unmarshal(jusers, &objects)
	if err != nil {
		//log.Error("client", "addUsers", err, "Parsing user data")
		fmt.Println(err)
	}

	for _, user := range objects {
		data, _ := json.Marshal(user)
		err = doRequest(methodPost, Config.urlUser, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}

	return nil
}

// AddComments  calls the API to add comments
func addComments(jcomments []byte) error {

	objects := []model.Comment{}

	err := json.Unmarshal(jcomments, &objects)
	if err != nil {
		fmt.Println(err)
		//log.Error("client", "addComment", err, "Parsing comments data")
	}

	for _, comment := range objects {
		data, _ := json.Marshal(comment)
		err = doRequest(methodPost, Config.urlComment, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func doRequest(method string, urlStr string, payload io.Reader) error {

	request, err := http.NewRequest(method, urlStr, payload)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		//	log.Error("client", "doRequest", err, "Processing request")
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()

	return nil
}
