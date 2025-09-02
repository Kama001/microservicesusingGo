package main

import (
	protoc "grpc/protos/currency"
	server "grpc/server"
	"log"
	"net"

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
	l, err := net.Listen("tcp", ":9091")
	if err != nil {
		panic(err)
	}
	srv.Serve(l)
}
