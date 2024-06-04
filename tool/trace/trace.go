package trace

import (
	"context"
	"github.com/rs/xid"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/geoip"
	"github.com/uber/jaeger-client-go/utils"
	"google.golang.org/grpc/metadata"
)

var (
	localIp net.IP
	seq     uint32
	spanSeq uint32
)

// init
//
//	@Description: 初始化, 本机ip，随机发数器
func init() {
	var err error
	localIp, err = geoip.GetLocalIPV4Obj()
	if err != nil {
		panic(err)
	}
	seq = rand.Uint32()
	spanSeq = rand.Uint32()
}

func Init(cfg *config.TraceConfig) {
	initJaeger(cfg.Jaeger)
}

// GetRequestId
// @description 获取request id, 优先从context中获取, 否则从当前协程获取
func GetRequestId(ctx context.Context) (reqId string) {
	if ctx != nil {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			reqData := md.Get(enum.RequestId)
			if len(reqData) > 0 {
				return reqData[0]
			}
		} else if md, ok = metadata.FromOutgoingContext(ctx); ok {
			reqData := md.Get(enum.RequestId)
			if len(reqData) > 0 {
				return reqData[0]
			}
		} else {
			v := ctx.Value(enum.RequestId)
			if v != nil {
				return v.(string)
			}
		}
	}

	return NewTraceID()
}

// ParseRequestId
// @description 获取request id, 优先从context中获取, 如果没有不会生成
func ParseRequestId(ctx context.Context) (reqId string) {
	if ctx != nil {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			reqData := md.Get(enum.RequestId)
			if len(reqData) > 0 {
				return reqData[0]
			}
		}

		v := ctx.Value(enum.RequestId)
		if v != nil {
			return v.(string)
		}

		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			reqData := md.Get(enum.RequestId)
			if len(reqData) > 0 {
				return reqData[0]
			}
		}
	}
	return ""
}

// ParseOrGenContext
//
//	@Description: 服务内 链路追踪 ctx 处理逻辑
//	@param ctx  body any true "-"
//	@return context.Context
//	@return string
//	@return string
//	@return string
func ParseOrGenContext(ctx context.Context) (context.Context, string, string, string) {
	traceId, spanId := ParseCurrentContext(ctx)
	if traceId != "" && spanId != "" {
		// traceId, spanId 都有，则不生成了
		// 但是要把新 spanId 设置到 out metadata中
		//ctx = metadata.AppendToOutgoingContext(ctx, enum.SpanId, newSpanId)
		var nextSpanId = NewSpanID()
		// 设置下一级 spanId
		ctx = NewWithOutCtx(ctx, traceId, spanId, nextSpanId)
		return ctx, traceId, spanId, nextSpanId
	}
	return NewTraceWithId(ctx, traceId, spanId, "")
}

// ParseOrGenGrpcClientContext
//
//	@Description: grpc client 端从ctx中获取traceId和spanId, 如果不存在则生成, 同时生成新spanId
//	@param ctx  body any true "-"
//	@return context.Context
//	@return string
//	@return string
//	@return string
func ParseOrGenGrpcClientContext(ctx context.Context) (context.Context, string, string, string) {
	traceId, parentId, currentSpanId := ParseGrpcInContextTraceSpans(ctx)
	if traceId != "" && parentId != "" {
		// traceId, spanId 都有，则不生成了
		// 但是要把新 spanId 设置到 out metadata中
		if currentSpanId == "" {
			currentSpanId = NewSpanID()
		}
		var nextSpanId = NewSpanID()
		nextCtx := NewWithOutCtx(ctx, traceId, currentSpanId, nextSpanId)
		return nextCtx, traceId, parentId, currentSpanId
	}
	return NewTraceWithId(ctx, traceId, parentId, "")
}

// ParseOrGenGrpcContext
//
//	@Description: grpc server端从ctx中获取traceId和spanId, 如果不存在则生成, 同时生成新spanId 设置到 out metadata中
//	@param ctx  body any true "-"
//	@return context.Context
//	@return string
//	@return string
//	@return string
func ParseOrGenGrpcContext(ctx context.Context) (context.Context, string, string, string) {
	traceId, parentId, currentSpanId := ParseGrpcInContextTraceSpans(ctx)
	if traceId != "" && parentId != "" {
		// traceId, spanId 都有，则不生成了
		// 但是要把新 spanId 设置到 out metadata中
		if currentSpanId == "" {
			currentSpanId = NewSpanID()
		}
		nextCtx := NewNextCtx(ctx, traceId, parentId, currentSpanId)
		return nextCtx, traceId, parentId, currentSpanId
	}
	return NewTraceWithId(ctx, traceId, parentId, "")
}

