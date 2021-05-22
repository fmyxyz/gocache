package store

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// RistrettoType represents the storage type as a string value
	RistrettoType = "ristretto"
	// RistrettoTagPattern represents the tag pattern to be used as a key in specified storage
	RistrettoTagPattern = "gocache_tag_%s"
)

// RistrettoClientInterface represents a dgraph-io/ristretto client
type RistrettoClientInterface interface {
	Get(key interface{}) (interface{}, bool)
	SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool
	Del(key interface{})
	Clear()
}

// RistrettoStore is a store for Ristretto (memory) library
type RistrettoStore struct {
	client  RistrettoClientInterface
	options *Options
}

// NewRistretto creates a new store to Ristretto (memory) library instance
func NewRistretto(client RistrettoClientInterface, options ...Option) *RistrettoStore {
	opts := &Options{}
	for _, option := range options {
		option(opts)
	}

	return &RistrettoStore{
		client:  client,
		options: opts,
	}
}

// Get returns data stored from a given key
func (s *RistrettoStore) Get(key interface{}) (interface{}, error) {
	var err error

	value, exists := s.client.Get(key)
	if !exists {
		err = errors.New("Value not found in Ristretto store")
	}

	return value, err
}

// GetWithTTL returns data stored from a given key and its corresponding TTL
func (s *RistrettoStore) GetWithTTL(key interface{}) (interface{}, time.Duration, error) {
	value, err := s.Get(key)
	return value, 0, err
}

// Set defines data in Ristretto memoey cache for given key identifier
func (s *RistrettoStore) Set(key interface{}, value interface{}, opts ...Option) error {
	var err error

	options := &Options{}
	if len(opts) == 0 {
		options = s.options
	}
	for _, opt := range opts {
		opt(options)
	}

	if set := s.client.SetWithTTL(key, value, options.CostValue(), options.ExpirationValue()); !set {
		err = fmt.Errorf("An error has occurred while setting value '%v' on key '%v'", value, key)
	}

	if err != nil {
		return err
	}

	if tags := options.TagsValue(); len(tags) > 0 {
		s.setTags(key, tags)
	}

	return nil
}

func (s *RistrettoStore) setTags(key interface{}, tags []string) {
	for _, tag := range tags {
		var tagKey = fmt.Sprintf(RistrettoTagPattern, tag)
		var cacheKeys = []string{}

		if result, err := s.Get(tagKey); err == nil {
			if bytes, ok := result.([]byte); ok {
				cacheKeys = strings.Split(string(bytes), ",")
			}
		}

		var alreadyInserted = false
		for _, cacheKey := range cacheKeys {
			if cacheKey == key.(string) {
				alreadyInserted = true
				break
			}
		}

		if !alreadyInserted {
			cacheKeys = append(cacheKeys, key.(string))
		}

		s.Set(tagKey, []byte(strings.Join(cacheKeys, ",")), Expiration(720*time.Hour))
	}
}

// Delete removes data in Ristretto memoey cache for given key identifier
func (s *RistrettoStore) Delete(key interface{}) error {
	s.client.Del(key)
	return nil
}

// Invalidate invalidates some cache data in Redis for given Options
func (s *RistrettoStore) Invalidate(options InvalidateOptions) error {
	if tags := options.TagsValue(); len(tags) > 0 {
		for _, tag := range tags {
			var tagKey = fmt.Sprintf(RistrettoTagPattern, tag)
			result, err := s.Get(tagKey)
			if err != nil {
				return nil
			}

			var cacheKeys = []string{}
			if bytes, ok := result.([]byte); ok {
				cacheKeys = strings.Split(string(bytes), ",")
			}

			for _, cacheKey := range cacheKeys {
				s.Delete(cacheKey)
			}
		}
	}

	return nil
}

// GetType returns the store type
func (s *RistrettoStore) GetType() string {
	return RistrettoType
}

// Clear resets all data in the store
func (s *RistrettoStore) Clear() error {
	s.client.Clear()
	return nil
}
