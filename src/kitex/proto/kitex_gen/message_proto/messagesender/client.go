// Code generated by Kitex v0.2.1. DO NOT EDIT.

package messagesender

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/streaming"
	"github.com/cloudwego/kitex/transport"
	"rpc_test/src/kitex/proto/kitex_gen/message_proto"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Send(ctx context.Context, Req *message_proto.MessageRequest, callOptions ...callopt.Option) (r *message_proto.MessageResponse, err error)
	StreamTest(ctx context.Context, callOptions ...callopt.Option) (stream MessageSender_StreamTestClient, err error)
}

type MessageSender_StreamTestClient interface {
	streaming.Stream
	Send(*message_proto.MessageRequest) error
	Recv() (*message_proto.MessageResponse, error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, client.WithTransportProtocol(transport.GRPC))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kMessageSenderClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kMessageSenderClient struct {
	*kClient
}

func (p *kMessageSenderClient) Send(ctx context.Context, Req *message_proto.MessageRequest, callOptions ...callopt.Option) (r *message_proto.MessageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Send(ctx, Req)
}

func (p *kMessageSenderClient) StreamTest(ctx context.Context, callOptions ...callopt.Option) (stream MessageSender_StreamTestClient, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.StreamTest(ctx)
}
