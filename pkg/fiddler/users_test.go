package fiddler

import (
	"database/sql"
	"testing"

	"github.com/coralproject/sponge/config"
	"github.com/stretchr/testify/assert"
)

// Signature New(table string) models.Model
func TestNewUser(t *testing.T) {

	// The func we are testing
	c, err := New("User")

	assert.Equal(t, User{}, c, "They should be equal.")
	assert.Nil(t, err, "It should build a User without errors.")
}

// Signature (c User) Transform(sd *sql.Rows, table config.Table) ([]Model, error)
func TestTransformUser(t *testing.T) {

	a, err := New("User")

	sd := mockMysqlRowsUsers()
	table := config.Table{}

	var Users []Model

	// The func we are testing
	Users, err = a.Transform(sd, table)

	assert.Nil(t, err, "It should not return any error.")

	assert.NotNil(t, Users, "It should return Users.")
}

// Something that implements sd *sql.Rows
func mockMysqlRowsUsers() *sql.Rows {

	var sd *sql.Rows

	// To Do: Needs to return a valid sql.Rows stubbed
	return sd
}
