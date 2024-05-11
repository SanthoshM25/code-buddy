package server

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String("port", ":8080", "http service port")

func Serve() {
	flag.Parse()

	http.HandleFunc("/ws/user/{userName}", handleConnection)
	http.HandleFunc("/ws/user/{userName}/hub/{hubId}", handleConnection)

	svr := &http.Server{
		Addr: *port,
	}
	log.Println("Server listening on port", *port)
	log.Fatal(svr.ListenAndServe())
}
