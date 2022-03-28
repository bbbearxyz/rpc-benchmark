package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"math/rand"
	"net"
	"os"
	"rpc_test/src/grpc/proto"
	"runtime"
	"strconv"
	"time"
)

type MessageSender struct {
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

func main() {

	var args = os.Args
	// 参数介绍
	// 只有一个参数 num stream worker
	// 如果为0 表示 没有worker 每一次stream都有一个go routine
	// 如果为1 表示 stream worker的数量为逻辑cpu数量
	useStreamWorker, _ := strconv.ParseInt(args[1], 10, 64)
	if useStreamWorker == 1 {
		println("use stream worker")
	} else {
		println("don't use stream worker")
	}
	// 支持流式
	// 生成流式数据
	// 生成10mb 反复用
	data = GetRandomString(10 * 1024 * 1024)
	listen, err := net.Listen("tcp", "127.0.0.1:9527")
	if err != nil {
		log.Fatalf("tcp listen failed:%v", err)
	}
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 1024),
		grpc.MaxSendMsgSize(1024 * 1024 * 1024),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 10,
			Timeout: 3,
		}),
		grpc.InitialWindowSize(1024 * 1024 * 1024),
		grpc.InitialConnWindowSize(1024 * 1024 * 1024),
		grpc.WriteBufferSize(32 * 1024 * 1024),
		grpc.ReadBufferSize(32 * 1024 * 1024),
		grpc.MaxConcurrentStreams(512),
	}
	if useStreamWorker == 1 {
		opts = append(opts, grpc.NumStreamWorkers(uint32(runtime.NumCPU())))
	}
	server := grpc.NewServer(opts...)
	fmt.Println("services start success")
	proto.RegisterMessageSenderServer(server, &MessageSender{})
	server.Serve(listen)

}

func (sender *MessageSender) Send(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	time.Sleep(time.Duration(request.Time) * time.Millisecond)
	return &proto.MessageResponse{
		ResponseSomething: "",
	}, nil
}

func (sender *MessageSender) StreamTest(stream proto.MessageSender_StreamTestServer) error {
	round := 1 * 1024 / 10 + 1
	for i := 0; i < round; i ++ {
		if i == round - 1 {
			stream.Send(&proto.MessageResponse{ResponseSomething: data[0: 4 * 1024 * 1024], IsEnd: true})
			break
		}
		stream.Send(&proto.MessageResponse{ResponseSomething: data, IsEnd: false})
		println(i)
	}
	stream.Recv()
	return nil
}