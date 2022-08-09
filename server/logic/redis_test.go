package logic

import (
	"mygrpc/proto"
	"testing"
)

func TestGet(t *testing.T) {
	DefaultProc = &ProcessSt{}
	redis := &proto.RedisSt{Hostport: "127.0.0.1:6379", Poolsize: 1, Timeout: 10}
	InitCache(redis, DefaultProc)
	_, err := DefaultProc.cache.GetData("1")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}
func TestSet(t *testing.T) {
	DefaultProc = &ProcessSt{}
	redis := &proto.RedisSt{Hostport: "127.0.0.1:6379", Poolsize: 1, Timeout: 10}
	InitCache(redis, DefaultProc)
	err := DefaultProc.cache.Set("1", "3")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}
