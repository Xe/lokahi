package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	// Add HTTP pprof routes
	_ "net/http/pprof"

	// Add tracing routes
	_ "golang.org/x/net/trace"

	// Expvar routes
	_ "expvar"
)

func init() {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Printf("manhole: cannot bind to 127.0.0.1:0: %v", err)
		return
	}

	log.Printf("manhole: Now listening on http://%s", l.Addr())

	rpc.HandleHTTP()
	go http.Serve(l, nil)
}
