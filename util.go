package main

import "strings"

func LowerString(s []byte) string {
	return strings.ToLower(string(s))
}

func RegisterHandler(db *server, key string, handler Handler) {
	db.Handlers[key] = handler
}
