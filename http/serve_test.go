package http

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeImg(t *testing.T) {
	img, err := ioutil.ReadFile("test_img/hello.png")
	assert.NoError(t, err)

	assert.NoError(t, ServeImage(img, 64, 48))
}
