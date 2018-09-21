package main

import (
	"testing"
	"github.com/go-redis/redis"
	"github.com/tidwall/redcon"
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

func TestHandler(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
	})

	pong, err := client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
	if pong != "PONG" {
		t.Fatal("Expect PONG, got " + pong)
	}
}
