package logic

import (
	"mygrpc/proto"
	"testing"

	"github.com/google/uuid"
)

func TestExchange(t *testing.T) {
	Init(&proto.MysqlSt{Hostport: "127.0.0.1:3306", Password: "xuelei123", Username: "root", Poolsize: 1, Idlesize: 1, Database: "kv"})
	uuid := uuid.New().String()
	err := Exchange(1, -10, 2, 10, uuid)
	if err != nil {
		t.Errorf("%s %s", err.Error(), uuid)
	}
}
