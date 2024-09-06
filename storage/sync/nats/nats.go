package nats

import (
	sPb "github.com/c12s/scheme/stellar"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
)

type NatsSync struct {
	nc    *nats.Conn
	topic string
}

func NewNatsSync(address, topic string) (*NatsSync, error) {
	nc, err := nats.Connect(address)
	if err != nil {
		return nil, err
	}

	return &NatsSync{
		nc:    nc,
		topic: topic,
	}, nil
}

func (ns *NatsSync) Sub(f func(u *sPb.LogBatch)) {
	ns.nc.Subscribe(ns.topic, func(msg *nats.Msg) {
		data := &sPb.LogBatch{}
		err := proto.Unmarshal(msg.Data, data)
		if err != nil {
			f(nil)
		}
		f(data)
	})
	ns.nc.Flush()
}
