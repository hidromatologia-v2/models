package models

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/stretchr/testify/assert"
)

func testController(t *testing.T, c *Controller) {
	defer func(tt *testing.T) {
		assert.Nil(tt, c.Close())
	}(t)
	defer c.Close()
}

func TestController(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testController(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testController(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}
