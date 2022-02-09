package main

import (
	"context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/stream/proto"
	"io"
	"log"
	"time"
)

const PORT = "9002"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err:%v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	err = printLists(client, &pb.StreamRequest{
		Pt: &pb.StreamPoint{Name: "grpc stream client:list", Value: 2018},
	})
	if err != nil {
		log.Fatalf("printLists.err:%v", err)
	}

	//err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "grpc stream client:record", Value: 2018}})
	//if err != nil {
	//	log.Fatalf("printREcord.err:%v", err)
	//}

	//err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "grpc stream client:Route", Value: 2018}})
	//if err != nil {
	//	log.Fatalf("printRoute.err:%v", err)
	//}

}

// 服务端流式rpc
func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		time.Sleep(5 * time.Second)
		/**
		RecvMsg 会从流中读取完整的 gRPC 消息体:
		（1）RecvMsg 是阻塞等待的
		（2）RecvMsg 当流成功/结束（调用了 Close）时，会返回 io.EOF
		（3）RecvMsg 当流出现任何错误时，流会被中止，错误信息会包含 RPC 错误码。而在 RecvMsg 中可能出现如下错误：
			io.EOF
			io.ErrUnexpectedEOF
			transport.ConnectionError
			google.golang.org/grpc/codes
		需要注意，默认的 MaxReceiveMessageSize 值为 1024 * 1024 * 4，建议不要超出
		*/
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp:pj.name:%s, pt.value:%d", resp.Pt.Name, resp.Pt.Value)
	}
	return nil
}

// 客户端流式rpc
func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp: pj.name:%s, pt.value:%d", resp.Pt.Name, resp.Pt.Value)
	return nil
}

// 双向流式rpc
func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name:%s, pt.value :%d", resp.Pt.Name, resp.Pt.Value)
	}
	// 调用CloseSend()方法，就可以关闭服务端的stream，让它停止发送数据。值得注意的是，调用CloseSend()后，若继续调用Recv()，
	// 会重新激活stream，接着当前的结果继续获取消息
	stream.CloseSend()
	return nil
}
