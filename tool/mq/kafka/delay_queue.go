package kafka

import (
	"context"
	"github.com/IBM/sarama"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	error2 "github.com/senyu-up/toolbox/tool/su_error"
	"time"
)

const DelayTopicName = "su_delay_mq_13144131"

type DelayProducerMsg struct {
	// 必填, 时间到期后, 将消息投递的目的topic
	Topic string // The Kafka topic for this message.
	// 非必填, 消息id, 保证全局唯一, 为空会使用uuid.v4 生成
	Key string
	// 必填, 消息内容
	Value []byte
	// 非必填, 头部消息
	Headers []*sarama.RecordHeader
	// 非必填, 元数据
	Metadata interface{}
	// 非必填,  指定到期后投递到的分区
	Partition int32
	// 延迟推送时间点
	TimeTo int64
}

type DelayQueueMsg struct {
	Topic     string    `json:"T"`
	Partition int32     `json:"P"`
	Data      []byte    `json:"D"`
	TimeTo    time.Time `json:"TT"`
	GmtCreate time.Time `json:"GT"`
}

func delayMsgFilter(data DelayProducerMsg) (msg sarama.ProducerMessage, err error) {
	if data.Key == "" {
		data.Key = uuid.NewV4().String()
	}
	if data.Topic == "" {
		err = error2.New(500400, "topic为空")
		return
	}
	if data.Topic == DelayTopicName {
		err = error2.New(500400, "topic名称冲突")
		return
	}
	if data.Value == nil {
		err = error2.New(500400, "value内容为空")
		return
	}
	var timeTo time.Time
	timeTo = time.Unix(data.TimeTo, 0)

	byteData, _ := jsoniter.Marshal(data)

	delayData := DelayQueueMsg{
		Data:      byteData,
		Topic:     data.Topic,
		Partition: data.Partition,
		TimeTo:    timeTo,
		GmtCreate: time.Now(),
	}

	val, _ := jsoniter.Marshal(delayData)

	return sarama.ProducerMessage{
		Topic: DelayTopicName,
		Key:   sarama.ByteEncoder(data.Key),
		Value: sarama.ByteEncoder(val),
	}, nil
}

// BatchDelayAsyncSendV2
//
//	@Description: 异步批量发送 延迟消息
//	@receiver s
//	@param ctx  body any true "-"
//	@param list  body any true "-"
//	@return err
func (s *Producer) BatchDelayAsyncSendV2(ctx context.Context, list []DelayProducerMsg) (err error) {
	//var msgList = make([]*sarama.ProducerMessage, 0, len(list))
	for i, _ := range list {
		msg, err1 := delayMsgFilter(list[i])
		if err1 != nil {
			return err1
		}
		s.PushAsyncRaw(ctx, &msg)
	}
	return
}

// BatchDelaySyncSend
//
//	@Description: 同步批量发送 延迟消息
//	@receiver s
//	@param ctx  body any true "-"
//	@param list  body any true "-"
//	@return err
func (s *Producer) BatchDelaySyncSend(ctx context.Context, list []DelayProducerMsg) (err error) {
	var msgList = make([]*sarama.ProducerMessage, 0, len(list))
	for i, _ := range list {
		msg, err1 := delayMsgFilter(list[i])
		if err1 != nil {
			return err1
		}
		msgList = append(msgList, &msg)
	}
	return s.PushSyncRawMsgs(ctx, msgList)
}

// DelaySendV2
// @description 延迟推送, msgId 为消息id, 需保证全局唯一, 如果为空, 会基于uuid.v4自动生成
// d: 再当前时间的基础上进行增加
func (s *Producer) DelaySendV2(ctx context.Context, data DelayProducerMsg) (err error) {
	msg, err := delayMsgFilter(data)
	if err != nil {
		return err
	}
	s.PushAsyncRaw(ctx, &msg)
	return err
}

// DelaySend
// @description 延迟推送, msgId 为消息id, 需保证全局唯一, 如果为空, 会基于uuid.v4自动生成
// d: 再当前时间的基础上进行增加
func (s *Producer) DelaySend(ctx context.Context, data *DelayProducerMsg, d time.Duration) (err error) {
	if data == nil {
		return error2.New(500400, "data 不允许为空")
	}

	var timeTo = data.TimeTo
	if data.TimeTo == 0 {
		if d < time.Second {
			return error2.New(500400, "延迟需>=1秒")
		}
		timeTo = time.Now().Add(d).Unix()
	}

	data.TimeTo = timeTo
	msg, err := delayMsgFilter(*data)
	if err != nil {
		return err
	}

	s.PushAsyncRaw(ctx, &msg)

	return err
}
