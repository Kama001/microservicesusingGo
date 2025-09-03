package main

import (
	"net"

	"github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/data"
	protos "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/protos/currency"
	"k8s.io/klog/v2"

	server "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	logger := klog.NewKlogr()
	rates := data.NewExchangeRates(&logger)

	// create a new gRPC server, use WithInsecure to allow http connections
	srv := grpc.NewServer()

	// create an instance of the Currency server
	cs := server.NewCurrency(&logger, rates)

	// register the currency server
	protos.RegisterCurrencyServer(srv, cs)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(srv)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		panic(err)
	}

	// listen for requests
	srv.Serve(l)
}
