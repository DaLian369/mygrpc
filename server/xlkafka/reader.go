package xlkafka

import (
	"fmt"
	"log"
	"mygrpc/proto"
	"runtime"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type none struct{}

type KafkaReader struct {
	consumer    *cluster.Consumer
	msgDecoder  MsgDecoder
	dying, dead chan none
}

func NewKafakReader(cfg *proto.KafkaReaderConfigSt) (reader *KafkaReader, err error) {
	hostports := strings.Split(cfg.Hostports, ",")
	if len(hostports) == 0 {
		err = fmt.Errorf("unexpected kafka hostports: %s", hostports)
		log.Printf("new kafka reader err: %v", err)
		return
	}
	if len(cfg.GroupId) == 0 {
		err = fmt.Errorf("unexpected kafka groupid: %s", cfg.GroupId)
		log.Printf("%v", err)
		return
	}
	if len(cfg.Topic) == 0 {
		err = fmt.Errorf("unexpected kafka topics: %s", cfg.Topic)
		log.Printf("%v", err)
		return
	}
	topics := strings.Split(cfg.Topic, ",")

	reader = &KafkaReader{
		dying: make(chan none),
		dead:  make(chan none),
	}
	reader.msgDecoder = newMsgDecoder(cfg.MsgDecoder)

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Retry.Backoff = time.Millisecond * 500
	config.Group.Return.Notifications = true // 重平衡通知
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.CommitInterval = time.Second
	// config.Consumer.Offsets.AutoCommit.Enable = false // 关闭自动提交
	reader.consumer, err = cluster.NewConsumer(hostports, cfg.GroupId, topics, config)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	return
}

func (r *KafkaReader) Open() <-chan interface{} {
	output := make(chan interface{})
	go func() {
		for err := range r.consumer.Errors() {
			log.Printf("%v", err)
		}
	}()
	go func() {
		for noti := range r.consumer.Notifications() {
			log.Printf("%v", noti)
		}
	}()
	go r.read(output)
	return output
}

func (r *KafkaReader) read(output chan interface{}) {
	defer close(r.dead)
	defer close(output)
	defer r.consumer.Close()
	for {
		select {
		case <-r.dying:
			return
		case msg, ok := <-r.consumer.Messages():
			if !ok {
				return
			}
			r.process(msg, output)
		}
	}
}

func (r *KafkaReader) process(msg *sarama.ConsumerMessage, output chan interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			buf := make([]byte, 10240)
			runtime.Stack(buf, false)
			log.Printf("stack: %v,\n%s", err, string(buf))
			err = fmt.Errorf("send msg panic: %v", err)
		}
	}()
	defer r.consumer.MarkOffset(msg, "")
	r.msgDecoder.decode(msg.Value, output)
}

func (r *KafkaReader) Close() {
	// 这里关闭dying，阻塞，等待read()方法里关闭dead，消费完消息再返回
	select {
	case <-r.dying:
		return
	default:
		close(r.dying)
	}
	<-r.dead
}
