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
	"google.golang.org/grpc"

	protos "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/protos/currency"
)

var bindAddress = ":9090"

func main() {

	conn, err := grpc.NewClient("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// cc = currency client
	cc := protos.NewCurrencyClient(conn)

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hp := handlers.NewProducts(l, cc)

	mu := mux.NewRouter()

	getRouter := mu.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", hp.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", hp.ListSingle)

	putRouter := mu.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", hp.Update)
	putRouter.Use(hp.MiddlewareProductValidation)

	postRouter := mu.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", hp.Create)
	postRouter.Use(hp.MiddlewareProductValidation)

	deleteRouter := mu.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", hp.Delete)

	s := http.Server{
		Addr:         bindAddress,
		Handler:      mu,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		l.Printf("Starting server on %s\n", bindAddress)
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
