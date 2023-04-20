package cache

import (
	"context"

	"github.com/eko/gocache/lib/v4/store"
)

type Cache interface {
	Set(key string, value any, options ...store.Option) error
	Get(key string, value any) error
	Del(key string) error
	Clear(context.Context) error
	Close() error
}
