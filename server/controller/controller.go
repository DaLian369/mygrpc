package controller

import (
	"context"
	"encoding/json"
	"log"
	"mygrpc/proto"
	"mygrpc/server/logic"
	"mygrpc/server/monitor"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Server struct {
	proto.UnimplementedGreeterServer
	cnt int64
}

// 一元拦截器
func UnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	remote, _ := peer.FromContext(ctx)
	remoteAddr := remote.Addr.String()

	in, _ := json.Marshal(req)
	inStr := string(in)
	log.Println("ip", remoteAddr, "access_start", info.FullMethod, "in", inStr)

	// prometheus req 自增1
	monitor.Reqs.WithLabelValues(info.FullMethod).Observe(1)

	start := time.Now()
	defer func() {
		out, _ := json.Marshal(resp)
		outStr := string(out)
		duration := int64(time.Since(start) / time.Millisecond)
		log.Println("ip", remoteAddr, "access_end", info.FullMethod, "in", inStr, "out", outStr, "err", err, "duration/ms", duration)
	}()

	resp, err = handler(ctx, req)

	return
}

func UnaryServerInterceptor2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	log.Println("UnaryServerInterceptor2")
	resp, err = handler(ctx, req)

	return
}

func (s *Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	if in.Id <= 0 {
		return &proto.HelloReply{Message: "param failed"}, nil
	}
	// 应该限制同一用户的uuid没有使用时，下发老的uuid，否则生成
	uid, err := logic.DefaultProc.GenUUID(in.Id)
	if err != nil {
		return &proto.HelloReply{Message: "uuid failed"}, nil
	}
	return &proto.HelloReply{Message: "Hello again " + in.GetName(), Uuid: uid}, nil
}

func (s *Server) Exchange(ctx context.Context, in *proto.ExchangeParam) (*proto.Resp, error) {
	if in.From <= 0 || in.To <= 0 || in.Value <= 0 || len(in.Key) <= 0 {
		return &proto.Resp{Ret: -1}, nil
	}
	err := logic.DefaultProc.Exchange(in.From, -in.Value, in.To, in.Value, in.Key)
	if err != nil {
		return &proto.Resp{Ret: -1}, err
	}
	return &proto.Resp{Ret: 1}, nil
}
