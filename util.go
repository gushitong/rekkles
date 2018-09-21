package main

import "strings"

func lower(s []byte) string {
	return strings.ToLower(string(s))
}

func RegisterHandler(db *aryadb, key string, handler Handler) {
	db.Handlers[key] = handler
}
