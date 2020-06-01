package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name    string
	Surname string
	Age     int
	Weight  float64
}

func TestNewRedisClient(t *testing.T) {
	client, err := NewRedisClient("localhost:6379", "", 0)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNew(t *testing.T) {
	client, err := NewRedisClient("localhost:6379", "", 0)
	cache := New(client, "account")
	assert.Nil(t, err)
	assert.NotNil(t, cache)
}

func TestCache(t *testing.T) {
	client, err := NewRedisClient("localhost:6379", "", 0)
	cache := New(client, "account")
	assert.Nil(t, err)

	u1 := &User{
		Name:    "crl",
		Surname: "r",
		Age:     25,
		Weight:  58.78,
	}
	u := &User{}

	users := []*User{u1, u1}
	us := []*User{}

	err = cache.Set("users:1", u1, 1000)
	assert.Nil(t, err)

	err = cache.Get("users:1", u)
	assert.Nil(t, err)
	assert.Equal(t, u.Name, "crl")
	assert.Equal(t, u.Weight, 58.78)

	err = cache.Set("users:hashlist", users, 10000)
	assert.Nil(t, err)

	err = cache.Get("users:hashlist", &us)
	assert.Nil(t, err)
	assert.Equal(t, us[0].Name, "crl")
	assert.Equal(t, us[0].Weight, 58.78)
}

func TestDelete(t *testing.T) {
	client, err := NewRedisClient("localhost:6379", "", 0)
	cache := New(client, "account")
	assert.Nil(t, err)

	u1 := &User{
		Name:    "crl",
		Surname: "r",
		Age:     25,
		Weight:  58.78,
	}
	err = cache.Set("users:1", u1, 1000)
	assert.Nil(t, err)

	err = cache.Delete("users:1")
	assert.Nil(t, err)

	err = cache.Get("users:1", u1)
	assert.EqualError(t, err, "redis.Get: redis: nil")
}

func TestPurge(t *testing.T) {
	client, err := NewRedisClient("localhost:6379", "", 0)
	cache := New(client, "purge")
	assert.Nil(t, err)

	u1 := &User{
		Name:    "crl",
		Surname: "r",
		Age:     25,
		Weight:  58.78,
	}
	err = cache.Set("users:1", u1, 1000)
	assert.Nil(t, err)
	err = cache.Set("users:2", u1, 1000)
	assert.Nil(t, err)
	err = cache.Set("users:3", u1, 1000)
	assert.Nil(t, err)

	total, err := cache.Purge()
	assert.Nil(t, err)
	assert.Equal(t, total, int64(3))

	err = cache.Get("users:1", u1)
	assert.EqualError(t, err, "redis.Get: redis: nil")
	err = cache.Get("users:2", u1)
	assert.EqualError(t, err, "redis.Get: redis: nil")
	err = cache.Get("users:3", u1)
	assert.EqualError(t, err, "redis.Get: redis: nil")
}
