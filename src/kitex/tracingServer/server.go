package main

import (
	"github.com/cloudwego/kitex/server"
	ko "github.com/kitex-contrib/tracer-opentracing"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"golang.org/x/net/context"
	"log"
	"rpc_test/src/kitex/proto/kitex_gen/message_proto"
	"rpc_test/src/kitex/proto/kitex_gen/message_proto/messagesender"
	"time"
)
type handler struct{}

const (
	ZIPKIN_RECORDER_HOST_PORT = "http://127.0.0.1:9411/api/v2/spans"
)

func (handler *handler) Send(ctx context.Context, req *message_proto.MessageRequest) (*message_proto.MessageResponse, error) {
	time.Sleep(time.Duration(req.Time) * time.Millisecond)
	return &message_proto.MessageResponse{ResponseSomething: ""}, nil
}

func (handler *handler) StreamTest(stream message_proto.MessageSender_StreamTestServer) error {
	return nil
}

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

	svr := messagesender.NewServer(&handler{}, server.WithSuite(ko.NewDefaultServerSuite()))
	svr.Run()
}
