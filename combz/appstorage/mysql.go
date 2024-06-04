package appstorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/db"
	"github.com/senyu-up/toolbox/tool/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type DSNCategory int

// 数据库链接类型控制
const (
	ReadDB DSNCategory = 1
	DB     DSNCategory = 2
	ALL    DSNCategory = 3
)

type dsn struct {
	DSN       string `json:"dsn"`
	DSNSlave  string `json:"dsn_slave"`
	AppSecret string `json:"app_secret"`
	AppKey    string `json:"app_key"`
}

const (
	// 初始化中, 维护中间状态
	StatusInitializing = 1
)

// DBStorage 管理接入项目的数据库连接实例
type DBStorage struct {
	// stage
	stage string
	//nsq/redis topic channel
	channelID string
	//存储实例
	ins sync.Map
	//读库存储实例
	readIns sync.Map
	//all dsn
	dsnMap sync.Map
	//wg
	wg sync.WaitGroup
	db *gorm.DB
	// 应用名
	appName string

	// 日志驱动
	log logger.Log
	// 日志登记
	logLevel glogger.LogLevel

	// 主库状态控制
	instInitialing sync.Map

	redisCli *redis.UniversalClient

	mysqlConf *config.MysqlConfig // 默认的 mysql config

	initImmediately bool // 是否立即初始化所有连接
	ctx             context.Context
}

func (d *DBStorage) getInstStatusKey(cg DSNCategory, app string) string {
	return fmt.Sprintf("%d_%s", cg, app)
}

// updateIns 修改传入游戏的数据库信息，因为原有setIns是存在的时候不会更新所以新增一个方法
func (d *DBStorage) updateIns(cg DSNCategory, app string, dbIns *gorm.DB) {
	switch cg {
	case ReadDB:
		d.readIns.Store(app, dbIns)
	case DB:
		d.ins.Store(app, dbIns)
	}
	conn, _ := dbIns.DB()
	conn.SetMaxIdleConns(100)
	conn.SetConnMaxLifetime(time.Hour)
}

func (d *DBStorage) setIns(cg DSNCategory, app string, dbIns *gorm.DB, update ...bool) {
	if len(update) > 0 && update[0] {
		d.updateIns(cg, app, dbIns)
		return
	}
	switch cg {
	case ReadDB:
		if ins, has := d.readIns.LoadOrStore(app, dbIns); has {
			conn, _ := ins.(*gorm.DB).DB()
			// SetMaxIdleConns 设置空闲连接池中连接的最大数量
			conn.SetMaxIdleConns(d.mysqlConf.MaxIdleConn)
			conn.SetMaxOpenConns(d.mysqlConf.MaxOpenConn) // 最大连接数量
			conn.SetConnMaxIdleTime(d.mysqlConf.MaxIdleTime)
			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			conn.SetConnMaxLifetime(d.mysqlConf.MaxLifeTime)
			//没有在使用的情况下,关闭链接
			//if conn.Stats().InUse == 0 {
			//	conn.Close()
			//}
		}
	case DB:
		if ins, has := d.ins.LoadOrStore(app, dbIns); has {
			conn, _ := ins.(*gorm.DB).DB()
			// SetMaxIdleConns 设置空闲连接池中连接的最大数量
			conn.SetMaxIdleConns(d.mysqlConf.MaxIdleConn)
			conn.SetMaxOpenConns(d.mysqlConf.MaxOpenConn) // 最大连接数量
			conn.SetConnMaxIdleTime(d.mysqlConf.MaxIdleTime)
			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			conn.SetConnMaxLifetime(d.mysqlConf.MaxLifeTime)
			//没有在使用的情况下,关闭链接
			//if conn.Stats().InUse == 0 {
			//	conn.Close()
			//}
		}
	}
}

func (d *DBStorage) asyncInit() error {
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
			if err := d.connectApp(ALL, app); err != nil {
				logger.Error("app connect err: ", err)
				return
			}
		}(key.(string))
		return true
	})
	return nil
}

