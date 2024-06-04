package su_error

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
)

type SUError struct {
	Code int32
	Msg  string
}

func (s *SUError) Error() string {
	return fmt.Sprintf("%d$%s", s.Code, s.Msg)
}

func New(code int32, msg string) error {
	return &SUError{
		Code: code,
		Msg:  msg,
	}
}

func NewSUError(code int32, msg string) *SUError {
	return &SUError{
		Code: code,
		Msg:  msg,
	}
}

func NewWithError(code int32, err error) error {
	return &SUError{
		Code: code,
		Msg:  err.Error(),
	}
}

// GetSUError 尝试转换获取 SUError 类型的 err
func GetSUError(err error) (*SUError, bool) {
	if val, ok := err.(*SUError); ok {
		return val, true
	} else {
		return nil, false
	}
}

func Parse(err error) (code int32, msg string) {
	s := err.Error()
	tmp := strings.SplitN(s, "$", 2)
	if len(tmp) == 2 {
		c, _ := strconv.Atoi(tmp[0])
		return int32(c), tmp[1]
	} else {
		return 0, s
	}
}

func IgnoreNoRecord(err error) error {
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err == redis.Nil {
		return nil
	} else if err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}

// EmptyError
// @description 判断data是否为空, 如果为空则按照 err, msg 的优先级返回对应的错误
func EmptyError(data interface{}, err error, module string, msg string, ignoreNoRecord bool) error {
	if ignoreNoRecord {
		err = IgnoreNoRecord(err)
	}

	if err != nil {
		if module != "" {
			return fmt.Errorf("module:%s err:%s", module, err.Error())
		} else {
			return err
		}
	}

	if validator.IsNilInterface(data) {
		if module != "" {
			return fmt.Errorf("moduel:%s err:%s", module, msg)
		} else {
			return errors.New(msg)
		}
	}

	return nil
}

type Entry struct {
	Key   string
	Value interface{}
}

func Sprintf(err error, args ...interface{}) error {
	return fmt.Errorf(err.Error(), args...)
}

func Wrap(err error, entries ...*Entry) error {
	if len(entries) == 0 {
		return err
	}
	msg := strings.Builder{}
	msg.WriteString("[param] ")
	for i, _ := range entries {
		if i > 0 {
			msg.WriteString("&")
		}
		if entries[i].Key != "" {
			msg.WriteString(entries[i].Key)
			msg.WriteString("=")
		}

		// 常见数据类型推断处理
		if entries[i].Value == nil {
			msg.WriteString("nil")
			continue
		} else if s, ok := entries[i].Value.(string); ok {
			msg.WriteString(s)
			continue
		} else if b, ok := entries[i].Value.([]byte); ok {
			msg.Write(b)
			continue
		} else if bl, ok := entries[i].Value.(bool); ok {
			if bl {
				msg.WriteString("true")
			} else {
				msg.WriteString("false")
			}
		} else {
			tpe := reflect.TypeOf(entries[i].Value)
			if tpe.Kind() == reflect.Ptr {
				tpe = tpe.Elem()
			}
			if tpe.Kind() == reflect.Ptr {
				tpe = tpe.Elem()
			}
			switch tpe.Kind() {
			case reflect.Struct, reflect.Map, reflect.Slice:
				curV, _ := jsoniter.MarshalToString(entries[i].Value)
				msg.WriteString(curV)
			default:
				msg.WriteString(fmt.Sprintf("%+v", entries[i].Value))
			}
		}
	}

	return errors.New(err.Error() + " " + msg.String())
}
