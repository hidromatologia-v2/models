package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testController(t *testing.T, c *Controller) {
	defer func(tt *testing.T) {
		assert.Nil(tt, c.Close())
	}(t)

}

func TestController(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testController(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testController(tt, c)
	})
}
