package main

import (
	"context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/mytest/proto"
	"log"
)

const PORT = "9001"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.dial err:%v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "HELLO",
	})
	if err != nil {
		log.Fatalf("client.Search err:%v", err)
	}
	log.Printf("resp:%s", resp.GetResponse())
}
