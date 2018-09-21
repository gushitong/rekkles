package main

import "strings"



func lower(s []byte) string {
	return strings.ToLower(string(s))
}

func RegisterHandler(db *db, key string, handler Handler) {
	db.Handlers[key] = handler
}
