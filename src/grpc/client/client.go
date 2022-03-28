package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"os"
	"rpc_test/src/grpc/proto"
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
var addr string
var isStream int64
var start time.Time
var end time.Time
var data []proto.MessageRequest

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
	for i := 0; i < taskNumber; i ++ {
		go func(taskId int) {
			request := &data[taskId]
			request.Time = latency
			for true {
				conn, _ := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				c := proto.NewMessageSenderClient(conn)
				ctx, _ := context.WithTimeout(context.Background(), time.Duration(testTime + 10) * time.Second)
				c.Send(ctx, request)
				success = atomic.AddInt64(&success, 1)
				conn.Close()
			}
		}(i)
	}
}
func StreamTask() {
	request := &data[0]
	conn, _ := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(testTime + 10) * time.Second)
	client := proto.NewMessageSenderClient(conn)
	stream, _ := client.StreamTest(ctx)
	stream.Send(request)
	for true {
		res, _ := stream.Recv()
		if res.IsEnd {
			break
		}
	}
	stream.CloseSend()
	conn.Close()
}
func GenerateData() {
	for i := 0; i < int(task); i ++ {
		request := proto.MessageRequest{}
		if field == 1 {
			request.Field1 = GetRandomString(int(payload))
		} else if field == 5 {
			averageLen := (payload) / field
			request.Field1 = GetRandomString(int(averageLen))
			request.Field2 = GetRandomString(int(averageLen))
			request.Field3 = GetRandomString(int(averageLen))
			request.Field4 = GetRandomString(int(averageLen))
			request.Field5 = GetRandomString(int((payload) - 4 * averageLen))
		} else if field == 10 {
			averageLen := (payload) / field
			request.Field1 = GetRandomString(int(averageLen))
			request.Field2 = GetRandomString(int(averageLen))
			request.Field3 = GetRandomString(int(averageLen))
			request.Field4 = GetRandomString(int(averageLen))
			request.Field5 = GetRandomString(int(averageLen))
			request.Field6 = GetRandomString(int(averageLen))
			request.Field7 = GetRandomString(int(averageLen))
			request.Field8 = GetRandomString(int(averageLen))
			request.Field9 = GetRandomString(int(averageLen))
			request.Field10 = GetRandomString(int((payload) - 9 * averageLen))
		}
		data = append(data, request)
	}
}
func main() {
	var args = os.Args

	// 参数介绍
	// 第一个参数 addr 表示server的地址
	// is stream -> 如果是streaming 执行流程是发一个256字节的request
	// response每次回10mb 一直到1GB
	// stream -> 如果不是streaming 后面有几个参数  task -> 表示并行任务的数量
	// field number -> 表示pb字段的数量(1, 5, 10)
	// payload 表示每个请求中携带的数据的大小 单位为B
	// latency 希望在server端增加的延迟的时间 单位为ms
	// testTime 表示测试的时间 单位为s
	addr = args[1]
	isStream, _ = strconv.ParseInt(args[2], 10, 64)
	if isStream > 0 {
		// 生成数据
		data = append(data, proto.MessageRequest{
			Field1: GetRandomString(256),
		})
		println("streaming test")
		println("addr is ", addr)
		println("test size is 1GB")
		start = time.Now()
		StreamTask()
		end = time.Now()
		cost := end.Sub(start).Seconds()
		println("throughput is ", float64(1 * 1024 * 1024) / cost, " KB/s")
	} else {
		task, _ = strconv.ParseInt(args[3], 10, 64)
		field, _ = strconv.ParseInt(args[4], 10, 64)
		payload, _ = strconv.ParseInt(args[5], 10, 64)
		latency, _ = strconv.ParseInt(args[6], 10, 64)
		testTime, _ = strconv.ParseInt(args[7], 10, 64)
		// 生成数据
		GenerateData()
		println("no streaming test")
		println("addr is ", addr)
		println("task is ", task, " field is ", field, " payload is ", payload, "B latency is ", latency, "ms testTime is ", testTime, "s")
		Task(int(task))
		ch := time.After(time.Duration(testTime) * time.Second)
		select {
		case <- ch:
		}
		println("throughput is ", success * payload / (testTime * 1024), " KB/s")

	}


}
