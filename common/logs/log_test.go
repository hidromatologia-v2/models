package logs

import (
	"fmt"
	"testing"
)

func TestLogOnError(t *testing.T) {
	t.Run("Nil", func(tt *testing.T) {
		LogOnError(nil)
	})
	t.Run("Error", func(tt *testing.T) {
		LogOnError(fmt.Errorf("an error"))
	})
}
