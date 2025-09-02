package main

import (
	"context"
	"fmt"
	proto "grpc/protoc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.TODO(), &proto.HelloRequest{SomeString: "hello from client"})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Response)
}
