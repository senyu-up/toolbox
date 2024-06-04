package appstorage

import (
	"context"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/runtime"
)

// nsq 通知类型
const (
	AddDataIns    DSNNotifyCategory = 1
	RemoveDataIns DSNNotifyCategory = 2 // 推送了删除 inst 消息后，原链接因为在连接池，还会保留一段时间，不要过分依赖 getDb 来判断是否删除
	UpdateDataIns DSNNotifyCategory = 3
)

type DSNNotifyMessage struct {
	AppKey   string
	Category DSNNotifyCategory
}

type DSNNotifyCategory int

// StorageSyncTopic nsq topic
const StorageSyncTopic = "SYNC_APP_STORAGE"
const MongoStorageSyncTopic = "SYNC_MONGO_APP_STORAGE"

func PushAppDSNChangeNotification(redisCli redis.UniversalClient, msg *DSNNotifyMessage) error {
	message, _ := jsoniter.MarshalToString(&msg)
	return redisCli.Publish(StorageSyncTopic, message).Err()
}

func PushMongoAppDSNChangeNotification(redisCli redis.UniversalClient, msg *DSNNotifyMessage) error {
	message, _ := jsoniter.MarshalToString(&msg)
	return redisCli.Publish(MongoStorageSyncTopic, message).Err()
}

func (s *DBStorage) runAppsStorageNotifierListener(redisCli redis.UniversalClient) {
	runtime.GOSafe(context.Background(), "app connect notify consumer", func() {
		list := redisCli.Subscribe(StorageSyncTopic).Channel()
		for true {
			info := &DSNNotifyMessage{}
			select {
			case msg := <-list:
				err := jsoniter.UnmarshalFromString(msg.Payload, &info)
				if err != nil {
					logger.Error("RunAppsStorageNotifierListener consumer unmarshal err: ", err)
				}
				switch info.Category {
				case AddDataIns:
					if err := s.addDSNForce(info.AppKey); err != nil {
						logger.Error("dsn notify consumer add dsn err: ", err)
						continue
					}
					if err := s.connectApp(ALL, info.AppKey); err != nil {
						logger.Error("dsn notify consumer connect dsn: ", err)
						continue
					}
					execAddEventFuncList(info.AppKey)
				case UpdateDataIns:
					if !s.checkConn(info.AppKey) {
						// 如果没用这个 app 数据库连接，则不更新
						continue
					}
					if err := s.addDSNForce(info.AppKey); err != nil {
						logger.Error("dsn notify consumer add dsn err: ", err)
						continue
					}
					if err := s.connectApp(ALL, info.AppKey, true); err != nil {
						logger.Error("dsn notify consumer connect dsn: ", err)
						continue
					}
					execUpdateEventFuncList(info.AppKey)
				case RemoveDataIns:
					s.removeApp(info.AppKey)
					execDelEventFuncList(info.AppKey)
				}
				logger.Warn("dsn notify consumer get msg:", info)
			}
		}
	})
}

func (d *MongoStorage) runAppsStorageNotifierListener(ctx context.Context, redisCli redis.UniversalClient) {
	runtime.GOSafe(ctx, "mongo app connect notify consumer", func() {
		var iCtx = context.Background()
		list := redisCli.Subscribe(MongoStorageSyncTopic).Channel()
		for true {
			info := &DSNNotifyMessage{}
			select {
			case msg := <-list:
				err := jsoniter.UnmarshalFromString(msg.Payload, &info)
				if err != nil {
					logger.Error("RunAppsStorageNotifierListener mongo consumer unmarshal err: ", err)
				}
				switch info.Category {
				case AddDataIns:
					if err := d.addDSNForce(info.AppKey); err != nil {
						logger.Error("mongo dsn notify consumer add dsn err: ", err)
						continue
					}
					if err := d.connectApp(iCtx, ALL, info.AppKey); err != nil {
						logger.Error("mongo dsn notify consumer connect dsn: ", err)
						continue
					}
					execAddEventFuncList(info.AppKey)
				case UpdateDataIns:
					if !d.checkConn(info.AppKey) {
						// 如果没用这个 app 数据库连接，则不更新
						continue
					}
					if err := d.addDSNForce(info.AppKey); err != nil {
						logger.Error("mongo dsn notify consumer add dsn err: ", err)
						continue
					}
					if err := d.connectApp(iCtx, ALL, info.AppKey, true); err != nil {
						logger.Error("mongo dsn notify consumer connect dsn: ", err)
						continue
					}
					execUpdateEventFuncList(info.AppKey)
				case RemoveDataIns:
					d.removeApp(info.AppKey)
					execDelEventFuncList(info.AppKey)
				}
				logger.Warn("mongo dsn notify consumer get msg:", info)
			}
		}
	})
}

type EventFunc func(appKey string) error

var (
	addEventList    []EventFunc
	updateEventList []EventFunc
	delEventList    []EventFunc
)

func RegisterAddEvent(functions ...EventFunc) {
	addEventList = append(addEventList, functions...)
}

func RegisterUpdateEvent(functions ...EventFunc) {
	updateEventList = append(updateEventList, functions...)
}

func RegisterDelEvent(functions ...EventFunc) {
	delEventList = append(delEventList, functions...)
}

func execAddEventFuncList(appKey string) {
	for i := 0; i < len(addEventList); i++ {
		if err := addEventList[i](appKey); err != nil {
			logger.Error("%v", err)
		}
	}
}

func execUpdateEventFuncList(appKey string) {
	for i := 0; i < len(updateEventList); i++ {
		if err := updateEventList[i](appKey); err != nil {
			logger.Error("%v", err)
		}
	}
}

func execDelEventFuncList(appKey string) {
	for i := 0; i < len(delEventList); i++ {
		if err := delEventList[i](appKey); err != nil {
			logger.Error("%v", err)
		}
	}
}
