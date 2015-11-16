package models

import (
	"database/sql"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Signature New(table string) models.Model {
func TestNew(t *testing.T) {
	// type Comments
	c, err := New("Comment")

	assert.Equal(c, Comment{}, "They should be equal.")
	assert.Nil(err)
}

// Test there is an error when model does not exist.
func TestErrorOnNew(t *testing.T) {
	// type Other does not exist
	c, err := New("Other")

	assert.Nil(c)
	assert.NotNil(err)
}

// Signature (c Comment) ProcessData(sd *sql.Rows) ([]Model, error) {
func TestCommentProcessData(t *testing.T) {

	c, err := New("Comment")
	//sd := mockMysql()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	// expect query to fetch order and user, match it with regexp
	mock.ExpectQuery("SELECT * from Comments").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1"))
	// expect transaction rollback, since order status is "cancelled"
	mock.ExpectRollback()

	comments, err = c.ProcessData(sd)

	assert.Nil(err)
	assert.NonNil(comments)
}

//* Utility Functions *//

type MyMockedSD struct {
	mock.Mock
}

func mockMysql() *sql.Rows {
	// To Do: this needs to be mocked
	// var mysql *source.MySQL
	// mysql = source.NewSource()
	//
	// db, err := mysql.open()
	// defer mysql.close(db)
	//
	// db, err := m.open()
	// defer m.close(db)
	//
	// sd, err := db.Query(query)
	//
	// defer sd.Close()

	sd := sql.Rows{}

	return sd
}
