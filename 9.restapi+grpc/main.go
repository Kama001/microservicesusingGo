package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"restapi/data"
	"restapi/handlers"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"

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

	l := klog.NewKlogr()

	// create database instance
	db := data.NewProductsDB(cc, &l)

	hp := handlers.NewProducts(&l, db)

	mu := mux.NewRouter()

	getRouter := mu.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", hp.ListAll)
	getRouter.HandleFunc("/products", hp.ListAll).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", hp.ListSingle)
	getRouter.HandleFunc("/products/{id:[0-9]+}", hp.ListSingle).Queries("currency", "{[A-Z]{3}}")

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
		l.Info("Starting server on ", bindAddress)
		err := s.ListenAndServe()
		if err != nil {
			l.Error(err, "Cannot listen on port %s", bindAddress)
		}
	}()
	<-sigChan
	l.Info("Received shutdown signal, shutting down server...")
	tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	s.Shutdown(tc)
}
