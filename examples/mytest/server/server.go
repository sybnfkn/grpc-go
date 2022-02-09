package main

import (
	"context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/mytest/proto"
	"log"
	"net"
)

type SearchService struct {
	// 不明白这行什么意思,不写这个，相当于下面Search实现不了接口，会出现报错
	pb.UnsafeSearchServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server "}, nil
}

const PORT = "9001"

func main() {
	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err:%v", err)
	}

	server.Serve(lis)
}
