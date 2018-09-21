package main

import (
	"github.com/tidwall/redcon"
	"fmt"
)

/**
 Redis: PING
 */
func ping(s Storage, conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

/**
 Redis: GET
 */
func get(s Storage, conn redcon.Conn, cmd redcon.Command) {

	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}

	val, err := s.Get(cmd.Args[1])
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %storage", fmt.Sprint(err)))
	}else {
		conn.WriteBulk(val)
	}
}

/**
 Redis: SET
 */
func set(s Storage, conn redcon.Conn, cmd redcon.Command) {
	// todo: deal with set with args
	if len(cmd.Args) != 3 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}

	err := s.Set(cmd.Args[1], cmd.Args[2])
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR %storage", fmt.Sprint(err)))
	}else {
		conn.WriteString("OK")
	}
}