package cache

import (
	"errors"
	"fmt"
	"github.com/fmyxyz/gocache/test/mocks"
	"testing"
	"time"

	"github.com/fmyxyz/gocache/codec"
	"github.com/fmyxyz/gocache/store"
	mocksStore "github.com/fmyxyz/gocache/test/mocks/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)

	// When
	cache := New(store)

	// Then
	assert.IsType(t, new(Cache), cache)
	assert.IsType(t, new(codec.Codec), cache.codec)

	assert.Equal(t, store, cache.codec.GetStore())
}

func TestCacheSet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	value := &struct {
		Hello string
	}{
		Hello: "world",
	}
	s := mocksStore.NewMockStoreInterface(ctrl)
	option := store.Expiration(5 * time.Second)
	s.EXPECT().Set("my-key", value, mocks.FuncEq(option)).Return(nil)

	cache := New(s)

	// When
	err := cache.Set("my-key", value, option)
	assert.Nil(t, err)
}

func TestCacheSetWhenErrorOccurs(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	value := &struct {
		Hello string
	}{
		Hello: "world",
	}

	storeErr := errors.New("An error has occurred while inserting data into store")

	s := mocksStore.NewMockStoreInterface(ctrl)
	option := store.Expiration(5 * time.Second)
	fmt.Printf("%p\n", option)
	s.EXPECT().Set("my-key", value, mocks.FuncEq(option)).Return(storeErr)
	fmt.Println("-----00----")
	cache := New(s)
	_ = cache
	// When
	err := cache.Set("my-key", value, option)
	fmt.Printf("\n------%p, %p\n-------\n", storeErr, err)
}

func TestCacheGet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &struct {
		Hello string
	}{
		Hello: "world",
	}

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Get("my-key").Return(cacheValue)

	cache := New(store)

	// When
	value, err := cache.Get("my-key")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
}

func TestCacheGetWhenNotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	returnedErr := errors.New("Unable to find item in store")

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Get("my-key").Return(nil, returnedErr)

	cache := New(store)

	// When
	value, err := cache.Get("my-key")

	// Then
	assert.Nil(t, value)
	assert.Equal(t, returnedErr, err)
}

func TestCacheGetWithTTL(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &struct {
		Hello string
	}{
		Hello: "world",
	}
	expiration := 1 * time.Second

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().GetWithTTL("my-key").
		Return(cacheValue, expiration, nil)

	cache := New(store)

	// When
	value, ttl, err := cache.GetWithTTL("my-key")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
	assert.Equal(t, expiration, ttl)
}

func TestCacheGetWithTTLWhenNotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	returnedErr := errors.New("Unable to find item in store")
	expiration := 0 * time.Second

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().GetWithTTL("my-key").
		Return(nil, expiration, returnedErr)

	cache := New(store)

	// When
	value, ttl, err := cache.GetWithTTL("my-key")

	// Then
	assert.Nil(t, value)
	assert.Equal(t, returnedErr, err)
	assert.Equal(t, expiration, ttl)
}

func TestCacheGetCodec(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)

	cache := New(store)

	// When
	value := cache.GetCodec()

	// Then
	assert.IsType(t, new(codec.Codec), value)
	assert.Equal(t, store, value.GetStore())
}

func TestCacheGetType(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)

	cache := New(store)

	// When - Then
	assert.Equal(t, CacheType, cache.GetType())
}

func TestCacheGetCacheKeyWhenKeyIsString(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)

	cache := New(store)

	// When
	computedKey := cache.getCacheKey("my-Key")

	// Then
	assert.Equal(t, "my-Key", computedKey)
}

func TestCacheGetCacheKeyWhenKeyIsStruct(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)

	cache := New(store)

	// When
	key := &struct {
		Hello string
	}{
		Hello: "world",
	}

	computedKey := cache.getCacheKey(key)

	// Then
	assert.Equal(t, "8144fe5310cf0e62ac83fd79c113aad2", computedKey)
}

func TestCacheDelete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Delete("my-key").Return(nil)

	cache := New(store)

	// When
	err := cache.Delete("my-key")

	// Then
	assert.Nil(t, err)
}

func TestCacheInvalidate(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocksStore.NewMockStoreInterface(ctrl)
	option := store.InvalidateTags("tag1")
	s.EXPECT().Invalidate(mocks.FuncEq(option)).Return(nil)

	cache := New(s)

	// When
	err := cache.Invalidate(option)

	// Then
	assert.Nil(t, err)
}

func TestCacheInvalidateWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unexpected error during invalidation")

	s := mocksStore.NewMockStoreInterface(ctrl)
	option := store.InvalidateTags("tag1")
	s.EXPECT().Invalidate(mocks.FuncEq(option)).Return(expectedErr)

	cache := New(s)

	// When
	err := cache.Invalidate(option)

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestCacheClear(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Clear().Return(nil)

	cache := New(store)

	// When
	err := cache.Clear()

	// Then
	assert.Nil(t, err)
}

func TestCacheClearWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unexpected error during invalidation")

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Clear().Return(expectedErr)

	cache := New(store)

	// When
	err := cache.Clear()

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestCacheDeleteWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unable to delete key")

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Delete("my-key").Return(expectedErr)

	cache := New(store)

	// When
	err := cache.Delete("my-key")

	// Then
	assert.Equal(t, expectedErr, err)
}
