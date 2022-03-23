package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"os"
	"rpc_test/src/proto"
	"strconv"
	"sync/atomic"
	"time"
)

var success int64
var task int64
var field int64
var payload int64
var latency int64
var testTime int64

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
			request := &proto.MessageRequest{}
			request.Time = latency
			if field == 1 {
				request.Field1 = GetRandomString(int(payload - 8))
			} else if field == 5 {
				averageLen := (payload - 8) / field
				request.Field1 = GetRandomString(int(averageLen))
				request.Field2 = GetRandomString(int(averageLen))
				request.Field3 = GetRandomString(int(averageLen))
				request.Field4 = GetRandomString(int(averageLen))
				request.Field5 = GetRandomString(int((payload - 8) - 4 * averageLen))
			} else if field == 10 {
				averageLen := (payload - 8) / field
				request.Field1 = GetRandomString(int(averageLen))
				request.Field2 = GetRandomString(int(averageLen))
				request.Field3 = GetRandomString(int(averageLen))
				request.Field4 = GetRandomString(int(averageLen))
				request.Field5 = GetRandomString(int(averageLen))
				request.Field6 = GetRandomString(int(averageLen))
				request.Field7 = GetRandomString(int(averageLen))
				request.Field8 = GetRandomString(int(averageLen))
				request.Field9 = GetRandomString(int(averageLen))
				request.Field10 = GetRandomString(int((payload - 8) - 9 * averageLen))
			} else {
				return
			}
			// defer conn.Close()
			// defer cancel()
			for true {
				conn, _ := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				c := proto.NewMessageSenderClient(conn)
				ctx, _ := context.WithTimeout(context.Background(), time.Duration(testTime + 10))
				c.Send(ctx, request)
				success = atomic.AddInt64(&success, 1)
				conn.Close()
			}
		}()
	}
}
func main() {
	var args = os.Args
	// 0 -> task 1 -> pb field number 2 -> payload 3 -> latency(ms) 4 -> testTime
	task, _ = strconv.ParseInt(args[1], 10, 64)
	field, _ = strconv.ParseInt(args[2], 10, 64)
	payload, _ = strconv.ParseInt(args[3], 10, 64)
	latency, _ = strconv.ParseInt(args[4], 10, 64)
	testTime, _ = strconv.ParseInt(args[5], 10, 64)
	println("task is ", task, " field is ", field, " payload is ", payload, "B latency is ", latency, "ms testTime is ", testTime, "s")
	Task(int(task))
	ch := time.After(time.Duration(testTime) * time.Second)
	select {
	case <- ch:
	}
	println("throughput is ", success * payload / (testTime * 1024), " KB/s")
}
