package main

import (
	"golang.org/x/net/context"
	"rpc_test/src/kitex/proto/kitex_gen/message_proto"
	"rpc_test/src/kitex/proto/kitex_gen/message_proto/messagesender"
	"time"
)
type handler struct{}

func (handler *handler) Send(ctx context.Context, req *message_proto.MessageRequest) (*message_proto.MessageResponse, error) {
	time.Sleep(time.Duration(req.Time) * time.Millisecond)
	return &message_proto.MessageResponse{ResponseSomething: ""}, nil
}

func (handler *handler) StreamTest(stream message_proto.MessageSender_StreamTestServer) error {
	return nil
}

func main() {
	svr := messagesender.NewServer(&handler{})
	svr.Run()
}
