package storage

import (
	"context"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"reflect"
	"testing"
)

var (
	ctx       = context.TODO()
	tia       = &TXImageAudit{}
	qn1       = "http://cdn.itjsz.com/202305151626887.png"
	qn2       = "http://cdn.itjsz.com/202305151218810.png"
	awsImage1 = "https://s3.us-west-2.amazonaws.com/test.download.file.dev/20230511/202305151218810.png?response-content-disposition=inline&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEOT%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaDGV1LWNlbnRyYWwtMSJGMEQCIC1z%2BiJmcrd5YyW6oW3RPUhpXy5ltug1yE5873pVFkmQAiA27hVLxhAjjBrfr4ByMhEKCWWVxOeI0P4PcMQ2vD%2BwayqEAwj9%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDUzMzE1MDgxMTY1MCIMfbN2jAbKIiUB3pp2KtgCemTO67CtQucJaRa6IlLKV44pouTc7tbkZs4GaNL71LqyZvcy%2FCliVD0Rz9GZkbjMbhM1Lexeu%2Bnb%2Fn9m3Xo%2Bg7VVOvWLtY1fusH04ShmUFbYrsFE1cVqt4NHSUNAo0qeT%2BOm6jr6tkGVZzKhiAvmXcc183Bu8NcRddnfVK%2Fm0cAFE9a4bS8GKIx7DHp%2F9KHyrzIeKz9GAFF83ee99Xsa0lEel7SltZIp8XfQDfBkRuFsRmiXXUKEfUUy3aA%2FgzjnwAW3by3sSOlQc7S4Z6z6iHOChnoH%2BNSV7iXQtUp3wEd%2B8Pahtn%2FE7uMdyYfh5Fc2WFOjv%2FaL%2BwNMosrzsWH2yVPdoVC6pknCQI1dXfZOyZOxlrLwx08p5nYMbLGN45GlOqUfztTLXDNkdqMtHry6SUKMrU8kPInej7KZ%2BbBY3Lhf%2B2jawosQiurgHpVPKBSNrm93q%2B5C7ZEwkN%2BGowY6tAJkCvM0tDFfGxMu5CWC%2BfemqduSEM82%2Fsy4QBZw8nkLOAEjpSy%2Bry3M06JeOMrhINeEgrOkEeuJOWC5UvtVHV8NqEhzhoylZ%2BEPpRD5TgAyn4hNqdEH8zE2omM0pCdowNI128%2BkS%2BIq2%2FZSnMutI%2BQV%2F6mhbOz%2BO%2BT%2F1zCd5A0X6VBdlMXJlbK5OnXPTUQ0zQgu1g8M0xosiW3Lso1uHRD%2FuEyhinsjn7aRcu83IEi1xHiSJ5Ka8rqsSVLRmoR8n%2FKJdADT41ix2CnxibWWd6i7pzKgjH%2FxFI%2FBw5lH%2BFdI%2BmQrcSb12aUnnMLbpXuaK7fGZp%2FQXxM8IF5LKMRGFPcOi%2Bx4tHYBuW2r2q7j%2BdQh6mzooJqsjulqT69l71dPlcets2SohXumEe0Cy94Tu9MFv90neQ%3D%3D&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20230515T084312Z&X-Amz-SignedHeaders=host&X-Amz-Expires=300&X-Amz-Credential=ASIAXYISDYYBCSYIUH4V%2F20230515%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Signature=20c62c8820941ff62796eb0be390dbb10f271cc9bf94c822144579144327c462"
)

func SetUp() {
	var conf = config.ImageAudit{
		Bucket: "centersys-image-1302688347",
		Regin:  "na-siliconvalley",

		SecretId:  "xxx",
		SecretKey: "xxx",

		BizType: "266614fb500e910431f5ec44acaa25c5",
	}
	tia = NewTXImageAudit(&conf)
}

func TestTXImageAudit_AuditImage(t *testing.T) {
	SetUp()
	type fields struct {
		conf   config.ImageAudit
		client *cos.Client
	}
	type args struct {
		ctx    context.Context
		req    *ImageAuditReq
		appKey string
		dataId string
		url    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImageAuditResult
		wantErr bool
	}{
		{
			name: "1-qn", args: args{
				req: &ImageAuditReq{
					AppKey: "ccc",
					DataId: "1",
					Url:    qn1,
				},
			},
			want: &ImageAuditResult{
				AppKey: "ccc",
				DataId: "1",
				Url:    qn1,
			},
		},
		{
			name: "2-aws", args: args{
				req: &ImageAuditReq{
					AppKey: "cccxx",
					DataId: "2",
					Url:    awsImage1,
				},
			},
			want: &ImageAuditResult{
				AppKey: "cccxx",
				DataId: "2",
				Url:    awsImage1,
			},
		},
		{
			name: "2-qn", args: args{
				req: &ImageAuditReq{
					AppKey: "ccb",
					DataId: "3",
					Url:    qn2,
				},
			},
			want: &ImageAuditResult{
				AppKey: "ccb",
				DataId: "3",
				Url:    qn2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tia.AuditImage(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuditImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if tt.want.AppKey != tt.args.req.AppKey || tt.want.DataId != tt.args.req.DataId || tt.want.Url != tt.args.req.Url {
					t.Errorf("AuditImage() got = %+v, want %+v", got, tt.want)
				}
			}
			t.Logf("re %+v", *got)
		})
	}
}

func TestTXImageAudit_AuditImages(t *testing.T) {
	SetUp()
	type fields struct {
		conf   config.ImageAudit
		client *cos.Client
	}
	type args struct {
		ctx    context.Context
		appKey string
		urls   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*ImageAuditResult
		wantErr bool
	}{
		{
			name: "1-qn", args: args{
				appKey: "cccc",
				urls:   map[string]string{"qn111": qn1, "qn222": qn2},
			},
		},
		{
			name: "2-aws", args: args{
				appKey: "cccc",
				urls:   map[string]string{"aws11": awsImage1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tia.AuditImages(ctx, tt.args.urls, tt.args.appKey, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("AuditImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, g := range got {

				for d, u := range tt.args.urls {
					if d == g.DataId {
						if g.AppKey != tt.args.appKey || u != g.Url {
							t.Errorf("AuditImages() got = %+v, want %+v", tt.args, g)
						}
						t.Logf("%+v", *g)
					}
				}
			}
		})
	}
}
