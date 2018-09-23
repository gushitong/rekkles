package main

import (
	"github.com/go-redis/redis"
	"github.com/tidwall/redcon"
	"testing"
)

var client *redis.Client

func init() {
	db := NewAryadbServer()
	go redcon.ListenAndServe(":6380",
		func(conn redcon.Conn, cmd redcon.Command) {
			db.Handle(conn, cmd)
		},
		func(conn redcon.Conn) bool {
			return true
		},
		func(conn redcon.Conn, err error) {

		},
	)
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
	})
}

func TestCmdAppend(t *testing.T) {
	key := "k_append"
	val := "val"

	client.Set(key, "", 0)

	n, err := client.Append(key, val).Result()
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(len(val)) {
		t.Fatal("Value length mismatch")
	}

	n, err = client.Append(key, val).Result()
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(len(val)*2) {
		t.Fatal("Value length mismatch")
	}
}

func TestCmdBitcount(t *testing.T) {
	key := "k_bitcount"
	val := "gushitong"

	client.Set(key, val, 0)
	n, err := client.BitCount(key, nil).Result()
	if err != nil {
		t.Fatal(err)
	}
	if n != 42 {
		t.Fatal("Bitcount result mismatch")
	}
}

func TestCmdDecr(t *testing.T) {
	key := "k_decr"
	val := 10

	client.Set(key, val, 0)
	v, err := client.Decr(key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if v != 9 {
		println(v)
		t.Fatal("Decr result mismatch.")
	}
}

func TestCmdDecrBy(t *testing.T)  {
	key := "k_decrby"
	val := 10

	client.Set(key, val, 0)
	v, err := client.DecrBy(key, 5).Result()
	if err != nil {
		t.Fatal(err)
	}
	if v != 5 {
		t.Fatal("Decrby result mismatch.")
	}
}

func TestCmdGetbit(t *testing.T) {
	key := "k_getbit"
	val := "g"

	client.Set(key, val, 0)
	v, err := client.GetBit(key, 2).Result()
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Fatal("Getbit result mismatch.")
	}

	v, err = client.GetBit(key, 3).Result()
	if err != nil {
		t.Fatal(err)
	}
	if v != 0 {
		t.Fatal("Getbit result mismatch.")
	}
}
