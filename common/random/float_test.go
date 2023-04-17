package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat(t *testing.T) {
	assert.NotEqual(t, Float(1000.0), Float(1000.0))
}
