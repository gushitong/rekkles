package main

import (
	"github.com/go-redis/redis"
	"github.com/tidwall/redcon"
	"testing"
)

func init() {
	db := NewAryaDB()

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
}

func TestCommand(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
	})

	key := ":arya_key:"
	val := ":arya_val:"

	pong, err := client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
	if pong != "PONG" {
		t.Fatal("Expect PONG, got " + pong)
	}

	ok, err := client.Set(key, val, 0).Result()
	if err != nil {
		t.Fatal(err)
	}
	if ok != "OK" {
		t.Fatal("Expect OK, got " + ok)
	}

	v, err := client.Get(key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if v != val {
		t.Fatal("Expect '" + val + "' got " + v)
	}

	n, err := client.Append(key, v).Result()
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(len(val)*2) {
		t.Fatal("Value length mismatch")
	}
}
