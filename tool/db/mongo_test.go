package db

import (
	"context"
	"github.com/rs/xid"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type Message struct {
	Type     int8
	Sender   string
	Receiver string
	//文本内容 “这是一个文本”
	Content []byte
	Seq     int64
	ID      int64
	//引用消息ID
	QuoID int64
	//通知类型
	NtfType int8
}

type Room struct {
	Type    int8
	ID      int64
	Members []int64
	MsgList []Message
	Attr    struct {
		Name string
		Img  string
	}
}

func GetMongo(ctx context.Context) (client *mongo.Client, err error) {
	var mgConf = &config.MongoConfig{
		Addr:       "172.16.10.86:27017",
		Password:   "",
		User:       "root",
		AuthSource: "",
		Db:         "test",
		IsCluster:  false,
		IsSrv:      false,
	}
	return MongoDB(mgConf, MgoOptWithLogDriver(logger.GetLogger()), MgoOptWithContext(ctx))
}

func GetMongo2(ctx context.Context) (client *mongo.Client, err error) {
	var mgConf = &config.MongoConfig{
		Dsn:        "mongodb://172.16.10.86:27017/?authMechanism=SCRAM-SHA-1",
		AuthSource: "",
		Db:         "test",
		LogLevel:   4,
		TraceOn:    true,
	}
	return MongoDB(mgConf, MgoOptWithLogDriver(logger.GetLogger()), MgoOptWithContext(ctx))
}

func TestConnect(t *testing.T) {
}

func TestInsert(t *testing.T) {
	MongoCli, _ = GetMongo2(context.Background())
	MongoCli.Database("test").Collection("room").Drop(context.TODO())
	var seq int64 = 20
	for i := 0; i < 100000; i++ {
		_, err := MongoCli.Database("test").Collection("room").InsertOne(context.Background(), Room{
			Type:    1,
			ID:      120,
			Members: []int64{111, 222},
			MsgList: []Message{
				{
					Type:     1,
					Sender:   "111",
					Receiver: "222",
					Content:  []byte(xid.New().String()),
					Seq:      seq + 1,
					ID:       100203011102 + seq + 1,
					QuoID:    0,
					NtfType:  1,
				},
				{
					Type:     1,
					Sender:   "222",
					Receiver: "111",
					Content:  []byte(xid.New().String()),
					Seq:      seq + 2,
					ID:       100203011102 + seq + 2,
					QuoID:    0,
					NtfType:  1,
				},
			},
			Attr: struct {
				Name string
				Img  string
			}{
				Name: "test room",
				Img:  "https://www.baidu.com",
			},
		})
		if err != nil {
			logger.Error("%v", err)
			break
		}
		//logger.INFO("id:", res.InsertedID)
		seq += 2
	}
}

func TestFind(t *testing.T) {
	cursor, err := MongoCli.Database("test").Collection("room").Find(context.TODO(), bson.D{})
	if err != nil {
		logger.SetErr(err).Error("%v", err)
		return
	}
	for cursor.Next(context.TODO()) {
		logger.Error("get:", cursor.Current.String())
	}
}

func TestUpdate(t *testing.T) {
	var ctx = trace.NewTrace()
	mgoCli, err := GetMongo2(ctx)
	if err != nil {
		logger.SetErr(err).Error("%v", err)
		return
	}
	if upRe, err := mgoCli.Database("test").Collection("room").
		UpdateOne(ctx, bson.D{{"id", 120}}, bson.D{{"$set", bson.D{{"type", 2}}}}); err != nil {
		logger.SetErr(err).Error("%v", err)
	} else {
		logger.Info("update re : %v", upRe.ModifiedCount)
	}
}

func TestLoopOperate(t *testing.T) {
	var ctx = trace.NewTrace()
	mgoCli, err := GetMongo2(ctx)
	if err != nil {
		logger.SetErr(err).Error("%v", err)
		return
	}

	for i := 0; i < 10000; i++ {
		if re, err := mgoCli.Database("test").Collection("room").InsertOne(ctx,
			bson.M{"_id": primitive.NewObjectID(), "name": "haha", "age": 18, "weight": 65.5}, options.InsertOne()); err != nil {
			logger.Error(" insert err", err)
		} else {
			logger.Info("insert re %+v", re)
		}

		if upRe, err := mgoCli.Database("test").Collection("room").
			UpdateOne(ctx, bson.D{{"id", 120}}, bson.D{{"$set", bson.D{{"type", 2}}}}); err != nil {
			logger.SetErr(err).Error("%v", err)
		} else {
			logger.Info("update re : %v", upRe.ModifiedCount)
		}

		if upRe, err := mgoCli.Database("test").Collection("room").
			UpdateOne(ctx, bson.D{{"cc", 65.3}}, bson.D{{"$set", bson.D{{"type", 2}}}}, options.Update().SetUpsert(true)); err != nil {
			logger.SetErr(err).Error("%v", err)
		} else {
			logger.Info("update re : %v", upRe.ModifiedCount)
		}

		cursor, err := MongoCli.Database("test").Collection("room").Find(context.TODO(), bson.D{})
		if err != nil {
			logger.SetErr(err).Error("%v", err)
			return
		}
		for cursor.Next(context.TODO()) {
			logger.Debug("get:", cursor.Current.String())
		}

		if delRe, err := mgoCli.Database("test").Collection("room").DeleteOne(ctx, bson.D{{"name", "haha"}}); err != nil {
			logger.SetErr(err).Error("%v", err)
		} else {
			logger.Info("delete re : %v", delRe.DeletedCount)
		}
	}
}

// 并发测试 map delete
func TestDelMgoMap(t *testing.T) {
	var mp = map[int64]*event.CommandStartedEvent{}

	go func() {
		var i int8 = 0
		for {
			mp[int64(i)] = &event.CommandStartedEvent{}
			i++
		}
	}()

	go func() {
		var i int8 = 0
		for {
			delete(mp, int64(i))
			i++
		}
	}()

	time.Sleep(time.Second * 100)
}
