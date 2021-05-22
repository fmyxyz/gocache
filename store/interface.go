package store

import (
	"time"
)

// StoreInterface is the interface for all available stores
type StoreInterface interface {
	Get(key interface{}) (interface{}, error)
	GetWithTTL(key interface{}) (interface{}, time.Duration, error)
	Set(key interface{}, value interface{}, options ...Option) error
	Delete(key interface{}) error
	Invalidate(options ...InvalidateOption) error
	Clear() error
	GetType() string
}
