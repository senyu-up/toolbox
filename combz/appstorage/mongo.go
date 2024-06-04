package appstorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/db"
	"github.com/senyu-up/toolbox/tool/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDsn struct {
	MongoDsn  string `json:"mongo_dsn"  gorm:"column:mongo_dsn"` // primary
	AppSecret string `json:"app_secret" gorm:"column:app_secret"`
	AppKey    string `json:"app_key"    gorm:"column:app_key"`
	AppId     int32  `json:"app_id"     gorm:"column:app_id"`
}

// MongoStorage 管理接入项目的数据库连接实例
type MongoStorage struct {
	DBStorage // 继承，复用部分方法，不复用的已重写
}

// updateIns 修改传入游戏的数据库信息，因为原有setIns是存在的时候不会更新所以新增一个方法
func (d *MongoStorage) updateIns(cg DSNCategory, app string, dbIns *mongo.Database) {
	switch cg {
	case ReadDB:
		d.readIns.Store(app, dbIns)
	case DB:
		d.ins.Store(app, dbIns)
	}
}

func (d *MongoStorage) setIns(ctx context.Context, cg DSNCategory, app string, dbIns *mongo.Database, update ...bool) {
	if len(update) > 0 && update[0] {
		d.updateIns(cg, app, dbIns)
		return
	}
	switch cg {
	case ReadDB:
		if ins, has := d.readIns.LoadOrStore(app, dbIns); has { // 加入到 从库 inst map
			if conn, ok := ins.(*mongo.Database); ok {
				if err := conn.Client().Ping(ctx, nil); err != nil {
					logger.Error("read db ping err: ", err)
					conn.Client().Disconnect(ctx)
				}
			}
		}
	case DB:
		if ins, has := d.ins.LoadOrStore(app, dbIns); has { // 加入到 主库 inst map
			if conn, ok := ins.(*mongo.Database); ok {
				if err := conn.Client().Ping(ctx, nil); err != nil {
					logger.Error("read db ping err: ", err)
					conn.Client().Disconnect(ctx)
				}
			}
		}
	}
}

func (d *MongoStorage) asyncInit(ctx context.Context) error {
	if err := d.allDSN(); err != nil {
		return err
	}
	d.instInitialing = sync.Map{}
	d.dsnMap.Range(func(key, _ interface{}) bool {
		// 主从同步
		d.instInitialing.Store(d.getInstStatusKey(ReadDB, key.(string)), StatusInitializing)
		d.instInitialing.Store(d.getInstStatusKey(DB, key.(string)), StatusInitializing)

		d.wg.Add(1)
		go func(app string) {
			defer func() {
				d.instInitialing.Delete(d.getInstStatusKey(ReadDB, app))
				d.instInitialing.Delete(d.getInstStatusKey(DB, app))
				if err := recover(); err != nil {
					logger.Error("recover app connect err: ", err)
				}
				d.wg.Done()
			}()
			if err := d.connectApp(context.Background(), ALL, app); err != nil {
				logger.Error("app connect err: ", err)
				return
			}
		}(key.(string))
		return true
	})
	return nil
}

// 通过 appKey 检查 dsn是否加载
func (d *MongoStorage) checkConn(app string) bool {
	_, ok := d.dsnMap.Load(app)
	return ok
}

// removeApp
//
//	@Description: 删除指定 app 的 mysql 连接实例，注：不会删除 dsn 记录
//	@receiver d
//	@param app  body any true "app key"
//	@return bool
func (d *MongoStorage) removeApp(app string) bool {
	var doDelete = false
	d.dsnMap.Delete(app)
	ins, ok := d.ins.Load(app)
	if ok && !isUsedMongoConn(ins.(*mongo.Database)) {
		d.ins.Delete(app)
		doDelete = true
	}
	rIns, ok := d.readIns.Load(app)
	if ok && !isUsedMongoConn(rIns.(*mongo.Database)) {
		d.readIns.Delete(app)
		doDelete = true
	}
	logger.Info("db storage removed app: ， delete op: ", app, doDelete)
	return doDelete
}

