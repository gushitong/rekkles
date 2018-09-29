package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	opt := new(Options)
	flag.StringVar(&opt.Addr, "b", ":6380", "listen address")
	flag.StringVar(&opt.Dir, "d", "/tmp/aryadb", "working dir")
	flag.StringVar(&opt.ValueDir, "v", "/tmp/aryadb", "value log dir")
	flag.StringVar(&opt.Auth, "a", "", "auth string")
	flag.BoolVar(&opt.SyncWrites, "s", true, "sync all writes to disk. Setting this to false would achieve better performance, but may cause data to be lost.")
	flag.Parse()

	server, err := NewAryadbServer(opt)
	if err != nil {
		fmt.Printf("Server not started: %s", err)
		return
	}
	go log.Printf("started server at %s", opt.Addr)
	log.Fatal(server.ListenAndSrv())
}
