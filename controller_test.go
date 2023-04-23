package models

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/stretchr/testify/assert"
)

func testController(t *testing.T, c *Controller) {
	defer func(tt *testing.T) {
		assert.Nil(tt, c.Close())
	}(t)

}

func TestController(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
		defer c.Close()
		testController(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testController(tt, NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String())))
	})
}
