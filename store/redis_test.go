package store

import (
	"context"
	"testing"
	"time"

	mocksStore "github.com/fmyxyz/gocache/test/mocks/store/clients"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	ctx := context.Background()
	options := &Options{
		Expiration: 6 * time.Second,
		Ctx:        ctx,
	}

	// When
	store := NewRedis(client, Expiration(6*time.Second), Ctx(ctx))

	// Then
	assert.IsType(t, new(RedisStore), store)
	assert.Equal(t, client, store.client)
	assert.Equal(t, options, store.options)
}

func TestRedisGet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	client.EXPECT().Get(context.Background(), "my-key").Return(&redis.StringCmd{})

	store := NewRedis(client)

	// When
	value, err := store.Get("my-key")

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, value)
}

func TestRedisSet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	client.EXPECT().Set(context.Background(), "my-key", cacheValue, 5*time.Second).Return(&redis.StatusCmd{})

	store := NewRedis(client, Expiration(6*time.Second))

	// When
	err := store.Set(cacheKey, cacheValue, Expiration(5*time.Second))

	// Then
	assert.Nil(t, err)
}

func TestRedisSetWhenNoOptionsGiven(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	client.EXPECT().Set(context.Background(), "my-key", cacheValue, 6*time.Second).Return(&redis.StatusCmd{})

	store := NewRedis(client, Expiration(6*time.Second))

	// When
	err := store.Set(cacheKey, cacheValue)

	// Then
	assert.Nil(t, err)
}

func TestRedisSetWithTags(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheKey := "my-key"
	cacheValue := "my-cache-value"

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	client.EXPECT().Set(context.Background(), cacheKey, cacheValue, time.Duration(0)).Return(&redis.StatusCmd{})
	client.EXPECT().SAdd(context.Background(), "gocache_tag_tag1", "my-key").Return(&redis.IntCmd{})
	client.EXPECT().Expire(context.Background(), "gocache_tag_tag1", 720*time.Hour).Return(&redis.BoolCmd{})

	store := NewRedis(client)

	// When
	err := store.Set(cacheKey, cacheValue, Tags("tag1"))

	// Then
	assert.Nil(t, err)
}

func TestRedisDelete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheKey := "my-key"

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	client.EXPECT().Del(context.Background(), "my-key").Return(&redis.IntCmd{})

	store := NewRedis(client)

	// When
	err := store.Delete(cacheKey)

	// Then
	assert.Nil(t, err)
}

func TestRedisInvalidate(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	options := InvalidateOptions{
		Tags: []string{"tag1"},
	}

	cacheKeys := &redis.StringSliceCmd{}

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	client.EXPECT().SMembers(context.Background(), "gocache_tag_tag1").Return(cacheKeys)
	client.EXPECT().Del(context.Background(), "gocache_tag_tag1").Return(&redis.IntCmd{})

	store := NewRedis(client)

	// When
	err := store.Invalidate(options)

	// Then
	assert.Nil(t, err)
}

func TestRedisClear(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocksStore.NewMockRedisClientInterface(ctrl)
	client.EXPECT().FlushAll(context.Background()).Return(&redis.StatusCmd{})

	store := NewRedis(client)

	// When
	err := store.Clear()

	// Then
	assert.Nil(t, err)
}

func TestRedisGetType(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocksStore.NewMockRedisClientInterface(ctrl)

	store := NewRedis(client)

	// When - Then
	assert.Equal(t, RedisType, store.GetType())
}
