package marshaler

import (
	"github.com/fmyxyz/gocache/cache"
	"github.com/fmyxyz/gocache/store"
	"github.com/vmihailenco/msgpack"
)

// Marshaler is the struct that marshal and unmarshal cache values
type Marshaler struct {
	cache.CacheInterface
	returnObj interface{}
}

// New creates a new marshaler that marshals/unmarshals cache values
func New(cache cache.CacheInterface) *Marshaler {
	return &Marshaler{
		CacheInterface: cache,
		//returnObj: map[string]interface{}{},
	}
}

// New creates a new marshaler that marshals/unmarshals cache values
func (c *Marshaler) ReturnObj(returnObj interface{}) *Marshaler {
	return &Marshaler{
		CacheInterface: c.CacheInterface,
		returnObj:      returnObj,
	}
}

// Get obtains a value from cache and unmarshal value with given object
func (c *Marshaler) Get(key interface{}) (interface{}, error) {
	result, err := c.CacheInterface.Get(key)
	if err != nil {
		return nil, err
	}

	switch v := result.(type) {
	case []byte:
		err = msgpack.Unmarshal(v, c.returnObj)
	case string:
		err = msgpack.Unmarshal([]byte(v), c.returnObj)
	}

	if err != nil {
		return nil, err
	}

	return c.returnObj, nil
}

// Set sets a value in cache by marshaling value
func (c *Marshaler) Set(key, object interface{}, options ...store.Option) error {
	bytes, err := msgpack.Marshal(object)
	if err != nil {
		return err
	}

	return c.CacheInterface.Set(key, bytes, options...)
}
