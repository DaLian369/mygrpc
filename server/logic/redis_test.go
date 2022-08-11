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

func TestSafeDel(t *testing.T) {
	DefaultProc = &ProcessSt{}
	redis := &proto.RedisSt{Hostport: "127.0.0.1:6379", Poolsize: 1, Timeout: 10}
	InitCache(redis, DefaultProc)
	c := DefaultProc.cache
	err := c.SafeDel(c.genUUIDKey(1), "97468eb7-d274-4789-8237-a9b3ccef7f81")
	if err != nil {
		t.Errorf("xxxx %s", err.Error())
	}
}
