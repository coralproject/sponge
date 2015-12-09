package fiddler

import (
	"database/sql"
	"testing"

	"github.com/coralproject/sponge/config"
	"github.com/stretchr/testify/assert"
)

// Signature New(table string) models.Model
func TestNewComment(t *testing.T) {

	// The func we are testing
	c, err := New("Comment")

	assert.Equal(t, Comment{}, c, "They should be equal.")
	assert.Nil(t, err, "It should build a comment without errors.")
}

// Signature (c Comment) Transform(sd *sql.Rows, table config.Table) ([]Model, error)
func TestTransformComment(t *testing.T) {

	a, err := New("Comment")

	sd := mockMysqlRowsComments()
	table := config.Table{}

	var comments []Model

	// The func we are testing
	comments, err = a.Transform(sd, table)

	assert.Nil(t, err, "It should not return any error.")

	assert.NotNil(t, comments, "It should return comments.")
}

// Something that implements sd *sql.Rows
func mockMysqlRowsComments() *sql.Rows {

	var sd *sql.Rows

	// To Do: Needs to return a valid sql.Rows stubbed
	return sd
}
