package controller

import (
	"context"
	"mygrpc/proto"
	"mygrpc/server/logic"
)

type Server struct {
	proto.UnimplementedGreeterServer
	cnt int64
}

func (s *Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func (s *Server) Exchange(ctx context.Context, in *proto.ExchangeParam) (*proto.Resp, error) {
	if in.From <= 0 || in.To <= 0 || in.Value <= 0 || len(in.Key) <= 0 {
		return &proto.Resp{Ret: -1}, nil
	}
	err := logic.Exchange(in.From, -in.Value, in.To, in.Value, in.Key)
	if err != nil {
		return &proto.Resp{Ret: -1}, err
	}
	return &proto.Resp{Ret: 1}, nil
}
