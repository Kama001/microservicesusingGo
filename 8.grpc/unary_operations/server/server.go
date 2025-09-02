package main

import (
	"context"
	"fmt"
	proto "grpc/protoc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedHelloServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	proto.RegisterHelloServiceServer(srv, &server{})
	reflection.Register(srv)
	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) SayHello(c context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	fmt.Println("request received from client:", req.SomeString)
	return &proto.HelloResponse{Response: "hello from server"}, nil
}
