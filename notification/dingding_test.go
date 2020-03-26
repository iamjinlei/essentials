package notification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDingding(t *testing.T) {
	t.Skip()
	assert.NoError(t, NewDingding("fill in webhook", time.Minute).Send("test", []string{"line1", "line2"}))
}
