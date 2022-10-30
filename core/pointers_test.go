package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	assert := assert.New(t)

	// Test const
	assert.True(*Bool(true))

	// Test var
	sample := true
	assert.True(*Bool(sample))

}
