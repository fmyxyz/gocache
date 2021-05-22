package codec

import (
	"errors"
	"github.com/fmyxyz/gocache/test/mocks"
	"testing"
	"time"

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
	codec := New(store)

	// Then
	assert.IsType(t, new(Codec), codec)
}

func TestGetWhenHit(t *testing.T) {
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

	codec := New(store)

	// When
	value, err := codec.Get("my-key")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)

	assert.Equal(t, 1, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestGetWithTTLWhenHit(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &struct {
		Hello string
	}{
		Hello: "world",
	}

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().GetWithTTL("my-key").Return(cacheValue, 1*time.Second, nil)

	codec := New(store)

	// When
	value, ttl, err := codec.GetWithTTL("my-key")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
	assert.Equal(t, 1*time.Second, ttl)

	assert.Equal(t, 1, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestGetWithTTLWhenMiss(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unable to find in store")

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().GetWithTTL("my-key").Return(nil, 0*time.Second, expectedErr)

	codec := New(store)

	// When
	value, ttl, err := codec.GetWithTTL("my-key")

	// Then
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)
	assert.Equal(t, 0*time.Second, ttl)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 1, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestGetWhenMiss(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unable to find in store")

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Get("my-key").Return(nil, expectedErr)

	codec := New(store)

	// When
	value, err := codec.Get("my-key")

	// Then
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 1, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestSetWhenSuccess(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &struct {
		Hello string
	}{
		Hello: "world",
	}

	s := mocksStore.NewMockStoreInterface(ctrl)
	option := store.Expiration(5 * time.Second)
	s.EXPECT().Set("my-key", cacheValue, mocks.FuncEq(option)).Return(nil)

	codec := New(s)

	// When
	err := codec.Set("my-key", cacheValue, option)

	// Then
	assert.Nil(t, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 1, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestSetWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &struct {
		Hello string
	}{
		Hello: "world",
	}

	option := store.Expiration(5 * time.Second)

	expectedErr := errors.New("Unable to set value in store")

	s := mocksStore.NewMockStoreInterface(ctrl)
	s.EXPECT().Set("my-key", cacheValue, mocks.FuncEq(option)).Return(expectedErr)

	codec := New(s)

	// When
	err := codec.Set("my-key", cacheValue, option)

	// Then
	assert.Equal(t, expectedErr, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 1, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestDeleteWhenSuccess(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Delete("my-key").Return(nil)

	codec := New(store)

	// When
	err := codec.Delete("my-key")

	// Then
	assert.Nil(t, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 1, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TesDeleteWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unable to delete key")

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Delete("my-key").Return(expectedErr)

	codec := New(store)

	// When
	err := codec.Delete("my-key")

	// Then
	assert.Equal(t, expectedErr, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 1, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestInvalidateWhenSuccess(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocksStore.NewMockStoreInterface(ctrl)
	option := store.InvalidateTags("tag1")
	s.EXPECT().Invalidate(mocks.FuncEq(option)).Return(nil)

	codec := New(s)

	// When
	err := codec.Invalidate(option)

	// Then
	assert.Nil(t, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 1, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestInvalidateWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unexpected error when invalidating data")

	s := mocksStore.NewMockStoreInterface(ctrl)
	option := store.InvalidateTags("tag1")
	s.EXPECT().Invalidate(mocks.FuncEq(option)).Return(expectedErr)

	codec := New(s)

	// When
	err := codec.Invalidate(option)

	// Then
	assert.Equal(t, expectedErr, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 1, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestClearWhenSuccess(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Clear().Return(nil)

	codec := New(store)

	// When
	err := codec.Clear()

	// Then
	assert.Nil(t, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 1, codec.GetStats().ClearSuccess)
	assert.Equal(t, 0, codec.GetStats().ClearError)
}

func TestClearWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unexpected error when clearing cache")

	store := mocksStore.NewMockStoreInterface(ctrl)
	store.EXPECT().Clear().Return(expectedErr)

	codec := New(store)

	// When
	err := codec.Clear()

	// Then
	assert.Equal(t, expectedErr, err)

	assert.Equal(t, 0, codec.GetStats().Hits)
	assert.Equal(t, 0, codec.GetStats().Miss)
	assert.Equal(t, 0, codec.GetStats().SetSuccess)
	assert.Equal(t, 0, codec.GetStats().SetError)
	assert.Equal(t, 0, codec.GetStats().DeleteSuccess)
	assert.Equal(t, 0, codec.GetStats().DeleteError)
	assert.Equal(t, 0, codec.GetStats().InvalidateSuccess)
	assert.Equal(t, 0, codec.GetStats().InvalidateError)
	assert.Equal(t, 0, codec.GetStats().ClearSuccess)
	assert.Equal(t, 1, codec.GetStats().ClearError)
}

func TestGetStore(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)

	codec := New(store)

	// When - Then
	assert.Equal(t, store, codec.GetStore())
}

func TestGetStats(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocksStore.NewMockStoreInterface(ctrl)

	codec := New(store)

	// When - Then
	expectedStats := &Stats{}
	assert.Equal(t, expectedStats, codec.GetStats())
}
