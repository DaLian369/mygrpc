package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mygrpc/proto"
	"sync"
	"time"

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

	var from int64 = 1
	var to int64 = 2
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// 发送请求并接受响应
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name, Id: from})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s, %s", r.GetMessage(), r.GetUuid())

	// 并发转账
	limit := 1
	wg := sync.WaitGroup{}
	wg.Add(limit)
	for i := 0; i < limit; i++ {
		go func(t int) {
			defer wg.Done()
			er, err := c.Exchange(ctx, &proto.ExchangeParam{From: from, To: to, Value: 1, Key: r.GetUuid()})
			if err != nil {
				fmt.Println("could not exchange:", err, t)
			}
			log.Printf("Exchange: %d %d", er.GetRet(), t)
		}(i)
	}
	wg.Wait()
}
