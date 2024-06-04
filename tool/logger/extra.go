package logger

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/trace"
	"reflect"
	"time"
)

type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Float     float64
	Boolean   bool
	Bytes     []byte
	Interface interface{}
}

// A FieldType indicates which member of the Field union struct should be used
// and how it should be serialized.
type FieldType uint8

const (
	// UnknownType is the default field type. Attempting to add it to an encoder will panic.
	UnknownType FieldType = iota
	// ArrayMarshalerType indicates that the field carries an ArrayMarshaler.
	ArrayMarshalerType
	// ObjectMarshalerType indicates that the field carries an ObjectMarshaler.
	ObjectMarshalerType
	// BinaryType indicates that the field carries an opaque binary blob.
	BinaryType
	// BoolType indicates that the field carries a bool.
	BoolType
	// ByteStringType indicates that the field carries UTF-8 encoded bytes.
	ByteStringType
	// Complex128Type indicates that the field carries a complex128.
	Complex128Type
	// Complex64Type indicates that the field carries a complex128.
	Complex64Type
	// DurationType indicates that the field carries a time.Duration.
	DurationType
	// Float64Type indicates that the field carries a float64.
	Float64Type
	// Float32Type indicates that the field carries a float32.
	Float32Type
	// Int64Type indicates that the field carries an int64.
	Int64Type
	// Int32Type indicates that the field carries an int32.
	Int32Type
	// Int16Type indicates that the field carries an int16.
	Int16Type
	// Int8Type indicates that the field carries an int8.
	Int8Type
	// StringType indicates that the field carries a string.
	StringType
	// TimeType indicates that the field carries a time.Time that is
	// representable by a UnixNano() stored as an int64.
	TimeType
	// TimeFullType indicates that the field carries a time.Time stored as-is.
	TimeFullType
	// Uint64Type indicates that the field carries a uint64.
	Uint64Type
	// Uint32Type indicates that the field carries a uint32.
	Uint32Type
	// Uint16Type indicates that the field carries a uint16.
	Uint16Type
	// Uint8Type indicates that the field carries a uint8.
	Uint8Type
	// UintptrType indicates that the field carries a uintptr.
	UintptrType
	// ReflectType indicates that the field carries an interface{}, which should
	// be serialized using reflection.
	ReflectType
	// NamespaceType signals the beginning of an isolated namespace. All
	// subsequent fields should be added to the new namespace.
	NamespaceType
	// StringerType indicates that the field carries a fmt.Stringer.
	StringerType
	// ErrorType indicates that the field carries an error.
	ErrorType
	// SkipType indicates that the field is a no-op.
	SkipType

	// InlineMarshalerType indicates that the field carries an ObjectMarshaler
	// that should be inlined.
	InlineMarshalerType
)

func NewExtras() *Extras {
	return &Extras{fields: []Field{}}
}

func E() *Extras {
	return NewExtras()
}

type Extras struct {
	fields []Field
}

func (e *Extras) String(key string, value string) *Extras {
	if value == "" {
		// 字符串为空，不记录
		return e
	}
	e.fields = append(e.fields, Field{Key: key, Type: StringType, String: value})
	return e
}

func (e *Extras) Bytes(key string, value []byte) *Extras {
	if value == nil || len(value) == 0 {
		// 字符串为空，不记录
		return e
	}
	e.fields = append(e.fields, Field{Key: key, Type: ByteStringType, Bytes: value})
	return e
}

func (e *Extras) Bool(key string, v bool) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: BoolType, Boolean: v})
	return e
}

func (e *Extras) Int(key string, v int) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Int64Type, Integer: int64(v)})
	return e
}

func (e *Extras) Int8(key string, v int8) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Int8Type, Integer: int64(v)})
	return e
}

func (e *Extras) Int16(key string, v int16) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Int16Type, Integer: int64(v)})
	return e
}

func (e *Extras) Int32(key string, v int32) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Int32Type, Integer: int64(v)})
	return e
}

func (e *Extras) Int64(key string, v int64) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Int64Type, Integer: v})
	return e
}

func (e *Extras) Uint(key string, v uint) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Uint64Type, Integer: int64(v)})
	return e
}

func (e *Extras) Uint8(key string, v uint8) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Uint8Type, Integer: int64(v)})
	return e
}

func (e *Extras) Uint16(key string, v uint16) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Uint16Type, Integer: int64(v)})
	return e
}

func (e *Extras) Uint32(key string, v uint32) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Uint32Type, Integer: int64(v)})
	return e
}

func (e *Extras) Uint64(key string, v uint64) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Uint64Type, Integer: int64(v)})
	return e
}

func (e *Extras) Float64(key string, v float32) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Float64Type, Float: float64(v)})
	return e
}

func (e *Extras) Float32(key string, v float32) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: Float32Type, Float: float64(v)})
	return e
}

func (e *Extras) Error(err error) *Extras {
	e.fields = append(e.fields, Field{Key: "error", Type: ErrorType, Interface: err})
	return e
}

func (e *Extras) Fields() []Field {
	var result = make([]Field, len(e.fields))
	copy(result, e.fields)
	return result
}

func (e *Extras) NowMs() *Extras {
	e.fields = append(e.fields, Field{Type: Int64Type, Key: "ms", Integer: time.Now().UnixNano() / 1e6})
	return e
}

func (e *Extras) NamedError(key string, err error) *Extras {
	if err == nil {
		e.fields = append(e.fields, Field{Type: SkipType})
		return e
	}
	e.fields = append(e.fields, Field{Type: ErrorType, Key: key, Interface: err})

	return e
}

func (e *Extras) Any(key string, value interface{}) *Extras {
	e.fields = append(e.fields, Field{Key: key, Type: ReflectType, Interface: value})
	return e
}

func (e *Extras) Interface(key string, v interface{}) *Extras {
	if v == nil {
		e.fields = append(e.fields, Field{Type: StringType, Key: key, String: ""})
	} else {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice, reflect.Array, reflect.Struct, reflect.Map, reflect.Ptr:
			// 对结构体, map, slice 进行特判
			s, err1 := jsoniter.MarshalToString(v)
			if err1 != nil {
				e.fields = append(e.fields, Field{Type: ErrorType, Key: key, Interface: v})
			} else {
				e.fields = append(e.fields, Field{Type: StringType, Key: key, String: s})
			}
		default:
			e.fields = append(e.fields, Field{Type: ReflectType, Key: key, Interface: v})
		}
	}
	return e
}

// Ctx
//
//	@Description: 设置 context，并尝试从 context 提取 spanId，traceId
//	@receiver e
//	@param ctx  body any true "-"
//	@return *Extras
func (e *Extras) Ctx(ctx context.Context) *Extras {
	traceId, spanId := trace.ParseCurrentContext(ctx)
	e.Trace(traceId)
	e.Span(spanId)
	return e
}

func (e *Extras) Trace(traceId string) *Extras {
	if traceId != "" {
		e.fields = append(e.fields, Field{Key: "trace", Type: StringType, String: traceId})
	}
	return e
}

func (e *Extras) Span(spanId string) *Extras {
	if spanId != "" {
		e.fields = append(e.fields, Field{Key: "span", Type: StringType, String: spanId})
	}
	return e
}