func (d *DBStorage) getDBInst() *gorm.DB {
	if d.db != nil {
		return d.db
	} else {
		panic("app storage db not init")
	}
}

// 检查该app是否连接
func (d *DBStorage) checkConn(app string) bool {
	_, ok := d.dsnMap.Load(app)
	return ok
}

func (d *DBStorage) addDSN(app string) error {
	info := dsn{}
	tx := d.getDBInst().Table("app_dsns")

	if err := tx.Where("app_dsns.app_key = ?", app).First(&info).Error; err != nil {
		return err
	}
	d.dsnMap.Store(app, info)
	return nil
}

// addDSNForce
//
//	@Description: 如果非要使用这个 app 的dsn，那就调用这个方法吧，去掉 model type 限制
//	@receiver d
//	@param app  body any true "app key"
//	@return error
func (d *DBStorage) addDSNForce(app string) error {
	info := dsn{}

	if err := d.getDBInst().Table("app_dsns").Where("app_key = ?", app).First(&info).Error; err != nil {
		return err
	}
	d.dsnMap.Store(app, info)
	return nil
}

// removeApp
//
//	@Description: 删除指定 app 的 mysql 连接实例，注：不会删除 dsn 记录
//	@receiver d
//	@param app  body any true "app key"
//	@return bool
func (d *DBStorage) removeApp(app string) bool {
	var doDelete = false
	d.dsnMap.Delete(app)
	ins, ok := d.ins.Load(app)
	if ok && !isUsedConn(ins.(*gorm.DB)) {
		d.ins.Delete(app)
		doDelete = true
	}
	rIns, ok := d.readIns.Load(app)
	if ok && !isUsedConn(rIns.(*gorm.DB)) {
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
func (d *DBStorage) connectApp(cg DSNCategory, app string, update ...bool) error {
	info, ok := d.dsnMap.Load(app)
	if !ok {
		return errors.New("no dsn info")
	}
	switch cg {
	case ReadDB:
		ins, err := d.connDSN(info.(dsn).DSNSlave)
		if err != nil {
			return err
		}
		d.setIns(cg, app, ins, update...)
	case DB:
		ins, err := d.connDSN(info.(dsn).DSN)
		if err != nil {
			return err
		}
		d.setIns(cg, app, ins, update...)
	case ALL:
		ins, err := d.connDSN(info.(dsn).DSN)
		if err != nil {
			return err
		}
		d.setIns(DB, app, ins, update...)
		rIns, err := d.connDSN(info.(dsn).DSNSlave)
		if err != nil {
			return err
		}
		d.setIns(ReadDB, app, rIns, update...)
	}
	logger.Info("app connect suc: ", app)
	return nil
}

func (d *DBStorage) connDSN(dsn string) (*gorm.DB, error) {
	var slow time.Duration
	if d.logLevel == 0 {
		if d.stage != enum.EvnStageLocal {
			d.logLevel = glogger.Error
		} else {
			d.logLevel = glogger.Info
			slow = time.Second * 5
		}
	}

	return gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{Logger: db.NewGormLogger(
			db.GormLogOptWithLevel(d.logLevel),
			db.GormLogOptWithTraceOn(true),
			db.GormLogOptWithLogDriver(d.log),
			db.GormLogOptWithSlowThreshold(slow))},
	)
}

// allDSN
//
//	@Description: 查询 app_dsns 获取所有 应用的 mysql dsn
//	@receiver d
//	@return error
func (d *DBStorage) allDSN() error {
	list := make([]dsn, 0)
	tx := d.getDBInst().Table("app_dsns")
	//if d.modelType != 0 {
	//	tx = tx.Joins("left join center_models on center_models.app_key = app_dsns.app_key").
	//		Where("model_type = ?", d.modelType).Group("app_dsns.app_key")
	//}

	if err := tx.Find(&list).Error; err != nil {
		return err
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
func (d *DBStorage) SetDB(app string) error {
	ins, has := d.ins.Load(app)
	if has && ins != nil {
		dbIns := ins.(*gorm.DB)
		sqlIns, _ := dbIns.DB()
		if sqlIns != nil {
			_ = sqlIns.Close()
		}
	}
	if err := d.addDSN(app); err != nil {
		logger.Error("dsn notify consumer add dsn err: %+v app:%+v", err, app)
		return err
	}
	err := d.connectApp(ALL, app, true)
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
func (d *DBStorage) RemoveDB(app string) bool {
	return d.removeApp(app)
}

// GetDB
//
//	@Description: 通过 app key 获取 gorm连接，如果没有则返回 nil，使用 gorm 链接前，记得判断空
//	@receiver d
//	@param app   body any true "app key"
//	@return *gorm.DB
func (d *DBStorage) GetDB(app string) *gorm.DB {
	// 判断当前app是否在初始化中
	ins, has := d.ins.Load(app)

	if has && ins != nil {
		return ins.(*gorm.DB)
	}
	// 判断是否在初始化中
	if _, exists := d.instInitialing.Load(d.getInstStatusKey(DB, app)); exists {
		time.Sleep(time.Millisecond * 100)
		return d.GetDB(app)
	}
	logger.Warn("try to add a not initial db, app: %s", app)
	if err := d.addDSNForce(app); err != nil {
		logger.Error("dsn notify consumer add dsn err: %+v app:%+v", err, app)
		return nil
	}
	err := d.connectApp(ALL, app)
	if err != nil {
		logger.Error("dsn notify consumer add dsn err: %+v app:%+v", err, app)
	}
	ins, _ = d.ins.Load(app)
	if ins != nil {
		return ins.(*gorm.DB)
	} else {
		d.removeApp(app)
	}
	return nil
}

// GetAllDBMap
//
//	@Description: 获取所有的 gorm 连接实例, 按照 app_key 作为 key 返回一个map。
//	@receiver d
//	@return map[string]*gorm.DB
func (d *DBStorage) GetAllDBMap() map[string]*gorm.DB {
	var re = make(map[string]*gorm.DB, 0)
	d.ins.Range(func(key, value interface{}) bool {
		re[key.(string)] = value.(*gorm.DB)
		return true
	})

	return re
}

func (d *DBStorage) GetReadDB(app string) *gorm.DB {
	ins, has := d.readIns.Load(app)
	if has && ins != nil {
		return ins.(*gorm.DB)
	}
	// 判断是否在初始化中
	if _, exists := d.instInitialing.Load(d.getInstStatusKey(ReadDB, app)); exists {
		time.Sleep(time.Millisecond * 100)
		return d.GetReadDB(app)
	}
	_ = d.connectApp(ALL, app)
	ins, _ = d.ins.Load(app)
	if ins != nil {
		return ins.(*gorm.DB)
	} else {
		d.removeApp(app)
	}
	return nil
}

// NewAppStorageDB
//
//	@Description: 初始化
//	@param opts	  body any false "options func"
//	@return *DBStorage
//	@return error
func NewAppStorageDB(opts ...StoreOption) (*DBStorage, error) {
	var storage = &DBStorage{
		ins: sync.Map{},
		wg:  sync.WaitGroup{},
		mysqlConf: &config.MysqlConfig{
			MaxOpenConn: 100,
			MaxIdleConn: 100,
			MaxIdleTime: time.Duration(60),
			MaxLifeTime: time.Duration(3600),
		},
	}
	for _, opt := range opts {
		opt(storage)
	}
	//if storage.modelType == 0 && 0 < len(storage.appName) {
	//	if mt, ok := GetAppNameModuleTypeMap()[storage.appName]; ok {
	//		storage.modelType = mt
	//	}
	//}
	//订阅 apps dsn change
	storage.runAppsStorageNotifierListener(*storage.redisCli)
	return storage, storage.asyncInit()
}

func isUsedConn(ins *gorm.DB) bool {
	i, _ := ins.DB()
	return i.Stats().InUse > 0
}
