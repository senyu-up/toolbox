package main

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/tool/http/req"
	"github.com/senyu-up/toolbox/tool/trace"
)

func ExampleReq() {
	var ctx = trace.NewTrace()
	fmt.Printf("send req with traceId: %+v \n", ctx)
	resp, err := req.New(context.TODO()).Context(ctx). // 带trace的 ctx 要通过 Context(ctx) 设置！
								Get("http://127.0.0.1:8181/trace")
	if err != nil {
		fmt.Printf("req.Get error: %v \n", err)
		return
	} else {
		fmt.Printf("resp: %s \n", resp)
	}

	resp, err = req.New(ctx).Get("http://127.0.0.1:8181/system/health")
	if err != nil {
		fmt.Printf("req.Get error: %v \n", err)
		return
	} else {
		fmt.Printf("resp: %s \n", resp)
	}
}

func main() {
	ExampleReq()
}
