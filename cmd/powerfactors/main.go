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

	//// Create a channel to listen for an interrupt signal
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	//
	//// Start the server in a goroutine
	//go func() {
	//	log.Printf("Server listening on %s\n", address)
	//	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("Server error: %v\n", err)
	//	}
	//}()
	//
	//// Wait for an interrupt signal
	//<-interrupt
	//log.Println("Shutting down server...")
	//
	//// Create a context with a timeout for graceful shutdown
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//
	//// Shutdown the server gracefully
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatalf("Server shutdown error: %v\n", err)
	//}
	//
	//log.Println("Server gracefully stopped")

}
