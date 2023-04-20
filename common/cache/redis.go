package cache

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	redis_v9 "github.com/redis/go-redis/v9"
)

type redisCache struct {
	c      *cache.Cache[string]
	client *redis_v9.Client
}

func (rc *redisCache) Set(key string, value any, options ...store.Option) error {
	valueBytes, mErr := json.Marshal(value)
	if mErr != nil {
		return mErr
	}
	return rc.c.Set(context.Background(), key, hex.EncodeToString(valueBytes), options...)
}

func (rc *redisCache) Get(key string, value any) error {
	hexValue, gErr := rc.c.Get(context.Background(), key)
	if gErr != nil {
		return gErr
	}
	valueBytes, dErr := hex.DecodeString(hexValue)
	if dErr != nil {
		return dErr
	}
	return json.Unmarshal(valueBytes, value)
}

func (rc *redisCache) Del(key string) error {
	return rc.c.Delete(context.Background(), key)
}

func (rc *redisCache) Clear(context context.Context) error {
	return rc.c.Clear(context)
}

func (rc *redisCache) Close() error {
	return rc.client.Close()
}

func Redis(options *redis_v9.Options) Cache {
	client := redis_v9.NewClient(options)
	pErr := client.Ping(context.Background()).Err()
	if pErr != nil {
		panic(pErr)
	}
	redisStore := redis_store.NewRedis(client)
	return &redisCache{c: cache.New[string](redisStore), client: client}
}

func RedisDefault() Cache {
	return Redis(&redis_v9.Options{Addr: "127.0.0.1:6379"})
}
