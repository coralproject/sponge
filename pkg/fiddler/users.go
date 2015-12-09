package fiddler

import (
	"fmt"
	"log"
	"time"

	"github.com/coralproject/shelf/pkg/srv/comment"
	configuration "github.com/coralproject/sponge/config"
	"github.com/oleiade/reflections"
)

//User is embedding the comment package to extend it
type User struct {
	comment.User
}

// Print only print information about the user
func (u User) Print() {
	fmt.Println("User: ", u.UserID, u.UserName)
}

// Transform get the data from sd
func (u User) Transform(sd []map[string]interface{}, table configuration.Table) ([]Transformer, error) {

	var user User
	var users []Transformer

	// To Do: it needs refactoring as my gut tells me that is quite inefficient
	for _, value := range sd {
		for coralField, f := range table.Fields {

			// convert field f with value value[f] into field coralField
			if f != "" {
				newValue := transformUserField(f, value[f], coralField)

				err := reflections.SetField(&user, coralField, newValue)
				if err != nil {
					log.Fatal(err)
					return nil, err
				}
			}
		}

		n := len(users)
		if len(users) == cap(users) {
			// Comments is full and we must expand
			// Double the size and add 1
			newUsers := make([]Transformer, len(users), 2*len(users)+1)
			copy(newUsers, users)
			users = newUsers
		}
		users = users[0 : n+1]
		//comment.Raw = strings.Split(raws, ",")
		users[n] = user
	}

	return users, nil
}

//Here we transform the record into what we want (based on the configuration)
// 1. convert types (values are all strings) into the struct
func transformUserField(sourceField string, oldValue interface{}, coralField string) interface{} {

	var newValue interface{}

	const longForm = "2015-11-02 12:26:05"

	// Right now this the simpler thing to do. Needs to look more into reflect to do it
	// dinamically (if it is worth it) or merge this into the comment package
	switch coralField {
	case "UserID": //string
		newValue = oldValue
	case "UserName": //string
		newValue = oldValue
	case "Avatar": //string
		newValue = oldValue
	case "LastLogin": //time.Time
		newValue, _ = time.Parse(longForm, oldValue.(string))
	case "MemberSince": //time.Time
		newValue, _ = time.Parse(longForm, oldValue.(string))
	case "TrustScore": //float64
		newValue = oldValue
	}

	return newValue
}
