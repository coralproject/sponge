package models

import (
	"database/sql"
	"testing"

	"github.com/coralproject/sponge/config"
	"github.com/stretchr/testify/assert"
)

// Signature New(table string) models.Model
func TestNewNote(t *testing.T) {

	// The func we are testing
	c, err := New("Note")

	assert.Equal(t, Note{}, c, "They should be equal.")
	assert.Nil(t, err, "It should build a note without errors.")
}

// Signature (c Note) Transform(sd *sql.Rows, table config.Table) ([]Model, error)
func TestTransformNote(t *testing.T) {

	a, err := New("Note")

	sd := mockMysqlRowsNotes()
	table := config.Table{}

	var notes []Model

	// The func we are testing
	notes, err = a.Transform(sd, table)

	assert.Nil(t, err, "It should not return any error.")

	assert.NotNil(t, notes, "It should return notes.")
}

// Something that implements sd *sql.Rows
func mockMysqlRowsNotes() *sql.Rows {

	var sd *sql.Rows

	// To Do: Needs to return a valid sql.Rows stubbed
	return sd
}
