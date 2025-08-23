package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"restapi/handlers"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hp := handlers.NewProducts(l)

	mu := mux.NewRouter()

	getRouter := mu.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", hp.GetProducts)

	putRouter := mu.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", hp.UpdateProducts)

	postRouter := mu.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", hp.AddProducts)

	s := http.Server{
		Addr:         ":9090",
		Handler:      mu,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	<-sigChan
	l.Println("Received shutdown signal, shutting down server...")
	tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	s.Shutdown(tc)
}
