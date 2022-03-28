// Code generated by Kitex v0.2.1. DO NOT EDIT.

package messagesender

import (
	"github.com/cloudwego/kitex/server"
	"rpc_test/src/kitex/proto/kitex_gen/message_proto"
)

// NewInvoker creates a server.Invoker with the given handler and options.
func NewInvoker(handler message_proto.MessageSender, opts ...server.Option) server.Invoker {
	var options []server.Option

	options = append(options, opts...)

	s := server.NewInvoker(options...)
	if err := s.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	return s
}