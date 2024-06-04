package marshaler

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/packet"
	"reflect"
)

var (
	ErrInvalidMsgType  = fmt.Errorf("invalid message type")
	ErrUnregisteredMsg = fmt.Errorf("unregistered message")
)

func Marshal(i proto.Message) []byte {
	data, _ := proto.Marshal(i)
	return data
}

func MarshalToString(i proto.Message) string {
	return proto.MarshalTextString(i)
}

func Unmarshal(b []byte, i proto.Message) {
	if err := proto.Unmarshal(b, i); err != nil {
		logger.Error("unmarshal err: ", err)
	}
}

type ProtoMarshaler struct{}

func (ProtoMarshaler) Marshal(i proto.Message) ([]byte, error) {
	return proto.Marshal(i)
}

func (ProtoMarshaler) Unmarshal(b []byte, i proto.Message) error {
	return proto.Unmarshal(b, i)
}

type Factory func() proto.Message

var msgFactories = map[string]Factory{}

func Register(factory Factory) {
	name := GetType(factory())
	if _, ok := msgFactories[name]; ok {
		panic(fmt.Sprintln("duplicate message factory", name))
	}
	msgFactories[name] = factory
}

func Encode(pb interface{}) ([]byte, error) {
	if msg, ok := pb.(proto.Message); ok {
		buffer := packet.Writer()
		data, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(GetType(msg))
		buffer.WriteRawBytes(data)
		return buffer.Data(), err
	}
	return nil, ErrInvalidMsgType
}

func Decode(data []byte) (interface{}, error) {
	buffer := packet.Reader(data)
	name, err := buffer.ReadString()
	if err != nil {
		return nil, err
	}
	factory, ok := msgFactories[name]
	if !ok {
		return nil, ErrUnregisteredMsg
	}
	msg := factory()
	return msg, proto.Unmarshal(buffer.RemainData(), msg)
}

func DecodeWithOut(data []byte, out proto.Message) error {
	buffer := packet.Reader(data)
	_, err := buffer.ReadString()
	if err != nil {
		return err
	}
	return proto.Unmarshal(buffer.RemainData(), out)
}

func GetType(msg interface{}) string {
	if t := reflect.TypeOf(msg); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
