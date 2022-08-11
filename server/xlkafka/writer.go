package xlkafka

import (
	"errors"
	"fmt"
	"log"
	"mygrpc/proto"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var (
	bestKafkaBatchSize = 128
	maxKafkaBatchSize  = 2048
	defaultChannelSize = 1024 * 10
)

var (
	errNotInitialized = errors.New("kafka is not initialized")
	errNotRunning     = errors.New("kafka is not running")
)

type Message struct {
	Topic string
	Key   string
	Value string
}

type Writer struct {
	output   chan *Message
	waitRun  *sync.WaitGroup
	stopChan chan bool
	running  bool
}

var defaultWriter *Writer

func InitWriter(cfg *proto.KafkaSt) (err error) {
	if !cfg.Enable {
		return
	}
	defaultWriter, err = NewKafkaProducer(cfg)
	return
}

func NewKafkaProducer(cfg *proto.KafkaSt) (w *Writer, err error) {
	w = new(Writer)
	err = w.init(cfg)
	if err != nil {
		log.Printf("NewKafkaProducer init err: %v", err)
		return
	}
	return
}

func SendMsg(topic, key, value string) error {
	if defaultWriter == nil {
		return errNotInitialized
	}
	return defaultWriter.SendMsg(topic, key, value)
}

func (w *Writer) SendMsg(topic, key, value string) (err error) {
	if !w.running {
		log.Printf(errNotRunning.Error())
		return errNotRunning
	}
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 10240)
			runtime.Stack(buf, false)
			log.Printf("stack: %v,\n%s", err, string(buf))
			err = fmt.Errorf("send msg panic: %v", r)
		}
	}()
	msg := &Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}
	w.output <- msg
	return
}

func (w *Writer) init(cfg *proto.KafkaSt) (err error) {
	w.output = make(chan *Message, defaultChannelSize)
	w.waitRun = new(sync.WaitGroup)
	w.stopChan = make(chan bool)

	hostports := strings.Split(cfg.Hostports, ",")
	if len(hostports) == 0 {
		err = fmt.Errorf("hostports is nil")
		log.Println(err)
		return
	}
	concurrent := int(cfg.Concurrent)
	if concurrent < 1 {
		concurrent = 1
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Errors = true
	config.Producer.Flush.Messages = bestKafkaBatchSize
	config.Producer.Flush.MaxMessages = maxKafkaBatchSize // 最大消息数
	config.Producer.Flush.Frequency = time.Second
	config.Producer.Flush.Bytes = 1024 * 1024 // 最大字节数
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Retry.Max = 6
	config.Producer.Retry.Backoff = time.Millisecond * 250

	producer, err := sarama.NewAsyncProducer(hostports, config)
	if err != nil {
		log.Printf("new producer err: %v", err)
		return
	}
	w.waitRun.Add(1)
	go w.processError(producer)
	for i := 0; i < int(cfg.Concurrent); i++ {
		w.waitRun.Add(1)
		go w.runKafka(producer, w.output)
	}
	log.Printf("kafka writer running, %d goroutines, hostports: %+v", concurrent, hostports)
	w.running = true
	return
}

func (w *Writer) runKafka(producer sarama.AsyncProducer, c <-chan *Message) {
	defer w.waitRun.Done()
	for msg := range c {
		key := sarama.StringEncoder(msg.Key)
		value := sarama.StringEncoder(msg.Value)
		message := &sarama.ProducerMessage{
			Key:   key,
			Value: value,
			Topic: msg.Topic,
		}
		producer.Input() <- message
	}
}

func (w *Writer) processError(producer sarama.AsyncProducer) {
	defer w.waitRun.Done()
	for {
		select {
		case err := <-producer.Errors():
			if err != nil && err.Msg != nil {
				log.Printf("Kafka write fail, topic %s, partition %d, length %d, error: %s",
					err.Msg.Topic, err.Msg.Partition, err.Msg.Value.Length(), err.Error())
			}
		case <-w.stopChan:
			err := producer.Close()
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
}
