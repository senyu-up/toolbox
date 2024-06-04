package appstorage

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/db"
	"github.com/senyu-up/toolbox/tool/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

var (
	appMgoDb *MongoStorage
)

func SetUpForMongo() {
	if appMgoDb != nil {
		return
	}
	var redConf = config.RedisConfig{
		//Addrs: []string{"127.0.0.1:6379"},
		Addrs: []string{"redis-node-0:6379"},
	}
	var err error
	redisCli = cache.InitRedisByConf(&redConf)

	var dbConf = config.MysqlConfig{
		Master: config.MysqlSingleConfig{
			Addr:     "172.16.10.86:3306",
			User:     "root",
			Password: "",
			Db:       "center_service",
		}}
	dbCli, err := db.NewMysql(&dbConf)
	if err != nil {
		fmt.Printf("conn mysql err %v \n", err)
		return
	}

	appMgoDb, err = NewMongoStorageDB(StoreOptionWithGorm(dbCli),
		StoreOptionWithRedisCli(&redisCli),
		StoreOptionWithChannelID("xh_test"),
		StoreOptionWithLogLevel(4),
		StoreOptionWithAppStage(enum.EvnStageLocal))
	if err != nil {
		fmt.Printf("new app storage  err %v \n", err)
		return
	}

	RegisterAddEvent(func(appKey string) error {
		fmt.Printf("test func call , add event %s", appkey)
		return nil
	})

	RegisterUpdateEvent(func(appKey string) error {
		fmt.Printf("test func call , update event %s", appkey)
		return nil
	})

	RegisterDelEvent(func(appKey string) error {
		fmt.Printf("test func call , del event %s", appkey)
		return nil
	})

	time.Sleep(time.Second)
}

func TestNewMongoStorageDB(t *testing.T) {
	SetUpForMongo()
	type args struct {
		opts []StoreOption
	}
	tests := []struct {
		name    string
		args    args
		want    *MongoStorage
		wantErr bool
	}{
		{
			name: "1", args: args{opts: []StoreOption{StoreOptionWithRedisCli(&redisCli)}}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMongoStorageDB(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMongoStorageDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NewMongoStorageDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoStorageQuery(t *testing.T) {
	SetUpForMongo()

	var ctx = context.Background()
	var dbClient = appMgoDb.GetReadDB(ctx, appkey)
	if err := dbClient.Client().Ping(ctx, nil); err != nil {
		t.Errorf("ping db err %v", err)
	} else {
		t.Logf("ping %s db ok", appkey)
	}
}

func TestMongoStorageInsert(t *testing.T) {
	SetUpForMongo()

	var ctx = context.Background()
	var dbClient = appMgoDb.GetDB(ctx, appkey, readpref.PrimaryPreferredMode)
	if re, err := dbClient.Collection("test_log").InsertOne(ctx,
		bson.M{"_id": primitive.NewObjectID(), "name": "娃haha", "age": 18, "weight": 65.5}, options.InsertOne()); err != nil {
		logger.Error(" insert err", err)
	} else {
		logger.Info("insert re %+v", re)
	}
}
