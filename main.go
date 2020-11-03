package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/vaibhav/CoffeeOnDemand/handlers"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	// handlers
	productHandler := handlers.NewProducts(l)

	// gorilla Router, generally named as router
	serveMux := mux.NewRouter()

	// sub-router
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productHandler.GetProducts)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productHandler.AddProduct)
	postRouter.Use(productHandler.MiddlewareValidateProduct) // will be executed prior to the Handler

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProducts)
	putRouter.Use(productHandler.MiddlewareValidateProduct)

	// serveMux.Handle("/products", productHandler)

	// create a new server
	server := http.Server{
		Addr:         ":8000",           // configure this bindAddress
		Handler:      serveMux,          // set the default handler
		ErrorLog:     l,                 // set the Logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from client
		WriteTimeout: 10 * time.Second,  // max time to write request to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP keep-alive
	}

	// start the server
	go func() {
		l.Println("starting server on port:8000")

		err := server.ListenAndServe()
		if err != nil {
			// l.Fatalf("Failed to start server, %s\n", err)
			l.Printf("Failed to start server, %s\n", err)
			os.Exit(1)
		}
	}()

	// trapping interupts and kill using Notify() and graceful shutdown of the server
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	l.Println("Recieved termination signal with graceful shutdown", sig)

	shutdownTimeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(shutdownTimeoutContext)
}
