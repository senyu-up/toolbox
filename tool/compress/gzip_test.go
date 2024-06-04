package compress

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
)

func TestGzip(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "1", args: args{
				data: []byte("23498refwskdafakj3wqrwiq874302923140985132452318534ujidfsavasfks"),
			},
			wantErr: false,
			want:    []byte("H4sIAAAAAAAEAA3HSQrAIAwAwHuhf9EsNHlOwAY0Jw2t32/nNoCksm7fGc3cYuCea/cpF2EBBaxUVLgiEP8RRnpGb572WnrkeXy7ThdsQgAAAA=="),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gzip(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Gzip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var want = make([]byte, len(tt.want))

			n, err := base64.StdEncoding.Decode(want, tt.want)
			if err != nil {
				fmt.Printf("Error decoding string: %s ", err.Error())
				return
			} else {
				want = want[:n]
			}
			if !reflect.DeepEqual(got, want) {
				if 0 != bytes.Compare(got, want) {
					t.Errorf("Gzip() got = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestGunzip(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "1", args: args{
				content: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 4, 192, 81, 10, 4, 33, 8, 6, 224, 43, 149, 191, 177, 122, 28, 161, 21, 202, 167, 146, 25, 175, 63, 31, 129, 85, 238, 223, 43, 99, 154, 91, 108, 212, 185, 181, 142, 252, 24, 141, 148, 208, 185, 169, 140, 14, 226, 65, 232, 50, 192, 207, 94, 211, 211, 94, 75, 143, 252, 2, 0, 0, 255, 255, 185, 96, 190, 191, 64, 0, 0, 0},
			},
			wantErr: false,
			want:    []byte("23498refwskdafakj3wqrwiq874302923140985132452318534ujidfsavasfks"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gunzip(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Gunzip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gunzip() got = %s, want %v", got, tt.want)
			}
		})
	}
}
