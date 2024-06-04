package trace

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/spf13/cast"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
)

type SuTracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func (s *SuTracer) Close() {
	s.closer.Close()
}

var Tracer = &SuTracer{}

/*initJaeger
* @Description: 初始化jaeger tracer
* @param cfg
 */
func initJaeger(cfg config.JaegerConf) {
	if !cfg.JaegerOn {
		return
	}
	var err error
	var tags []opentracing.Tag
	for k, v := range cfg.Tags {
		tags = append(tags, opentracing.Tag{Key: k, Value: v})
	}

	if cfg.CollectorEndpoint == "" {
		fmt.Println("jaeger collector endpoint is empty")
		return
	}

	sender, err := jaeger.NewUDPTransport(cfg.CollectorEndpoint, 0)
	if err != nil {
		fmt.Println("jaeger new udp transport error", err)
		return
	}
	/*
		const：全量采集，采样率设置0,1 分别对应打开和关闭
		probabilistic：概率采集，默认万份之一，取值可在 0 至 1 之间，例如设置为 0.5 的话意为只对 50% 的请求采样
		rateLimiting：限速采集，每秒只能采集一定量的数据，如设置2的话，就是每秒采集2个链路数据
		remote ：是遵循远程设置，取值的含义和 probabilistic 相同，都意为采样的概率，只不过设置为 remote 后，Client 会从 Jaeger Agent 中动态获取采样率设置。
		guaranteedThroughput:复合采样，至少每秒采样lowerBound次(rateLimiting),超过lowerBound次的话，按照samplingRate概率来采样(probabilistic)
	*/
	var sampler jaeger.Sampler
	if cfg.RateLimitPerSecond > 0 {
		sampler = jaeger.NewRateLimitingSampler(cfg.RateLimitPerSecond)
	} else {
		if cfg.SamplerFreq >= 1 {
			// 始终采集
			sampler = jaeger.NewConstSampler(true)
		} else {
			// 0.0 and 1.0
			if cfg.SamplerFreq <= 0 {
				cfg.SamplerFreq = 0.0001
			}
			sampler, err = jaeger.NewProbabilisticSampler(cfg.SamplerFreq)
			if err != nil {
				panic(err)
			}
		}
	}

	// ZipkinSharedRPCSpan
	fmt.Println("service name", cfg.AppName)
	t, c := jaeger.NewTracer(cfg.AppName,
		sampler,
		jaeger.NewRemoteReporter(sender),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)

	opentracing.SetGlobalTracer(t)

	Tracer = &SuTracer{
		tracer: t,
		closer: c,
	}

	return
}

func StartSpan(ctx context.Context, tag map[string]interface{}, carrier opentracing.TextMapCarrier) opentracing.Span {
	s := opentracing.SpanFromContext(ctx)
	for k, v := range tag {
		s.SetTag(k, v)
	}
	for k, v := range carrier {
		s.SetBaggageItem(k, v)
	}

	return s
}

// 这个方法废弃，后面传traceId，spanId都用 string 类型
// deprecated
func NewSpanCtx(ctx context.Context, opName string, parentId uint64, tags map[string]interface{}, baggage map[string]string) (opentracing.Span, context.Context) {
	tId, sId := ParseFromContext(ctx)
	traceId := cast.ToUint64(tId)
	spanId := cast.ToUint64(sId)
	opts := []opentracing.StartSpanOption{
		opentracing.SpanReference{
			Type: 99,
			ReferencedContext: jaeger.NewSpanContext(jaeger.TraceID{Low: traceId}, jaeger.SpanID(spanId),
				jaeger.SpanID(parentId), false, baggage),
		},
	}

	if parentId != 0 {
		parent := jaeger.NewSpanContext(jaeger.TraceID{Low: traceId}, jaeger.SpanID(parentId), jaeger.SpanID(parentId), false, baggage)
		opts = append(opts, opentracing.ChildOf(parent))
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, opName, opts...)
	span.SetTag("traceId", traceId)
	span.SetTag("spanId", spanId)
	span.SetTag("parentId", parentId)

	if len(tags) > 0 {
		for k, v := range tags {
			span.SetTag(k, v)
		}
	}

	return span, ctx
}

func NewJaegerSpan(opName string, traceId, spanId, parentId string, tags map[string]interface{}, baggage map[string]string) opentracing.Span {
	var (
		err            error
		jaegerTraceId  jaeger.TraceID
		jaegerSpanId   jaeger.SpanID
		jaegerParentId jaeger.SpanID
	)
	if jaegerTraceId, err = jaeger.TraceIDFromString(traceId); err != nil {
		return opentracing.StartSpan(opName)
	}
	if jaegerSpanId, err = jaeger.SpanIDFromString(spanId); err != nil {
		log.Printf("jaeger.SpanIDFromString(spanId) error %v", err)
	}
	jaegerParentId, _ = jaeger.SpanIDFromString(parentId) // parentId 可能为空， 是允许的
	opts := []opentracing.StartSpanOption{
		opentracing.SpanReference{
			Type:              99,
			ReferencedContext: jaeger.NewSpanContext(jaegerTraceId, jaegerSpanId, jaegerParentId, false, baggage),
		},
	}

	if jaegerParentId != 0 {
		parent := jaeger.NewSpanContext(jaegerTraceId, jaegerParentId, jaegerParentId, false, baggage)
		opts = append(opts, opentracing.ChildOf(parent))
	}
	span := opentracing.StartSpan(opName, opts...)
	span.SetTag("traceId", traceId)
	span.SetTag("spanId", spanId)
	span.SetTag("parentId", parentId)

	if len(tags) > 0 {
		for k, v := range tags {
			span.SetTag(k, v)
		}
	}

	return span
}

// GrpcUnaryJaegerInterceptor
//
//	@Description: grpc client interceptor
//	@return unc
func GrpcUnaryJaegerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {

	nextCtx, reqId, spanId, nextSpanId := ParseOrGenContext(ctx)
	span := NewJaegerSpan("grpc_c:"+method, reqId, nextSpanId, spanId, nil, nil)
	defer span.Finish()

	return invoker(nextCtx, method, req, reply, cc, opts...)
}

// GrpcStreamJaegerInterceptor
//
//	@Description: grpc client stream interceptor
//	@param ctx  body any true "-"
//	@param desc  body any true "-"
//	@param cc  body any true "-"
//	@param method  body any true "-"
//	@param streamer  body any true "-"
//	@param opts  body any true "-"
//	@return grpc.ClientStream
//	@return error
func GrpcStreamJaegerInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	nextCtx, reqId, spanId, nextSpanId := ParseOrGenContext(ctx)
	span := NewJaegerSpan("grpc_c:"+method, reqId, nextSpanId, spanId, nil, nil)
	defer span.Finish()

	return streamer(nextCtx, desc, cc, method, opts...)
}
