package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	bigcache_store "github.com/eko/gocache/store/bigcache/v4"
)

type bigCache struct {
	c *cache.Cache[[]byte]
}

func (bc *bigCache) Set(key string, value any, options ...store.Option) error {
	valueBytes, mErr := json.Marshal(value)
	if mErr != nil {
		return mErr
	}
	return bc.c.Set(context.Background(), key, valueBytes, options...)
}

func (bc *bigCache) Get(key string, value any) error {
	valueBytes, gErr := bc.c.Get(context.Background(), key)
	if gErr != nil {
		return gErr
	}
	return json.Unmarshal(valueBytes, value)
}

func (bc *bigCache) Del(key string) error {
	return bc.c.Delete(context.Background(), key)
}

func (bc *bigCache) Clear(context context.Context) error {
	return bc.c.Clear(context)
}

func (bc *bigCache) Close() error {
	return bc.Clear(context.Background())
}

func Bigcache() Cache {
	bigcacheClient, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	bigcacheStore := bigcache_store.NewBigcache(bigcacheClient)
	return &bigCache{c: cache.New[[]byte](bigcacheStore)}
}
