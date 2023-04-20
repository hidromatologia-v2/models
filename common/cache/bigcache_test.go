package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigCache(t *testing.T) {
	c := Bigcache()
	defer c.Close()
	t.Run("Set", func(tt *testing.T) {
		defer c.Del("TEST")
		assert.Nil(tt, c.Set("TEST", "VALUE"))
		var value string
		assert.Nil(tt, c.Get("TEST", &value))
		assert.Equal(tt, "VALUE", value)
	})
	t.Run("Get", func(tt *testing.T) {
		defer c.Del("TEST")
		assert.Nil(tt, c.Set("TEST", "VALUE"))
		var value string
		assert.Nil(tt, c.Get("TEST", &value))
		assert.Equal(tt, "VALUE", value)
	})
	t.Run("Del", func(tt *testing.T) {
		assert.Nil(tt, c.Set("TEST", "VALUE"))
		var value string
		assert.Nil(tt, c.Get("TEST", &value))
		assert.Equal(tt, "VALUE", value)
		assert.Nil(tt, c.Del("TEST"))
	})
	t.Run("Clear", func(tt *testing.T) {
		defer c.Del("TEST")
		assert.Nil(tt, c.Set("TEST", "VALUE"))
		assert.Nil(tt, c.Clear(context.Background()))
		var value string
		assert.NotNil(tt, c.Get("TEST", &value))
	})
}
