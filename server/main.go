package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/santhoshm25/code-buddy/server/internal"
)

var port = flag.String("port", ":8080", "http service port")

func main() {
	http.HandleFunc("/ws", handler.handleConnection)

	svr := &http.Server{
		Addr: *port,
	}

	log.Fatal(svr.ListenAndServe())
}
