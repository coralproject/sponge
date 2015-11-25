package models

import (
	"database/sql"
	"testing"

	"github.com/coralproject/sponge/config"
	"github.com/stretchr/testify/assert"
)

// Signature New(table string) models.Model
func TestNewAction(t *testing.T) {

	// The func we are testing
	c, err := New("Action")

	assert.Equal(t, Action{}, c, "They should be equal.")
	assert.Nil(t, err, "It should build an action without errors.")
}

// Signature (a Action) Transform(sd *sql.Rows, table config.Table) ([]Model, error)
func TestTransformAction(t *testing.T) {

	a, err := New("Action")

	sd := mockMysqlRowsActions()
	table := config.Table{}

	var actions []Model

	// The func we are testing
	actions, err = a.Transform(sd, table)

	assert.Nil(t, err, "It should not return any error.")

	assert.NotNil(t, actions, "It should return actions.")
}

// Something that implements sd *sql.Rows
func mockMysqlRowsActions() *sql.Rows {

	var sd *sql.Rows

	// To Do: Needs to return a valid sql.Rows stubbed
	return sd
}