// connectApp
//
//	@Description: 通过 dsn 类型创建 gorm db 链接
//	@receiver d
//	@param cg   body any true "好像没有 ReadDB，DB 场景，默认都是 ALL"
//	@param app  body any true "app key"
//	@param update  body any false "已存在时，是否强制更新"
//	@return error
func (d *MongoStorage) connectApp(ctx context.Context, cg DSNCategory, app string, update ...bool) error {
	info, ok := d.dsnMap.Load(app)
	if !ok {
		update = []bool{true}
		// 第一遍没查到，查库
		if err := d.addDSN(app); err != nil {
			logger.Error("[connectApp] add dsn err: %+v app:%+v", err, app)
			return err
		}
	}
	info, ok = d.dsnMap.Load(app)
	if !ok || info == nil { // 第二遍没查到，或者为空
		return errors.New("no dsn info")
	}
	var mgoInfo = info.(mongoDsn)
	switch cg {
	case ReadDB:
		ins, err := d.connDSN(ctx, mgoInfo.MongoDsn, mgoInfo.AppId, readpref.SecondaryPreferredMode, false)
		if err != nil {
			return err
		}
		d.setIns(ctx, cg, app, ins, update...)
	case DB:
		ins, err := d.connDSN(ctx, mgoInfo.MongoDsn, mgoInfo.AppId, readpref.PrimaryMode, true)
		if err != nil {
			return err
		}
		d.setIns(ctx, cg, app, ins, update...)
	case ALL:
		ins, err := d.connDSN(ctx, mgoInfo.MongoDsn, mgoInfo.AppId, readpref.PrimaryMode, true)
		if err != nil {
			return err
		}
		d.setIns(ctx, DB, app, ins, update...)
		rIns, err := d.connDSN(ctx, mgoInfo.MongoDsn, mgoInfo.AppId, readpref.SecondaryPreferredMode, false)
		if err != nil {
			return err
		}
		d.setIns(ctx, ReadDB, app, rIns, update...)
	}
	logger.Info("app connect suc: ", app)
	return nil
}

func (d *MongoStorage) connDSN(ctx context.Context, dsn string, appId int32, prefer readpref.Mode, primary bool) (*mongo.Database, error) {
	var slow time.Duration
	if d.stage != enum.EvnStageLocal {
		slow = time.Second * 60
	} else {
		slow = time.Second * 5
	}

	prefer = readpref.SecondaryPreferredMode
	if primary {
		prefer = readpref.PrimaryMode
	}
	client, err := db.MongoDB(&config.MongoConfig{Dsn: dsn, LogLevel: int64(d.logLevel), TraceOn: true},
		db.MgoOptWithContext(ctx), db.MgoOptWithLogDriver(d.log), db.MgoOptWithSlowThreshold(slow),
		db.MgoOptWithPref(prefer))
	return client.Database(GetMongoDBName(appId)), err
}

// allDSN
//
//	@Description: 查询 app_dsns 获取所有 应用的 mysql dsn
//	@receiver d
//	@return error
func (d *MongoStorage) allDSN() error {
	//list := d.getAllDsn(d.modelType, "")
	list := make([]mongoDsn, 0)
	if re := d.getDBInst().Table("app_dsns").Select("app_id,app_key,mongo_dsn,app_secret").Find(&list); re.Error != nil {
		return re.Error
	}
	for _, info := range list {
		d.dsnMap.Store(info.AppKey, info)
	}
	return nil
}

// SetDB
//
//	@Description: 将指定 app_key 加入到连接实例，如果没有，则读取数据获取 dsn，并加入到连接实例中，如果查不到报错
//	@receiver d
//	@param app   body any true "app key"
//	@return error
func (d *MongoStorage) SetDB(ctx context.Context, app string) error {
	ins, has := d.ins.Load(app)
	if has && ins != nil {
		dbIns := ins.(*mongo.Database)
		if dbIns != nil {
			_ = dbIns.Client().Disconnect(ctx)
		}
	}
	if err := d.addDSN(app); err != nil {
		logger.Error("dsn notify consumer add dsn err: %+v app:%+v", err, app)
		return err
	}
	err := d.connectApp(ctx, ALL, app, true)
	if err != nil {
		return err
	}
	ins, ok := d.ins.Load(app)
	if !ok || ins == nil {
		d.removeApp(app)
		logger.Error("db write map fail")
		return errors.New("db write map fail")
	}
	return nil
}

// RemoveDB
//
//	@Description: 删除指定 app 的 mysql 连接实例，注：不会删除 dsn 记录. 2. 删除时会检查是否还在用，一般都是在用，所以不会进行实际的 inst 删除。
//	@receiver d
//	@param app   body any true "app key"
//	@return bool
func (d *MongoStorage) RemoveDB(app string) bool {
	return d.removeApp(app)
}

