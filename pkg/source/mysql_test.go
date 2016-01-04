/* package source_test is doing unit tests for the source package */
package source_test

import "testing"

// Mock this global configuration variables
// var strategy = strategy.Strategy{}
// var credential = strategy.Credential{}
// var config *strategy.Config

// NewSource returns a new MySQL struct
// Signature: NewSource() *MySQL
// It depends on the credentials to get the connection string
func TestNewSource(t *testing.T) {

	// credential.Database = "test"
	// credential.Username = "testuser"
	// credential.Password = "testpassword"
	// credential.Host = ""
	// credential.Port = ""
	// credential.Adapter = ""
	// credential.Type = "mysql"
	//
	// config.Name = "Test"
	// config.Strategy = strategy
	// config.Credentials[0] = credential
	//
	// var m *source.MySQL
	//
	// m = source.NewSource() // function being tested
	//
	// // m should have a valid connection string
	// if m.Connection == "" {
	// 	t.Error("Connection string should not be nil.")
	// }
	//
	// // m should not have a database connection
	// if m.Database != nil {
	// 	t.Error("Database should be nil.")
	// }
}

// GetData returns the raw data from the tableName
// Signature (m *MySQL) GetData(tableName string, modelName string) util.Data
func TestGetData(t *testing.T) {

	// credential.Database = "test"
	// credential.Username = "testuser"
	// credential.Password = "testpassword"
	// credential.Host = ""
	// credential.Port = ""
	// credential.Adapter = ""
	// credential.Type = "mysql"
	//
	// config.Name = "Test"
	// config.Strategy = strategy
	// config.Credentials[0] = credential
	//
	// var m *source.MySQL
	// m = source.NewSource()
	//
	// d := m.GetData("testModel") // function being tested
	//
	// if d.Error != nil {
	// 	t.Error("Error should be nil.")
	// }
	//
	// if d.Type != "testModel" {
	// 	t.Error("The type should be the model name.")
	// }
	//
	// // Check that d.Rows has the fixtures in it
}

// runQuery run the query on the Database
// Signature runQuery(db *sql.DB, model string, query string) ([]models.Model, error)
func TestRunQuery(t *testing.T) {
	//
	// credential.Database = "test"
	// credential.Username = "testuser"
	// credential.Password = "testpassword"
	// credential.Host = ""
	// credential.Port = ""
	// credential.Adapter = ""
	// credential.Type = "mysql"
	//
	// config.Name = "Test"
	// config.Strategy = strategy
	// config.Credentials[0] = credential
	//
	// var db *sql.DB //mock the db
	// model := "testModel"
	// query := "a query to test"
	//
	// var models []models.Model
	// var err error
	//
	// models, err = source.RunQuery(db, model, query)
	//
	// if err != nil {
	// 	t.Error("Error should be nil.")
	// }
	// // check models is not empty
	// if models[0] == nil {
	// 	t.Error("Models does not have models...")
	// }
}

// other func to test?
// connection
// open
// close
