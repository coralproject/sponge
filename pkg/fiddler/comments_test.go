package fiddler

import (
	"testing"

	str "github.com/coralproject/sponge/strategy"
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

	c, err := New("Comment")

	sd := mockMysqlRowsComments()
	var table str.Table

	var comments []Transformer

	// The func we are testing
	comments, err = c.Transform(sd, table)

	assert.Nil(t, err, "It should not return any error.")
	assert.NotNil(t, comments, "It should return comments.")
}

// Something that implements sd *sql.Rows
func mockMysqlRowsComments() []map[string]interface{} {

	var sd []map[string]interface{}

	// To Do: Needs to return a valid sql.Rows stubbed
	return sd
}
