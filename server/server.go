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

func (Sender *MessageSender) Send(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	time.Sleep(time.Duration(request.Time) * time.Millisecond)
	return &proto.MessageResponse{
		ResponseSomething: "",
	}, nil
}