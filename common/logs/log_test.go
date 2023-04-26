package logs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogOnError(t *testing.T) {
	t.Run("Nil", func(tt *testing.T) {
		LogOnError(nil)
	})
	t.Run("Error", func(tt *testing.T) {
		LogOnError(fmt.Errorf("an error"))
	})
}

func TestPanicOnError(t *testing.T) {
	t.Run("Nil", func(tt *testing.T) {
		PanicOnError(nil)
	})
	t.Run("Error", func(tt *testing.T) {
		defer func(ttt *testing.T) {
			assert.NotNil(ttt, recover())
		}(tt)
		PanicOnError(fmt.Errorf("an error"))
	})
}
