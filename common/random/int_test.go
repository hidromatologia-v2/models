package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	assert.NotEqual(t, Int(1000), Int(1000))
}