// ParseFromContext
// @description 基于context获取链路id以及spanId
func ParseFromContext(ctx context.Context) (traceId string, spanId string) {
	defer func() {
		_ = recover()
	}()
	if ctx == nil {
		return
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqData := md.Get(enum.RequestId)
		if len(reqData) > 0 {
			traceId = reqData[0]
		}
		reqData = md.Get(enum.SpanId)
		if len(reqData) > 0 {
			spanId = reqData[0]
		}

		if traceId != "" {
			return
		}
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		reqData := md.Get(enum.RequestId)
		if len(reqData) > 0 {
			traceId = reqData[0]
		}
		reqData = md.Get(enum.SpanId)
		if len(reqData) > 0 {
			spanId = reqData[0]
		}

		if traceId != "" {
			return
		}
	}

	v := ctx.Value(enum.RequestId)
	if v != nil {
		traceId = v.(string)
	}
	v = ctx.Value(enum.SpanId)
	if v != nil {
		spanId = v.(string)
	}
	return
}

// ParseCurrentContext
// @description 解析当前登记的 context的 traceId，spanId
func ParseCurrentContext(ctx context.Context) (traceId string, spanId string) {
	if ctx == nil {
		return
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		reqData := md.Get(enum.RequestId)
		if len(reqData) > 0 {
			traceId = reqData[0]
		}
		reqData = md.Get(enum.SpanId)
		if len(reqData) > 0 {
			spanId = reqData[0]
		}
	}

	if traceId == "" {
		v := ctx.Value(enum.RequestId)
		if v != nil {
			traceId = v.(string)
		}
	}

	if spanId == "" {
		v := ctx.Value(enum.SpanId)
		if v != nil {
			spanId = v.(string)
		}
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if traceId == "" {
			reqData := md.Get(enum.RequestId)
			if len(reqData) > 0 {
				traceId = reqData[0]
			}
		}
		if spanId == "" {
			reqData := md.Get(enum.SpanId)
			if len(reqData) > 0 {
				spanId = reqData[0]
			}
		}
	}
	return
}

// ParseAndGenContext
// @description 从 ctx 获取 traceId 和 spanId, nextSpanId， 如果没有就没有
func ParseGrpcInContextTraceSpans(ctx context.Context) (traceId string, spanId string, nextSpanId string) {
	defer func() {
		_ = recover()
	}()
	if ctx == nil {
		return
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqData := md.Get(enum.RequestId)
		if len(reqData) > 0 {
			traceId = reqData[0]
		}
		reqData = md.Get(enum.SpanId)
		if len(reqData) > 0 {
			spanId = reqData[0]
		}
	}

	if traceId == "" {
		v := ctx.Value(enum.RequestId)
		if v != nil {
			traceId = v.(string)
		}
	}
	if spanId == "" {
		v := ctx.Value(enum.SpanId)
		if v != nil {
			spanId = v.(string)
		}
	}

	return
}

// NewContextWithRequestIdAndSpanId
// @description 基于context创建新的context, 并且添加request id和span id
func NewContextWithRequestIdAndSpanId(ctx context.Context, reqId, spanId string) context.Context {
	return AssignToContext(ctx, enum.RequestId, reqId, enum.SpanId, spanId)
}

// NewContextWithRequestId
// @description 基于context创建新的context, 并且添加request id
func NewContextWithRequestId(ctx context.Context, reqId string) context.Context {
	return AssignToContext(ctx, enum.RequestId, reqId)
}

// NewContextWithSpanId
// @description 基于context创建新的context, 并且添加span id
func NewContextWithSpanId(ctx context.Context, spanId string) context.Context {
	return AssignToContext(ctx, enum.SpanId, spanId)
}

// NewTrace
// @description new trace
func NewTrace() context.Context {
	reqId := NewTraceID()
	spanId := NewSpanID()

	return NewContextWithRequestIdAndSpanId(context.Background(), reqId, spanId)
}

// NewTraceWithId
//
//	@Description: 基于context创建新的context, 并且添加request id和span id
//	@param ctx  body any true "-"
//	@param reqIdStr  body any true "-"
//	@param spanIdStr  body any true "-"
//	@param nextSpanIdStr  body any true "-"
//	@return context.Context
//	@return string
//	@return string
//	@return string
func NewTraceWithId(ctx context.Context, reqIdStr, spanIdStr, nextSpanIdStr string) (context.Context, string, string, string) {
	if nil == ctx {
		ctx = context.Background()
	}
	// 空检查
	if reqIdStr == "" {
		reqIdStr = NewTraceID()
	}
	if spanIdStr == "" {
		spanIdStr = NewSpanID()
	}
	if nextSpanIdStr == "" {
		nextSpanIdStr = NewSpanID()
	}

	// 一般 ctx value
	context.WithValue(ctx, enum.RequestId, reqIdStr)
	context.WithValue(ctx, enum.SpanId, spanIdStr)

	// incoming ctx value
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Set(enum.RequestId, reqIdStr)
		md.Set(enum.SpanId, spanIdStr)
		ctx = metadata.NewIncomingContext(ctx, md)
	} else {
		var md = metadata.MD{}
		md.Set(enum.RequestId, reqIdStr)
		md.Set(enum.SpanId, spanIdStr)
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	// outgoing ctx value
	ctx = metadata.AppendToOutgoingContext(ctx, enum.RequestId, reqIdStr, enum.SpanId, nextSpanIdStr)
	return ctx, reqIdStr, spanIdStr, nextSpanIdStr
}

