package main

import (
	"context"
	"fmt"
	config2 "github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/trace"
	"github.com/spf13/cast"
	"time"
)

func Example1() {
	originCtx := trace.NewTrace()
	ctx, cancel := context.WithCancel(originCtx)
	newCtx1, tId, sp, _ := trace.ParseOrGenGrpcContext(ctx)
	cancel()
	if err := newCtx1.Err(); err != nil {
		fmt.Printf("newCtx1 err: %v\n", err)
	}
	fmt.Printf("newCtx1 traceId %s, id1 :%s, sp1: %s\n", trace.GetRequestId(originCtx), tId, sp)

	ctx, cancel = context.WithCancel(originCtx)
	var newCtx2 = trace.ReNewCtxWithTrace(ctx)
	newCtx3, tId2, sp2, _ := trace.ParseOrGenGrpcContext(newCtx2)
	cancel()
	if err := newCtx3.Err(); err != nil {
		fmt.Printf("newCtx3 err: %v\n", err)
	}
	fmt.Printf("newCtx3 traceId %s, id2 :%s, sp2: %s\n", trace.GetRequestId(originCtx), tId2, sp2)
}

func main1() {
	// 初始化配置
	var config = &config2.TraceConfig{
		ServerLogOn:    true,
		ClientLogOn:    true,
		ClientLogLevel: "info",

		Jaeger: config2.JaegerConf{
			JaegerOn: true,           // 开启状态才会初始化 jaeger
			AppName:  "example_test", // 当前记录到jeager时，上报的 app name。

			// jeager 采集器地址信息
			CollectorEndpoint: "127.0.0.1:6831", // 采集器地址 UDP 采集！
			//CollectorEndpoint: "http://localhost:14268/api/traces",
			//AgentPort:         "14268",
			User:     "",
			Password: "",

			SamplerFreq:              1,
			QueueFlushIntervalSecond: 1,
		},
	}
	trace.Init(config) // 通过配置初始化 trace，jaeger 客户端

	var ctx = trace.NewTrace()                                                        // 直接生成带 trace 的 context
	trace.NewContextWithRequestIdAndSpanId(ctx, "xxx_xxx_trace_id", "xxx_xxx_spanId") // 在给定 context里设置 requestId 和 spanId

	// 从 context 中获取 traceId 和 spanId
	traceId, spanId := trace.ParseFromContext(ctx)
	fmt.Printf("traceId: %s, spanId: %s", traceId, spanId)

	// 生成新的 traceId、spanId
	spanId = trace.NewSpanID()
	traceId = trace.NewTraceID()

	for {
		ctx = trace.NewTrace() // 生成一个新的 trace

		span, ctx := trace.NewSpanCtx(ctx, "my_test2", cast.ToUint64(spanId), nil, nil)
		span.Finish()

		traceId, spanId = trace.ParseFromContext(ctx)

		// 记录一条 子span
		span = trace.NewJaegerSpan("my_test", traceId, trace.NewSpanID(), spanId, nil, nil)
		span.Finish()

		fmt.Print("finished~\n")
		time.Sleep(time.Second)
	}
}

func main() {
	main1()
	//Example1()
}
