package main

import (
	"context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/mytest/proto"
)

type SearchService struct {
	// 不明白这行什么意思
	pb.UnsafeSearchServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server "}, nil
}

const PORT = "9001"

func main() {
	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})
}
