package ks

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/senyu-up/toolbox/tool/logger"
	"sync"
	"time"
)

var FirehoseRecordNumberError = errors.New("[Firehose] Push record err. At least one is needed! ")

type FirehoseHandle struct {
	sess    *session.Session
	client  *firehose.Firehose
	carrier sync.Map
}

func InitFirehose(region string, id, secret string) (handler *FirehoseHandle, err error) {
	handler = &FirehoseHandle{}
	sess, err := GetAwsSession(region, id, secret)
	if err != nil {
		return nil, err
	}
	handler.sess = sess
	handler.client = firehose.New(sess)
	return handler, nil
}

func (p *FirehoseHandle) PushBatch(name string, record [][]byte, timeout time.Duration) error {
	if len(record) == 0 {
		logger.Info("[Firehose] Push record err. At least one is needed!")
		return FirehoseRecordNumberError
	}
	//withTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	//defer func() {
	//	cancel()
	//}()
	//_, err := p.client.PutRecordBatchWithContext(withTimeout, &firehose.PutRecordBatchInput{
	//	DeliveryStreamName: aws.String(name),
	//	Records:            p.Records(name, record),
	//})
	//默认 40s 超时
	_, err := p.client.PutRecordBatch(&firehose.PutRecordBatchInput{
		DeliveryStreamName: aws.String(name),
		Records:            p.Records(name, record),
	})
	if err != nil {
		return err
	}
	return err
}

func (p *FirehoseHandle) Records(name string, record [][]byte) []*firehose.Record {
	load, ok := p.carrier.Load(name)
	if !ok {
		w := &wrapper{}
		w.init(len(record))
		load = w
		p.carrier.Store(name, w)
	}
	for i, bytes := range record {
		load.(*wrapper).get(len(record))[i].Data = bytes
	}
	return load.(*wrapper).get(len(record))[:len(record)]
}

type wrapper struct {
	len    int
	record []*firehose.Record

	sync.Mutex
}

func (p *wrapper) init(len int) {
	p.len = len
	p.record = make([]*firehose.Record, len, len)
	for i, _ := range p.record {
		p.record[i] = &firehose.Record{}
	}
}

func (p *wrapper) cap(len int) {
	if p.expand(len) {
		p.Lock()
		defer p.Unlock()
		old := p.record
		p.record = make([]*firehose.Record, len, len)
		for i, _ := range p.record {
			p.record[i] = &firehose.Record{}
		}
		p.len = len
		copy(p.record, old)
	}
}

func (p *wrapper) expand(target int) bool {
	return p.len < target
}

func (p *wrapper) get(len int) []*firehose.Record {
	p.cap(len)
	p.Lock()
	defer p.Unlock()
	return p.record
}
