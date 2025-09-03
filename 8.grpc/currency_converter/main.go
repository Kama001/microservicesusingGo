package main

import (
	"log"
	"net"

	protoc "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/protos/currency"

	server "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	srv := grpc.NewServer()
	logger := log.New(log.Writer(), "currency ", log.LstdFlags)
	cs := server.NewCurrency(logger)
	protoc.RegisterCurrencyServer(srv, cs)
	// this will display the api's getting served
	reflection.Register(srv)
	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		panic(err)
	}
	srv.Serve(l)
}
