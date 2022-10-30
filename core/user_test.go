package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserMention(t *testing.T) {
	assert := assert.New(t)

	// Create user object
	user := &User{
		ID: "sample_id",
	}

	assert.Equal("<@sample_id>", user.Mention())

	// Test null user
	var nullUser *User = nil
	assert.NotEmpty(nullUser.Mention())

}
