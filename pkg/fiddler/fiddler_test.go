package fiddler

import (
	"testing"

	// Decisions: testify versus gocheck
	"github.com/stretchr/testify/assert"
)

// Test there is an error when model does not exist.
func TestErrorOnNew(t *testing.T) {
	// type Other does not exist
	c, err := New("Other")

	assert.Nil(t, c, "Any other type of model should return nil.")
	assert.NotNil(t, err, "Any other type of model should return an error.")
}

// The tests for the different models and in their own model_test file.
