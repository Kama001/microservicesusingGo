package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testhandlers/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	hg := handlers.NewGoodbye(l)

	mu := http.NewServeMux()

	mu.HandleFunc("/", hh.HelloHandler)
	mu.HandleFunc("/goodbye", hg.GoodbyeHandler)
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
