package logger

import (
	"log"
	"mygrpc/proto"
	"os"
)

var Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile)

func Init(cfg *proto.ConfigSt) (err error) {
	return
}
