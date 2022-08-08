package main

import (
	"context"
	"flag"
	"log"
	"mygrpc/proto"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	//  创建连接
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 创建客户端
	c := proto.NewGreeterClient(conn)

	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 发送请求并接受响应
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	er, err := c.Exchange(ctx, &proto.ExchangeParam{From: 1, To: 2, Value: 1, Key: uuid.New().String()})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Exchange: %d", er.GetRet())
}
