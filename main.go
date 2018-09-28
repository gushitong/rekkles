package main

import (
	"log"

	"github.com/tidwall/redcon"
)

var addr = ":6380"

func main() {

	server := NewAryadbServer()
	go log.Printf("started server at %s", addr)

	err := redcon.ListenAndServe(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			server.Handle(conn, cmd)
		},
		func(conn redcon.Conn) bool {
			// use this function to accept or deny the connection.
			// log.Printf("accept: %stor", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// this is called when the connection has been closed
			// log.Printf("closed: %stor, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
