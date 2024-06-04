package struct_tool

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// 测试一个结构体转map函数
func TestStructToMapByTag(t *testing.T) {
	type AStruct struct {
		S string `json:"s"`
	}
	type TestStruct struct {
		Er   error   `json:"er" gorm:"er"`
		S    string  `json:"s"`
		I    int     `json:"i"`
		F    float32 `json:"f"`
		Bi   byte    `json:"bi,byte"`
		B    bool    `json:"b" gorm:"column:b"`
		Bits []byte  `gorm:"column:bits" json:"bits"`
		R    rune    `json:"r" gorm:"r"`

		AS AStruct `json:"as" gorm:"as"`
	}

	type args struct {
		data          interface{}
		tag           string
		ignoreFields  []string
		notIgnoreZero bool
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "1", args: args{
				data: &TestStruct{S: "1", I: 12, F: 1.2, Bi: 1, B: true, Bits: []byte{1, 2, 3}},
				tag:  "json", ignoreFields: []string{"i"}, notIgnoreZero: false,
			},
			want: map[string]interface{}{"s": "1" /*"i": 12,*/, "f": 1.2, "bi": 1, "b": true, "bits": []byte{1, 2, 3}},
		},
		{
			name: "1-2", args: args{
				data: &TestStruct{S: "1", I: 12, F: 1.2, Bi: 1, B: true, Bits: []byte{}, AS: AStruct{}},
				tag:  "json", ignoreFields: []string{}, notIgnoreZero: false,
			},
			want: map[string]interface{}{"s": "1", "i": 12, "f": 1.2, "bi": 1, "b": true},
		},
		{
			name: "2", args: args{
				data: TestStruct{S: "1", I: 12, Bi: 1, B: true, Bits: []byte{1, 2, 3}, R: 'a', AS: AStruct{}},
				tag:  "json", ignoreFields: []string{"i"}, notIgnoreZero: false,
			},
			want: map[string]interface{}{"s": "1" /*"i": 12,*/, "bi": 1, "b": true, "bits": []byte{1, 2, 3}, "r": 'a'},
		},
		{
			name: "2-2", args: args{
				data: TestStruct{S: "1", I: 12, Bi: 1, B: true, Bits: []byte{1, 2, 3}, R: 'a', AS: AStruct{"0"}},
				tag:  "gorm", notIgnoreZero: true,
			},
			want: map[string]interface{}{"as": AStruct{S: "0"}, "column:b": true, "column:bits": []byte{1, 2, 3}, "er": nil, "r": 'a'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StructToMapByTag(tt.args.data, tt.args.tag, tt.args.ignoreFields, tt.args.notIgnoreZero); !reflect.DeepEqual(got, tt.want) {
				gotJ, _ := json.Marshal(got)
				wantJ, _ := json.Marshal(tt.want)
				if 0 != strings.Compare(string(gotJ), string(wantJ)) {
					t.Errorf("StructToMapByTag() = %v, \n\twant %v", got, tt.want)
				}
			}
		})
	}
}
