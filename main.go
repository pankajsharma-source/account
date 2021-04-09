package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pankajsharma-source/env"
	"github.com/pankajsharma-source/account/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
	// hello this is nikitas commenttttt! dont ever erasse it pls. my favorite animal is the koala or sea otter or dolphin or horse or puppy or dog or cat or kitten or turtle or parrot or gold-fish
	env.Parse()

	l := log.New(os.Stdout, "user-profile-api ", log.LstdFlags)

	// create the handlers
	uh := handlers.NewUser(l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	//	getRouter.HandleFunc("/{id:[0-9]+}", uh.GetUser)
	getRouter.HandleFunc("/", uh.GetUser)


	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", uh.AddUser)
	postRouter.Use(uh.MiddlewareValidateUser)

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}