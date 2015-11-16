/* package source_test is doing unit tests for the source package */
package source

import "testing"

// NewSource returns a new MySQL struct
// Signature: NewSource() *MySQL
// It depends on the credentials to get the connection string
func TestNewSource(t *testing.T) {

	t.Skip("Pass")
	//m := NewSource()

	// m should be type MySQL

	// m should have a valid connection string
	// m should not have a database connection
}

// GetData returns the raw data from the tableName
// Signature (m *MySQL) GetData(tableName string) util.Data
func TestGetDatA(t *testing.T) {
	//d := m.GetData("test")
	t.Skip("Pass")
}

// GetNewData returns the data requested for the table
// Signature GetNewData() utils.Data {
func TestGetNewData(t *testing.T) {

	t.Skip("Pass")
	//
	// expectedType := type(utils.Data)
	// expectedData := [
	// 	{
	// 		CommentID: 1,
	// 		AssetID: 1,
	// 		StatusID : 0,
	// 		CommentTitle: "Titulo 1",
	// 		CommentBody: "Body 1 "
	// 	}
	// ]
	//
	// // fake Credentials and fake NewSource
	// d := GetNewData()
	//
	// // d has to be type utils.Data
	// if d.type != expectedType {
	// 	t.Fatalf("Expected type to be %s but it was %s", expectedType, d.type)
	// }
	//
	// // d has to has 2 Comments
	// if len(d.Comments) != len(expectedData) {
	// 	f.Fatalf("Expected only %d comments but there were %d.", len(expectedData), len(d.Comments))
	// }
	//
	// // d.error has to be nil
	// if d.Error != nil {
	// 	f.Fatalf("Expected to have no error but there was this one: %s .", d.Error)
	// }
}

/* ON HOW TO USE THIS PACKAGE */

// // ExampleMySQL on how to use the MySQL
// func ExampleMySQL() {
//
// 	// Creates a new mysql source
// 	m := NewSource()
//
// 	// Opens a connection to the database (it uses configuration package with a json configuration file)
// 	db := m.Open()
// 	// Closes connection when func finishes
// 	defer m.Close()
//
// 	// Returns rows related to that query
// 	d := m.Get(db, "SELECT * FROM nyt_comments LIMIT 1")
// 	for d.Next() {
// 		var comment string
// 		if d.Scan(&comment); err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("The comment is '%s'.", comment)
// 	}
// }
