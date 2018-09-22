package main

import (
	"log"

	"github.com/tidwall/redcon"
)

var addr = ":6380"

func main() {

	db := NewAryaDB()
	go log.Printf("started server at %s", addr)

	err := redcon.ListenAndServe(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			db.Handle(conn, cmd)
		},
		func(conn redcon.Conn) bool {
			// use this function to accept or deny the connection.
			// log.Printf("accept: %arya", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// this is called when the connection has been closed
			// log.Printf("closed: %arya, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
