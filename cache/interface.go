package cache

import (
	"time"

	"github.com/fmyxyz/gocache/codec"
	"github.com/fmyxyz/gocache/store"
)

// CacheInterface represents the interface for all caches (aggregates, metric, memory, redis, ...)
type CacheInterface interface {
	Get(key interface{}) (interface{}, error)
	Set(key, object interface{}, options ...store.Option) error
	Delete(key interface{}) error
	Invalidate(options ...store.InvalidateOption) error
	Clear() error
	GetType() string
}

// SetterCacheInterface represents the interface for caches that allows
// storage (for instance: memory, redis, ...)
type SetterCacheInterface interface {
	CacheInterface
	GetWithTTL(key interface{}) (interface{}, time.Duration, error)

	GetCodec() codec.CodecInterface
}
