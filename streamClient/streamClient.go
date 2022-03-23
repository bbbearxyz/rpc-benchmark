package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"rpc_test/src/proto"
	"time"
)
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
func Task(taskNumber int) {
	// 8 * 64
	for i := 0; i < taskNumber; i ++ {
		go func() {
			addr := "localhost:9527"
			// Set up a connection to the server.
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			client := proto.NewMessageSenderClient(conn)
			// x := GetRandomString(256)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			stream, _ := client.StreamTest(ctx)
			stream.Send(&proto.MessageRequest{SaySomething: "hello"})
			res, _ := stream.Recv()
			println(res.ResponseSomething)
			stream.CloseSend()
		}()
	}
}
func main() {
	Task(1)
	for true {

	}
}
