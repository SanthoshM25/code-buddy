package server

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String("port", ":8080", "http service port")

func StartServer() {
	http.HandleFunc("/ws", handleConnection)

	svr := &http.Server{
		Addr: *port,
	}

	log.Fatal(svr.ListenAndServe())
}
