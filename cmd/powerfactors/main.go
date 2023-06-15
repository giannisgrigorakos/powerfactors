package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"powerFactors/internal/server"
)

func main() {

	// TODO:
	// check edge cases sta test
	// docker + docker-compose
	// extensible (some comments maybe on Server and new supported periods)
	// README
	// makefile
	// Documentation
	// check that the passed flags work
	// delete utils
	// refactor to be more concise

	// I might want to have 0.0.0.0 because of docker
	ip := flag.String("address", "127.0.0.1", "The IP address that this server will use.")
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
