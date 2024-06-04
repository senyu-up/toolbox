package ks

import (
	"github.com/aws/aws-sdk-go/service/firehose"
	"time"
)

// IFirehose aws firehose handler interface
type IFirehose interface {
	PushBatch(name string, record [][]byte, duration time.Duration) error
	Records(name string, record [][]byte) []*firehose.Record
}
