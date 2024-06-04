package loader

import (
	"github.com/senyu-up/toolbox/tool/config"
	"reflect"
	"testing"
)

func TestInitConf(t *testing.T) {
	type args struct {
		loader Loader
		param  []ConfOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yaml"), ConfOptWithType("yaml")}},
			wantErr: false,
		},
		{
			name: "2", args: args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yaml")}},
			wantErr: false,
		},
		{
			name:    "3",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yml")}},
			wantErr: false,
		},
		{
			name:    "4",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yml"), ConfOptWithType("yaml")}},
			wantErr: false,
		},
		{
			name:    "e1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config2.yaml"), ConfOptWithType("yaml")}},
			wantErr: true,
		},
		{
			name:    "e2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config2.yaml")}},
			wantErr: true,
		},
		{
			name:    "t1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.toml"), ConfOptWithType("toml")}},
			wantErr: false,
		},
		{
			name:    "t2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.toml")}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InitConf(tt.args.loader, tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConfig_Get(t *testing.T) {
	type fields struct {
		loader Loader
	}
	type args struct {
		loader Loader
		param  []ConfOption
		key    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yaml"), ConfOptWithType("yaml")}, key: "redisdb.addr"},
			want:    "localhost:6379",
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yml")}, key: "env"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "3",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yaml")}, key: "run_env"},
			want:    "local",
			wantErr: false,
		},
		{
			name:    "4",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yml"), ConfOptWithType("yaml")}, key: "mongodb"},
			want:    map[string]string{"addr": "localhost:27017"},
			wantErr: false,
		},
		{
			name:    "5",
			args:    args{loader: &File{}, param: []ConfOption{}, key: "mysql.user"},
			want:    "root",
			wantErr: false,
		},
		{
			name:    "e1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config2.yaml"), ConfOptWithType("yaml")}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "e2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config2.yaml")}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "t1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.toml"), ConfOptWithType("toml")}, key: "redisdb.iscluster"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "t2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.toml")}, key: "redisdb.db"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "t3",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.json")}, key: "redisdb.iscluster"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := InitConf(tt.args.loader, tt.args.param...)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Unmarshal(t *testing.T) {
	type AppConf struct {
		RedisDb config.RedisConfig
		Mysql   config.MysqlConfig
		MongoDb config.MongoConfig
	}
	type fields struct {
		loader Loader
	}
	type args struct {
		loader Loader
		param  []ConfOption
		dst    interface{}
	}
	var (
		appConf1    = &AppConf{}
		errConf1    = &fields{}
		wantAppConf = &AppConf{
			RedisDb: config.RedisConfig{
				Addrs:     []string{"localhost:6379"},
				IsCluster: true,
			},
			MongoDb: config.MongoConfig{
				Addr: "localhost:27017",
			},
			Mysql: config.MysqlConfig{
				Dsn: "root:xxx@tcp(172.16.10.40:30006)/center_service??timeout=90s&collation=utf8mb4_unicode_ci",
			},
		}
	)
	wantAppConf2 := *wantAppConf
	wantAppConf2.RedisDb.DB = 1
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yaml"), ConfOptWithType("yaml")}, dst: appConf1},
			want:    wantAppConf,
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yml")}, dst: appConf1},
			want:    wantAppConf,
			wantErr: false,
		},
		{
			name:    "3",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yaml")}, dst: appConf1},
			want:    wantAppConf,
			wantErr: false,
		},
		{
			name:    "4",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yml"), ConfOptWithType("yaml")}, dst: appConf1},
			want:    wantAppConf,
			wantErr: false,
		},
		{
			name:    "e1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config.yaml"), ConfOptWithType("yaml")}, dst: errConf1},
			want:    &fields{},
			wantErr: false,
		},
		{
			name:    "e2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./config2.yaml")}, dst: errConf1},
			want:    &fields{},
			wantErr: true,
		},
		{
			name:    "t1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.toml"), ConfOptWithType("toml")}, dst: appConf1},
			want:    &wantAppConf2,
			wantErr: false,
		},
		{
			name:    "t2",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.toml")}, dst: appConf1},
			want:    &wantAppConf2,
			wantErr: false,
		},
		{
			name:    "j1",
			args:    args{loader: &File{}, param: []ConfOption{ConfOptWithPath("./conf.json")}, dst: appConf1},
			want:    &wantAppConf2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c, err := InitConf(tt.args.loader, tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := c.Unmarshal(tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.want) {
				t.Errorf("Unmarshal() got = %+v, want %+v", tt.args.dst, tt.want)
			}
		})
	}
}
