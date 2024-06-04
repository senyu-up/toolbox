package runtime

import (
	"fmt"
	"testing"
)

type MyErr struct {
	Msg string
}

func (MyErr) Error() string {
	return "my error"
}

func TestIsZero(t *testing.T) {
	type exampleStruct struct {
		Field1 int
		Field2 string
	}

	type args struct {
		x interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil value",
			args: args{x: nil},
			want: true,
		},
		{
			name: "zero value of int",
			args: args{x: 0},
			want: true,
		},
		{
			name: "non-zero value of int",
			args: args{x: 1},
			want: false,
		},
		{
			name: "zero value of float",
			args: args{x: 0.0},
			want: true,
		},
		{
			name: "non-zero value of float",
			args: args{x: 3.14},
			want: false,
		},
		{
			name: "empty string",
			args: args{x: ""},
			want: true,
		},
		{
			name: "non-empty string",
			args: args{x: "hello"},
			want: false,
		},
		{
			name: "zero value of bool",
			args: args{x: false},
			want: true,
		},
		{
			name: "non-zero value of bool",
			args: args{x: true},
			want: false,
		},
		{
			name: "zero value of pointer",
			args: args{x: (*int)(nil)},
			want: true,
		},
		{
			name: "non-zero value of pointer",
			args: args{x: new(int)},
			want: false,
		},
		{
			name: "zero value of error",
			args: args{x: (error)(nil)},
			want: true,
		},
		{
			name: "non-zero value of error",
			args: args{x: fmt.Errorf("an error occurred")},
			want: false,
		},
		{
			name: "zero value of map",
			args: args{x: make(map[string]int)},
			want: true,
		},
		{
			name: "non-zero value of map",
			args: args{x: map[string]int{"a": 1}},
			want: false,
		},
		{
			name: "zero value of slice",
			args: args{x: make([]int, 0)},
			want: true,
		},
		{
			name: "non-zero value of slice",
			args: args{x: []int{1}},
			want: false,
		},
		{
			name: "zero value of array",
			args: args{x: [3]int{}},
			want: true,
		},
		{
			name: "non-zero value of array",
			args: args{x: [3]int{1, 2, 3}},
			want: false,
		},
		{
			name: "zero value of struct",
			args: args{x: exampleStruct{}},
			want: true,
		},
		{
			name: "non-zero value of struct",
			args: args{x: exampleStruct{Field1: 1, Field2: "hello"}},
			want: false,
		},

		{
			name: "zero value of err struct",
			args: args{x: MyErr{}},
			want: true,
		},
		{
			name: "non-zero value err struct",
			args: args{x: MyErr{Msg: "hello"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsZero(tt.args.x); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
