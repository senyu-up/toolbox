package su_logger

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func BenchmarkError(b *testing.B) {
	var ctx = context.TODO()
	var err = errors.New("This is Error ")
	type args struct {
		ctx   context.Context
		err   error
		msg   string
		extra []*Extra
	}
	tests := []struct {
		name string
		args args
	}{
		{
			// BenchmarkError/1-8 	1000000000	         0.0000101 ns/op
			name: "1", args: args{ctx: ctx, err: err, msg: "Got Err !", extra: nil},
		},
		{
			// BenchmarkError/2-8 	1000000000	         0.0000193 ns/op
			name: "2", args: args{ctx: ctx, err: err, msg: fmt.Sprintf("Got Err ! %v", err), extra: nil},
		},
		{
			// BenchmarkError/3-8 	1000000000	         0.0000146 ns/op
			name: "3", args: args{ctx: ctx, err: err, msg: fmt.Sprintf("Got Err ! %v", err),
				extra: []*Extra{E().String("name", "xh").Int("age", 18)}},
		},
		{
			// BenchmarkError/4-8 	1000000000	         0.0000145 ns/op
			name: "4", args: args{ctx: ctx, err: err,
				msg:   fmt.Sprintf("Got Req %v Err ! %v ", map[string]string{"haha": "cc"}, err),
				extra: []*Extra{E().String("name", "xh").Any("info", map[string]string{"haha1": "a"})}},
		},
	}

	b.ResetTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(t *testing.B) {
			Error(tt.args.ctx, tt.args.err, tt.args.msg, tt.args.extra...)
		})
	}
	b.StopTimer()
}
