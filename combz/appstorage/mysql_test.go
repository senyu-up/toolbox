package appstorage

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"sync"
	"testing"
)

var (
	appkey = "VUpqSE1YTWJhYToxNjc1ODQ2NDIzOmRldmVsb3A="
)

func TestDBStorage_GetDB(t *testing.T) {
	SetUp()
	type fields struct {
		ins            sync.Map
		readIns        sync.Map
		channelID      string
		dsnMap         sync.Map
		wg             sync.WaitGroup
		db             *gorm.DB
		instInitialing sync.Map
		redisCli       *redis.UniversalClient
	}
	type args struct {
		app string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *gorm.DB
	}{
		{
			name: "1", args: args{
				app: appkey,
			},
			want: &gorm.DB{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := appDb.GetDB(tt.args.app); nil == got {
				t.Errorf("GetDB() = %v, want %v", got, tt.want)
			} else {
				if rows, err := got.Table("account_info").Select("*").Where("1=1").Rows(); err != nil {
					t.Errorf("query err")
				} else {
					var cols, _ = rows.Columns()
					t.Logf("get table cols %v", cols)
					for rows.Next() {
						var id int
						var name string
						if err = rows.Scan(&id, &name); err != nil {
							t.Errorf("mysql scan err %v", err)
						}
						t.Logf("id %d app_key %s", id, name)
					}
				}
			}
		})
	}
}

func TestDBStorage_SetDB(t *testing.T) {
	SetUp()
	type fields struct {
		ins            sync.Map
		readIns        sync.Map
		channelID      string
		dsnMap         sync.Map
		wg             sync.WaitGroup
		db             *gorm.DB
		instInitialing sync.Map
		redisCli       *redis.UniversalClient
	}
	type args struct {
		app string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{app: "Z1FLdkg2djlhYToxNjc1OTM1MDM1OmRldmVsb3A="},
			wantErr: false,
		},
		{
			name: "2", args: args{app: "1FLdkg2djlhYToxNjc1OTM1MDM1OmRldmVsb3A="},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := appDb.SetDB(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("SetDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBStorage_RemoveDB(t *testing.T) {
	SetUp()
	type fields struct {
		ins            sync.Map
		readIns        sync.Map
		channelID      string
		dsnMap         sync.Map
		wg             sync.WaitGroup
		db             *gorm.DB
		instInitialing sync.Map
		redisCli       *redis.UniversalClient
	}
	type args struct {
		app string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "1", args: args{
				app: appkey,
			},
			//want: true,
			want: false,
		},

		{
			name: "1", args: args{
				app: "xxxxx",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := appDb.RemoveDB(tt.args.app); got != tt.want {
				t.Errorf("RemoveDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStorage_GetReadDB(t *testing.T) {
	SetUp()
	type fields struct {
		ins            sync.Map
		readIns        sync.Map
		channelID      string
		dsnMap         sync.Map
		wg             sync.WaitGroup
		db             *gorm.DB
		instInitialing sync.Map
		redisCli       *redis.UniversalClient
	}
	type args struct {
		app string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gorm.DB
		wantNil bool
	}{
		{
			name: "1", args: args{
				app: appkey,
			},
			wantNil: false,
		},

		{
			name: "1", args: args{
				app: "xxxxx",
			},
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := appDb.GetReadDB(tt.args.app); (got == nil && !tt.wantNil) || (got != nil && tt.wantNil) {
				t.Errorf("GetReadDB() = %v, want %v", got, tt.want)
			} else {
				t.Logf("GetReadDB() = %v, wantNil: %v", got, tt.wantNil)
			}
		})
	}
}
