package xlkafka

import "log"

const (
	raw       = "raw"
	rawAction = "raw_action"
)

type MsgDecoder interface {
	decode(msg []byte, out chan<- interface{})
}

func newMsgDecoder(name string) MsgDecoder {
	switch name {
	case raw:
		return newRawDecoder()
	default:
		log.Printf("Unknown kafka message decoder!")
	}
	return newRawDecoder()
}

type RawDecoder struct {
}

func newRawDecoder() MsgDecoder {
	d := new(RawDecoder)
	return d
}

func (d *RawDecoder) decode(msg []byte, out chan<- interface{}) {
	if len(msg) == 0 {
		return
	}
	out <- msg
}
