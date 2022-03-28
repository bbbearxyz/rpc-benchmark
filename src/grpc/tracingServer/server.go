package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"rpc_test/src/grpc/proto"
	"time"
)


type MessageSender struct{
	proto.UnimplementedMessageSenderServer
}

var data string

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i ++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

const (
	ZIPKIN_RECORDER_HOST_PORT = "http://127.0.0.1:9411/api/v2/spans"
)

func main() {
	reporter := zipkinhttp.NewReporter(ZIPKIN_RECORDER_HOST_PORT)
	defer reporter.Close()

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)

	// 支持流式
	// 生成流式数据
	// 生成10mb 反复用
	data = GetRandomString(10 * 1024 * 1024)
	listen, err := net.Listen("tcp", "localhost:9527")
	if err != nil {
		log.Fatalf("tcp listen failed:%v", err)
	}
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		),
	}
	server := grpc.NewServer(opts...)
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

func (Sender *MessageSender) StreamTest(stream proto.MessageSender_StreamTestServer) error {
	round := 1 * 1024 / 10 + 1
	for i := 0; i < round; i ++ {
		if i == round - 1 {
			stream.Send(&proto.MessageResponse{ResponseSomething: data[0: 4 * 1024 * 1024], IsEnd: true})
			break
		}
		stream.Send(&proto.MessageResponse{ResponseSomething: data, IsEnd: false})
	}
	stream.Recv()
	return nil
}

