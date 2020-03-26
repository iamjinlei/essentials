package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	bytes, err := Get("https://github.com", nil)
	assert.NoError(t, err)
	t.Log(string(bytes))
}
