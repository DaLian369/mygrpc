package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"mygrpc/conf"
	"mygrpc/proto"
	"mygrpc/server/controller"
	"mygrpc/server/logic"
	"mygrpc/server/monitor"
	"mygrpc/server/xlkafka"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"net/http"
	_ "net/http/pprof"
)

var (
	port = flag.Int("port", 50051, "The server port")
)
var DefaultKafkaReader *xlkafka.KafkaReader

// 读配置done go pprof done  数据库 连接池 事务 缓存 消息队列 日志 监控 websocket prometheus

func main() {
	// go pprof信息
	go http.ListenAndServe(":9000", nil) // 127.0.0.1:9000/debug/pprof
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	// 获取配置文件路径
	var configPath string
	flag.StringVar(&configPath, "c", "", "config path")
	flag.Parse()
	if configPath == "" {
		PrintAndDie("not config")
	}
	// 加载配置文件
	cfg := new(proto.ConfigSt)
	err := conf.Init(configPath, cfg)
	if err != nil {
		PrintAndDie(err.Error())
	}

	// 初始化kafka写入
	err = xlkafka.InitWriter(&cfg.KafkaWriter)
	if err != nil {
		PrintAndDie(err.Error())
	}
	// 初始化kafka消费
	DefaultKafkaReader, err = xlkafka.NewKafakReader(&cfg.KafkaReader)
	if err != nil {
		PrintAndDie(err.Error())
	}
	if DefaultKafkaReader != nil {
		go KafkaReaderRun()
	}
	// 初始化prometheus
	err = monitor.InitPrometheus()
	if err != nil {
		PrintAndDie(err.Error())
	}

	// 初始化数据库
	err = logic.Init(cfg)
	if err != nil {
		PrintAndDie(err.Error())
	}

	// 监听tcp端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	// 创建grpc server
	s := grpc.NewServer()
	// 注册handler
	proto.RegisterGreeterServer(s, &controller.Server{})
	log.Printf("server listening at %v", lis.Addr())
	// grpcurl 启动反射服务，Protobuf 本身具有反射功能，可以在运行时获取对象的 Proto 文件。gRPC 同样也提供了一个名为 reflection 的反射包，用于为 gRPC 服务提供查询。
	reflection.Register(s)
	// sinal chan
	singalChan := make(chan os.Signal)
	signal.Notify(singalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGTSTP)
	// 服务启动
	go s.Serve(lis)
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	<-singalChan
	s.Stop()
	stopAll()
}

func KafkaReaderRun() {
	log.Println("server run success")
	for d := range DefaultKafkaReader.Open() {
		data, ok := d.([]byte)
		if !ok {
			log.Printf("bad log entry: %v", d)
			continue
		}
		o := proto.Order{}
		json.Unmarshal(data, &o)
		log.Printf("kafka reader entry: %+v", o)
		// 在这里如果要保证多线程
	}
}

func stopAll() {
	DefaultKafkaReader.Close()
}

func PrintAndDie(msg string) {
	// 报告函数调用信息
	pc, file, line, _ := runtime.Caller(1)
	log.Printf("err: %v, file %s, func: %s, line %d, %s", os.Stderr, file, runtime.FuncForPC(pc).Name(), line, msg)
	time.Sleep(1)
	// 程序退出，0表示成功，非0表示失败
	os.Exit(1)
}
