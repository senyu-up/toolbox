package aws_kafka

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/storage"
)

var (
	connKafka = Kafka{}
	accessId  = "AKIAXYISDYYBBS"
	secret    = "/F27UkHBbkq3SpO9KA"
	brokers   = []string{"b-2-public.testkafkaiam.kafka.us-west-2.amazonaws.com:9198",
		"b-1-public.testkafkaiam.kafka.us-west-2.amazonaws.com:9198"}
	topic1 = "test_topic"
	topic2 = "test_topic22"
)

func SetUp() {
	var ctx = context.Background()
	var conf = &config.AwsKafkaConfig{
		Brokers: brokers,
		Topic:   topic1,
	}
	conf.SASL.Enable = true
	conf.SASL.Mechanism = "AWS_MSK_IAM"
	conf.SASL.Region = storage.AWSUwRegion
	conf.SASL.AccessId = accessId
	conf.SASL.SecretKey = secret
	k, err := New(ctx, conf)
	log.Printf("kafka conn: %v, err: %v", k, err)
	connKafka = *k
}

func TestKafka_CreateTopic(t *testing.T) {
	SetUp()
	type fields struct {
		consRwLock sync.RWMutex
		conf       config.KafkaConfig
		Conn       *kafka.Conn
	}
	type args struct {
		name           string
		numPartition   int32
		numReplication int16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1", fields: fields{}, wantErr: false, args: args{
				name:           topic2,
				numPartition:   1,
				numReplication: 1,
			},
		},
		{ // 重复创建不会报错
			name: "1-1", fields: fields{}, wantErr: false, args: args{
				name:           topic2,
				numPartition:   1,
				numReplication: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := connKafka.CreateTopic(tt.args.name, tt.args.numPartition, tt.args.numReplication); (err != nil) != tt.wantErr {
				t.Errorf("CreateTopic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKafka_DelTopic(t *testing.T) {
	SetUp()
	type fields struct {
		consRwLock sync.RWMutex
		conf       config.KafkaConfig
		Conn       *kafka.Conn
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1", fields: fields{}, wantErr: false, args: args{
				name: topic2,
			},
		},
		{
			name: "2", fields: fields{}, wantErr: true, args: args{
				name: "test_topic22333",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := connKafka.DelTopic(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DelTopic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKafka_DescribeTopic(t *testing.T) {
	SetUp()

	type fields struct {
		consRwLock sync.RWMutex
		conf       config.KafkaConfig
		Conn       *kafka.Conn
	}
	type args struct {
		name string
	}
	var brokers = []kafka.Broker{
		{
			"b-2-public.testkafkaiam.34xcfb.c13.kafka.us-west-2.amazonaws.com",
			9198,
			2,
			"usw2-az2",
		},
		{
			"b-1-public.testkafkaiam.34xcfb.c13.kafka.us-west-2.amazonaws.com",
			9198,
			2,
			"usw2-az1",
		},
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantDetail kafka.Partition
		wantErr    bool
	}{
		{
			name: "1", fields: fields{}, wantErr: true, wantDetail: kafka.Partition{}, args: args{
				name: topic2,
			},
		},
		{
			name: "2", fields: fields{}, wantErr: false, wantDetail: kafka.Partition{
				Topic: "xxx-topic",
				ID:    1,
				Leader: kafka.Broker{
					"b-2-public.testkafkaiam.34xcfb.c13.kafka.us-west-2.amazonaws.com",
					9198,
					2,
					"usw2-az2",
				},
				Replicas: brokers,
			}, args: args{
				name: "xxx-topic",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := connKafka.DescribeTopic(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(gotDetail.Topic == tt.wantDetail.Topic &&
				gotDetail.ID == tt.wantDetail.ID &&
				gotDetail.Leader == tt.wantDetail.Leader) {
				t.Errorf("DescribeTopic() gotDetail = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestKafka_ListTopics(t *testing.T) {
	SetUp()

	type fields struct {
		consRwLock sync.RWMutex
		conf       config.KafkaConfig
		Conn       *kafka.Conn
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]kafka.Partition
		wantErr bool
	}{
		{
			name: "1", fields: fields{}, wantErr: false, want: map[string]kafka.Partition{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connKafka.ListTopics()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTopics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) <= 0 {
				//if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTopics() got = %v, want %v", got, tt.want)
			} else {
				for _, tp := range got {
					t.Logf("ListTopics() got = %v \n", tp)
				}
			}
		})
	}
}
