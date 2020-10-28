package xredis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert := require.New(t)
	c := Config{
		Addr:     ":6379",
		Password: "",
		DB:       0,
	}
	redis := New(c)
	s := redis.Set("solar", "solar", time.Second*30)
	assert.NoError(s.Err())
	r, err := redis.Get("solar").Result()
	assert.NoError(err)
	assert.Equal("solar", r)
}