// GetReadDB 获取读连接
func (d *MongoStorage) GetReadDB(ctx context.Context, app string) *mongo.Database {
	ins, has := d.readIns.Load(app)
	if has && ins != nil {
		return ins.(*mongo.Database)
	}
	// 判断是否在初始化中
	if _, exists := d.instInitialing.Load(d.getInstStatusKey(ReadDB, app)); exists {
		time.Sleep(time.Millisecond * 100)
		return d.GetReadDB(ctx, app)
	}
	// 没有初始化，则初始化
	err := d.connectApp(ctx, ALL, app)
	if err != nil {
		logger.Error("[GetReadDB] connectApp err: %+v app:%+v", err, app)
	}
	ins, _ = d.readIns.Load(app)
	if ins != nil {
		return ins.(*mongo.Database)
	} else {
		d.removeApp(app)
	}
	return nil
}

// GetWriteDB
//
//	@Description: 通过 app key 获取 mongo 主库连接，如果没有则返回 nil，使用前记得判断空
//	@receiver d
//	@param app   body any true "app key"
//	@return *mongo.Database
func (d *MongoStorage) GetWriteDB(ctx context.Context, app string) *mongo.Database {
	// 判断当前app是否在初始化中
	ins, has := d.ins.Load(app)

	if has && ins != nil {
		return ins.(*mongo.Database)
	}
	// 判断是否在初始化中
	if _, exists := d.instInitialing.Load(d.getInstStatusKey(DB, app)); exists {
		time.Sleep(time.Millisecond * 100)
		return d.GetWriteDB(ctx, app)
	}
	logger.Warn("try to add a not initial db, app: %s", app)
	if err := d.addDSNForce(app); err != nil {
		logger.Error("dsn notify consumer add dsn err: %+v app:%+v", err, app)
		return nil
	}
	err := d.connectApp(ctx, ALL, app)
	if err != nil {
		logger.Error("dsn notify consumer add dsn err: %+v app:%+v", err, app)
	}
	ins, _ = d.ins.Load(app)
	if ins != nil {
		return ins.(*mongo.Database)
	} else {
		d.removeApp(app)
	}
	return nil
}

// GetDB 通过 app_key 和 readPref 获取 mongo 读或写连接，如果没有则返回 nil，使用前记得判断空
// prefer 字段，如果传0，则返回默认读连接
// prefer 不支持  readpref.NearestMode ⚠️
func (d *MongoStorage) GetDB(ctx context.Context, app string, prefer readpref.Mode) *mongo.Database {
	if prefer == readpref.NearestMode {
		logger.Ctx(ctx).Error("[GetDB]not supported prefer mode [%d]", prefer)
		return nil // 不支持这个
	}
	if prefer == 0 {
		prefer = readpref.SecondaryPreferredMode
	}
	if prefer == readpref.SecondaryPreferredMode || prefer == readpref.SecondaryMode {
		return d.GetReadDB(ctx, app)
	} else if prefer == readpref.PrimaryPreferredMode || prefer == readpref.PrimaryMode {
		return d.GetWriteDB(ctx, app)
	}
	return nil
}

// GetAllDBMap
//
//	@Description: 获取所有的 gorm 连接实例, 按照 app_key 作为 key 返回一个map。
//	@receiver d
//	@return map[string]*mongo.Database
func (d *MongoStorage) GetAllDBMap() map[string]*mongo.Database {
	var re = make(map[string]*mongo.Database, 0)
	d.ins.Range(func(key, value interface{}) bool {
		re[key.(string)] = value.(*mongo.Database)
		return true
	})

	return re
}

// NewMongoStorageDB
//
//	@Description: 初始化
//	@param opts	  body any false "options func"
//	@return *MongoStorage
//	@return error
func NewMongoStorageDB(opts ...StoreOption) (*MongoStorage, error) {
	var originDb = DBStorage{ins: sync.Map{}, readIns: sync.Map{}, dsnMap: sync.Map{}, wg: sync.WaitGroup{}, instInitialing: sync.Map{}}
	for _, opt := range opts {
		opt(&originDb)
	}
	var storage = &MongoStorage{DBStorage: originDb}
	var ctx = context.Background()
	if storage.ctx != nil {
		ctx = storage.ctx
	}

	//订阅 apps dsn change
	storage.runAppsStorageNotifierListener(ctx, *storage.redisCli)
	if storage.initImmediately {
		return storage, storage.asyncInit(ctx)
	}
	return storage, nil // 默认懒加载，不立即初始化
}

func isUsedMongoConn(ins *mongo.Database) bool {
	//var ctx = context.Background()	// TODO
	//ins.Database("test").RunCommand(ctx, bson.D{{"ping", 1}})
	//i, _ := ins.DB()
	//return i.Stats().InUse == 0
	return false
}

func GetMongoDBName(appId int32) string {
	return fmt.Sprintf("GP_%d", appId)
}
