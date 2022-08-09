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
