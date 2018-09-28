package main

import (
	"testing"
)

func TestHashDel(t *testing.T) {

	key := "hash_del"
	hashKey := "key"
	hashVal := "val"

	client.HDel(key, hashKey)

	client.HSet(key, hashKey, hashVal).Result()
	val, err := client.HDel(key, hashKey).Result()
	if err != nil {
		t.Fatal(err)
	}
	if val != 1 {
		t.Fatal("HashDel result mismatch")
	}
}

func TestHashExists(t *testing.T) {
	key := "hash_exists"
	hashKey := "key"
	hashVal := "val"

	client.HDel(key, hashKey)

	ok, err := client.HExists(key, hashKey).Result()
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("HashExist result mismatch")
	}
	client.HSet(key, hashKey, hashVal)
	ok, err = client.HExists(key, hashKey).Result()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("HashExist result mismatch")
	}
}

func TestHashGet(t *testing.T) {
	key := "hash_get"
	hashKey := "key"
	hashVal := "val"

	client.HDel(key, hashKey)

	val, err := client.HGet(key, hashKey).Result()
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Fatal("HashGet result mismatch")
	}
	client.HSet(key, hashKey, hashVal)
	if err != nil {
		t.Fatal(err)
	}
	if val != val {
		t.Fatal("HashGet result mismatch")
	}
}

func TestHashGetall(t *testing.T) {
	key := "hash_getall"
	hashKey := "key"
	hashVal := "val"

	client.HDel(key, hashKey).Result()

	val, err := client.HGetAll(key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(val) != 0 {
		t.Fatal("HashGetall result mismatch")
	}
	client.HSet(key, hashKey, hashVal)

	val, err = client.HGetAll(key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(val) != 1 || val[hashKey] != hashVal {
		t.Fatal("HashGetall result mismatch")
	}
}

func TestHashIncrby(t *testing.T) {
	key := "hash_incrby"
	hashKey := "key"
	hashVal := "10"

	client.HSet(key, hashKey, hashVal)
	i, err := client.HIncrBy(key, hashKey, 10).Result()
	if err != nil {
		t.Fatal(err)
	}
	if i != 20 {
		t.Fatal("HashIncrby result mismatch.")
	}
}

func TestHashIncrbyFloat(t *testing.T) {
	key := "hash_incrbyfloat"
	hashKey := "key"
	hashVal := "10.5"

	client.HSet(key, hashKey, hashVal)
	i, err := client.HIncrByFloat(key, hashKey, 10.6).Result()
	if err != nil {
		t.Fatal(err)
	}
	if i != 21.1 {
		t.Fatal("HashIncrbyFloat result mismatch.")
	}
}

func TestHashKeys(t *testing.T) {
	key := "hash_keys"
	hashKey1 := "key1"
	hashKey2 := "key2"
	hashVal := "val"

	client.HSet(key, hashKey1, hashVal)
	client.HSet(key, hashKey2, hashVal)
	keys, err := client.HKeys(key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || keys[0] != hashKey1 {
		t.Fatal("HashIncrby result mismatch.")
	}
}

func TestHashLen(t *testing.T) {
	key := "hash_len"
	hashKey := "key"
	hashVal := "val"

	client.HSet(key, hashKey, hashVal)
	i, err := client.HLen(key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatal("HashIncrby result mismatch.")
	}
}

func TestHashMget(t *testing.T) {
	key := "hash_mget"
	hashKey1 := "key1"
	hashKey2 := "key2"
	hashVal1 := "val1"
	hashVal2 := "val2"

	client.HSet(key, hashKey1, hashVal1)
	client.HSet(key, hashKey2, hashVal2)
	vals, err := client.HMGet(key, hashKey1, hashKey2).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 || vals[0] != hashVal1 {
		t.Fatal("HashIncrby result mismatch.")
	}
}

func TestHashMset(t *testing.T) {
	key := "hash_mset"
	hashMap := map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}

	client.Del(key)
	ok, err := client.HMSet(key, hashMap).Result()
	if err != nil {
		t.Fatal(err)
	}
	if ok != "OK" {
		t.Fatal("HashIncrby result mismatch.")
	}

	vals, err := client.HMGet(key, "key1", "key2").Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 || vals[0] != "val1" {
		t.Fatal("HashIncrby result mismatch.")
	}
}

func TestHashSetnx(t *testing.T) {
	key := "hash_setnx"
	hashKey := "key"
	hashVal1 := "val1"
	hashVal2 := "val2"

	client.HDel(key, hashKey)
	ok, err := client.HSetNX(key, hashKey, hashVal1).Result()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("HashSetnx result mismatch.")
	}
	ok, err = client.HSetNX(key, hashKey, hashVal2).Result()
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("HashSetnx result mismatch.")
	}
}

func TestHashStrlen(t *testing.T) {
	key := "hash_strlen"
	hashKey := "key1"
	hashVal := "val1"

	client.Del(key)
	client.HSet(key, hashKey, hashVal)
	//No HStrlen Method
}

func TestHashVals(t *testing.T) {
	key := "hash_vals"
	hashKey1 := "key1"
	hashKey2 := "key2"
	hashVal1 := "val1"
	hashVal2 := "val2"

	client.HSet(key, hashKey1, hashVal1)
	client.HSet(key, hashKey2, hashVal2)
	vals, err := client.HVals(key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 || vals[0] != hashVal1 {
		t.Fatal("HashIncrby result mismatch.")
	}
}
