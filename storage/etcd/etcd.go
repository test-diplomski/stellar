package etcd

import (
	"context"
	"fmt"
	sPb "github.com/c12s/scheme/stellar"
	s "github.com/c12s/stellar-go"
	"github.com/c12s/stellar/model"
	sync "github.com/c12s/stellar/storage/sync"
	"github.com/c12s/stellar/storage/sync/nats"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"time"
)

type DB struct {
	Kv     clientv3.KV
	Client *clientv3.Client
	s      sync.Syncer
}

func New(conf *model.Config, timeout time.Duration) (*DB, error) {
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: timeout,
		Endpoints:   conf.Endpoints,
	})

	if err != nil {
		return nil, err
	}

	ns, err := nats.NewNatsSync(conf.Syncer, conf.STopic)
	if err != nil {
		return nil, err
	}

	return &DB{
		Kv:     clientv3.NewKV(cli),
		Client: cli,
		s:      ns,
	}, nil
}

func (db *DB) List(ctx context.Context, req *sPb.ListReq) (*sPb.ListResp, error) {
	span, _ := s.FromGRPCContext(ctx, "stellar.list")
	defer span.Finish()
	fmt.Println(span)

	// Pass one: do Query by kv pairs sent from gateway
	// all pairs sent, are used as a lookup so all pairs
	// must be present in query item

	childSpan := span.Child("etcd.Get.WithPrefix.WithSort")
	resp, err := db.Kv.Get(ctx, lookupKey(), clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		childSpan.AddLog(&s.KV{"etvd error", err.Error()})
		return nil, err
	}
	childSpan.Finish()
	fmt.Println(childSpan)

	index := []string{}
	for _, item := range resp.Kvs {
		elem := &sPb.Tags{}
		err = proto.Unmarshal(item.Value, elem)
		if err != nil {
			span.AddLog(&s.KV{"unmarshall error", err.Error()})
			return nil, err
		}

		for key, value := range req.Query {
			if val, ok := elem.Tags[key]; !ok || val != value {
				continue
			}
		}
		index = append(index, extractTraceKey(string(item.Key)))
	}

	// Pass two: based on index of keys that satisfied query
	// extrat spans and return tham
	traces := []*sPb.GetResp{}
	for _, key := range index {
		t, err := db.Get(s.NewTracedContext(ctx, span), &sPb.GetReq{TraceId: key})
		if err != nil {
			span.AddLog(&s.KV{"stellar.get error", err.Error()})
			return nil, err
		}
		traces = append(traces, t)
	}

	return &sPb.ListResp{
		Traces: traces,
	}, nil
}

func (db *DB) Get(ctx context.Context, req *sPb.GetReq) (*sPb.GetResp, error) {
	span, _ := s.FromGRPCContext(ctx, "stellar.get")
	defer span.Finish()
	fmt.Println(span)

	childSpan := span.Child("etcd.Get.WithPrefix.WithSort")
	trace := []*sPb.Span{}
	resp, err := db.Kv.Get(ctx, traceKey(req.TraceId), clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		childSpan.AddLog(&s.KV{"etcd get error", err.Error()})
		return nil, err
	}
	childSpan.Finish()
	fmt.Println(childSpan)

	for _, item := range resp.Kvs {
		elem := &sPb.Span{}
		err = proto.Unmarshal(item.Value, elem)
		if err != nil {
			span.AddLog(&s.KV{"unmarshall err", err.Error()})
			return nil, err
		}
		trace = append(trace, elem)
	}

	return &sPb.GetResp{
		Trace: trace,
	}, nil
}

func (db *DB) StartCollector(ctx context.Context) {
	db.s.Sub(func(msg *sPb.LogBatch) {
		go func(batch *sPb.LogBatch) {
			for _, log := range batch.Batch {
				logData, _ := proto.Marshal(log)
				keyParts := []string{log.SpanContext.TraceId, log.SpanContext.SpanId}
				key := formKey(keyParts)
				fmt.Println(key)
				_, err := db.Kv.Put(ctx, key, string(logData))
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}(msg)
	})
}
