package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"powerFactors/internal/server"
)

func main() {
	// We use as default the 0.0.0.0 address to bypass the network isolation of docker.
	ip := flag.String("address", "0.0.0.0", "The IP address that this server will use.")
	port := flag.String("port", "3000", "The port that this server will listen to.")

	flag.Parse()

	address := *ip + ":" + *port

	serverManager := server.NewServer()

	srv := &http.Server{
		Handler:      serverManager.Router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
