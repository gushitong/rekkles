package main

import (
	"testing"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestZsetZadd(t *testing.T) {
	key := "zadd"
	k1, k2, k3 := "k1", "k2", "k3"
	s1, s2, s3 := 1.0, 2.0, 3.0
	client.Del(key)

	n ,err := client.ZAdd(key, redis.Z{s1, k1}).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), n)
	client.ZAdd(key, redis.Z{s2, k2})
	client.ZAdd(key, redis.Z{s3, k3})

	n, err = client.ZCard(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), n)
}

func TestZsetZcard(t *testing.T) {
	key := "zcard"
	k1, k2, k3 := "k1", "k2", "k3"
	s1, s2, s3 := 1.0, 2.0, 3.0
	client.Del(key)

	n, err := client.ZCard(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), n)

	client.ZAdd(key, redis.Z{s1, k1})
	client.ZAdd(key, redis.Z{s2, k2})
	client.ZAdd(key, redis.Z{s3, k3})

	n, err = client.ZCard(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), n)
}

func TestZsetZcount(t *testing.T) {
	key := "zcount"
	k1, k2, k3 := "k1", "k2", "k3"
	s1, s2, s3 := 1.0, 2.0, 3.0
	client.Del(key)

	n, err := client.ZCount(key, "0", "3").Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), n)

    client.ZAdd(key, redis.Z{s1, k1})
	client.ZAdd(key, redis.Z{s2, k2})
	client.ZAdd(key, redis.Z{s3, k3})

	n, err = client.ZCount(key, "0", "3").Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), n)
}

func TestZsetZincrby(t *testing.T) {
	key := "zincrby"
	k := "k"
	client.Del(key)
	n, err := client.ZIncrBy(key, 1.0, k).Result()
	assert.Nil(t, err)
	assert.Equal(t, float64(1), n)
	n, err = client.ZIncrBy(key, 2.0, k).Result()
	assert.Nil(t, err)
	assert.Equal(t, float64(3), n)
}

func TestZRange(t *testing.T) {
	key := "zrange"
	k1, k2, k3 := "k1", "k2", "k3"
	s1, s2, s3 := 1.0, 2.0, 3.0
	client.Del(key)
	client.ZAdd(key, redis.Z{s1, k1})
	client.ZAdd(key, redis.Z{s2, k2})
	client.ZAdd(key, redis.Z{s3, k3})
	val, err := client.ZRange(key, 0, -1).Result()
	assert.Nil(t, err)
	assert.Len(t, val, 3)

	val, err = client.ZRange(key, 0, 1).Result()
	assert.Nil(t, err)
	assert.Len(t, val, 2)
}

func TestZRangebyscore(t *testing.T) {
	key := "zrangebyscore"
	k1, k2, k3 := "k1", "k2", "k3"
	s1, s2, s3 := 2.0, 5.0, 7.0
	client.Del(key)
	client.ZAdd(key, redis.Z{s1, k1})
	client.ZAdd(key, redis.Z{s2, k2})
	client.ZAdd(key, redis.Z{s3, k3})
	val, err := client.ZRangeByScore(key, redis.ZRangeBy{Min:"7", Max:"15"}).Result()
	assert.Nil(t, err)
	assert.Len(t, val, 1)
	val, err = client.ZRangeByScore(key, redis.ZRangeBy{Min:"0", Max:"15"}).Result()
	assert.Nil(t, err)
	assert.Len(t, val, 3)
}


func TestZRank(t *testing.T) {
	key := "zrank"
	k1, k2, k3 := "k1", "k2", "k3"
	s1, s2, s3 := 2.0, 5.0, 7.0
	client.Del(key)
	client.ZAdd(key, redis.Z{s1, k1})
	client.ZAdd(key, redis.Z{s2, k2})
	client.ZAdd(key, redis.Z{s3, k3})
	val, err := client.ZRank(key, k1).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), val)
	val, err = client.ZRank(key, k3).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), val)
}

func TestZScore(t *testing.T) {
	key := "zscore"
	k1, k2, k3 := "k1", "k2", "k3"
	s1, s2, s3 := 2.0, 5.0, 7.0
	client.Del(key)
	client.ZAdd(key, redis.Z{s1, k1})
	client.ZAdd(key, redis.Z{s2, k2})
	client.ZAdd(key, redis.Z{s3, k3})
	val, err := client.ZScore(key, k1).Result()
	assert.Nil(t, err)
	assert.Equal(t, 2.0, val)
	val, err = client.ZScore(key, k3).Result()
	assert.Nil(t, err)
	assert.Equal(t, 7.0, val)
}

