package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Signature New(table string) models.Model {
func TestNew(t *testing.T) {
	// type Comments
	c, err := New("Comment")

	assert.Equal(t, c, Comment{}, "They should be equal.")
	assert.Nil(t, err)
}

func TestErrorOnNew(t *testing.T) {
	c, err := New("Other")

	assert.Nil(t, c)
	assert.NotNil(t, err)
}