/*GetFromContext
* @Description: 从ctx中获取指定的key, 优先从普通ctx, 再从metadata中获取
* @param ctx
* @param k
* @return string
 */
func GetFromContext(ctx context.Context, k string) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(k)
	if v != nil {
		return v.(string)
	} else if md, ok := metadata.FromIncomingContext(ctx); ok {
		rs := md.Get(k)
		if len(rs) > 0 {
			return rs[0]
		}
	} else if md, ok = metadata.FromOutgoingContext(ctx); ok {
		rs := md.Get(k)
		if len(rs) > 0 {
			return rs[0]
		}
	}

	return ""
}

// AssignToContext
// @description 赋值到context中, 支持 grpc 以及 常规的context, 使用追加模式, 如果有冲突, 新值覆盖旧值
func AssignToContext(ctx context.Context, kv ...string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for i := 0; i < len(kv); i += 2 {
			md.Set(kv[i], kv[i+1])
		}
		return metadata.NewIncomingContext(ctx, md)
	} else if md, ok = metadata.FromOutgoingContext(ctx); ok {
		for i := 0; i < len(kv); i += 2 {
			md.Set(kv[i], kv[i+1])
		}
		return metadata.NewOutgoingContext(ctx, md)
	} else {
		l := len(kv)
		for i := 0; i < l; i += 2 {
			ctx = context.WithValue(ctx, kv[i], kv[i+1])
		}

		return ctx
	}
}

func IdGeneratorFactory() func() uint64 {
	seedGenerator := utils.NewRand(time.Now().UnixNano())
	pool := sync.Pool{
		New: func() interface{} {
			return rand.NewSource(seedGenerator.Int63())
		},
	}

	return func() uint64 {
		generator := pool.Get().(rand.Source)
		number := uint64(generator.Int63())
		pool.Put(generator)
		return number
	}
}

var idGenerator = IdGeneratorFactory()

// NewTraceID creates and returns a trace ID.
func NewTraceID() (traceID string) {
	return xid.New().String()
}

// NewSpanID creates and returns a span ID.
func NewSpanID() (spanID string) {

	return xid.New().String()
}

// WithContext
// @description 将链路id填充到context
func WithContext(ctx context.Context) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqData := md.Get(enum.RequestId)
		if len(reqData) > 0 {
			return
		}
		md.Set(enum.RequestId, NewTraceID())
		// 填充request id
		el := reflect.ValueOf(ctx).Elem().FieldByName("val")
		v := reflect.NewAt(el.Type(), unsafe.Pointer(el.UnsafeAddr())).Elem()
		v.Set(reflect.ValueOf(md))
	}
}

// 解析传入的 ctx，解析出 trace 信息，注入到新ctx中
// 用于取消 ctx 的cancel，和提出其他value信息。
func ReNewCtxWithTrace(ctx context.Context) (newCtx context.Context) {
	newCtx = context.Background()
	if ctx == nil {
		return
	}
	// 检查当前ctx是否携带链路id
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		newCtx = metadata.NewIncomingContext(newCtx, md)
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		newCtx = metadata.NewOutgoingContext(newCtx, md)
	}

	v := ctx.Value(enum.RequestId)
	if v != nil {
		newCtx = context.WithValue(newCtx, enum.RequestId, v)
	}
	v = ctx.Value(enum.SpanId)
	if v != nil {
		newCtx = context.WithValue(newCtx, enum.SpanId, v)
	}
	return
}

func NewNextCtx(ctx context.Context, traceId string, spanId string, nextSpanId string) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}

	// 检查当前ctx是否携带链路id
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Set(enum.RequestId, traceId)
		md.Set(enum.SpanId, spanId)
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		md.Set(enum.RequestId, traceId)
		md.Set(enum.SpanId, nextSpanId)
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(enum.RequestId, traceId, enum.SpanId, nextSpanId))
	}

	return context.WithValue(context.WithValue(ctx, enum.RequestId, traceId), enum.SpanId, spanId)
}

func NewWithOutCtx(ctx context.Context, traceId string, spanId string, nextSpanId string) context.Context {
	if ctx == nil {
		ctx = NewTrace()
		return ctx
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		md.Set(enum.RequestId, traceId)
		md.Set(enum.SpanId, nextSpanId)
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(enum.RequestId, traceId, enum.SpanId, nextSpanId))
	}

	// 检查当前ctx是否携带链路id
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Set(enum.RequestId, traceId)
		md.Set(enum.SpanId, spanId)
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	return context.WithValue(context.WithValue(ctx, enum.RequestId, traceId), enum.SpanId, spanId)
}
