package fiddler

//
// import (
// 	"database/sql"
// 	"fmt"
// 	"time"
//
// 	"github.com/coralproject/sponge/config"
// )
//
// //* ACTIONS *//
//
// // Action has information on all the actions that a user performs on elements like recommends, likes, comments, notes, share, etc.
// type Action struct {
// 	ID         string    `json:"id" bson:"_id"`
// 	Type       string    `json:"type" bson:"type"`
// 	UserID     string    `json:"userid" bson:"userid"`
// 	Value      string    `json:"value" bson:"value"`
// 	CreateDate time.Time `json:"createdate" bson:"createdate"`
// 	UpdateDate time.Time `json:"updatedate" bson:"updatedate"`
// }
//
// // Print only print information about the action
// func (a Action) Print() {
// 	fmt.Println("Action: ", a.UserID, a.Type, a.Value)
// }
//
// // Transform get the data from sd
// func (a Action) Transform(sd *sql.Rows, table config.Table) ([]Model, error) {
// 	var action Action
// 	var actions []Model
//
// 	columns, e := sd.Columns()
// 	fmt.Printf("***************** COLUMNS: %s", columns) // Debugging
// 	if e != nil {
// 		fmt.Println(e)
// 	}
//
// 	var createDate string
// 	var updateDate string
//
// 	for sd.Next() {
// 		err := sd.Scan(&action.ID, &action.Type, &action.UserID, &action.Value, &createDate, &updateDate)
// 		if err != nil {
// 			return nil, scanError{error: err}
// 		}
//
// 		action.CreateDate, _ = time.Parse("2006-01-02", createDate) // To Do: I need to see how to dinamically discover what is the dateTime layout
// 		action.UpdateDate, _ = time.Parse("2006-01-02", updateDate)
//
// 		n := len(actions)
// 		if len(actions) == cap(actions) {
// 			// Comments is full and we must expand
// 			// Double the size and add 1
// 			newActions := make([]Model, len(actions), 2*len(actions)+1)
// 			copy(newActions, actions)
// 			actions = newActions
// 		}
// 		actions = actions[0 : n+1]
// 		actions[n] = action
// 	}
//
// 	return actions, nil
// }
