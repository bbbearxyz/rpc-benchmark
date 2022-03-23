package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"rpc_test/src/proto"
	"time"
)

type MessageSender struct{
	proto.UnimplementedMessageSenderServer
}

var num int

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:9527")
	if err != nil {
		log.Fatalf("tcp listen failed:%v", err)
	}
	server := grpc.NewServer()
	fmt.Println("services start success")
	proto.RegisterMessageSenderServer(server, &MessageSender{})
	server.Serve(listen)

}

func (sender *MessageSender) Send(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	num += 1
	println("num wait ", num)
	time.Sleep(10 * time.Second)
	return &proto.MessageResponse{
		ResponseSomething: request.SaySomething + "yyy",
	}, nil
}

func (sender *MessageSender) StreamTest(stream proto.MessageSender_StreamTestServer) error {
	for {
		stream.Send(&proto.MessageResponse{ResponseSomething: "world"})
		res, _ := stream.Recv()
		println(res.SaySomething)
	}
}