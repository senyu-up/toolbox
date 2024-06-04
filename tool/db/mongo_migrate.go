package db

import (
	"context"
	"strings"

	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/su_slice"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	primaryKeyName = "_id_"       //主键key
	hashedKeyName  = "_id_hashed" //切片所有
)

// CollIdx 集合结构体 需要实现接口
type CollIdx interface {
	// CollectionName 集合名
	CollectionName() string
	// Index 集合索引
	Index() []mongo.IndexModel
}

// CollectionIndex 注册索引的结构体，记录需要注册的索引的表集合
type CollectionIndex struct {
	// database
	db *mongo.Database
	// 集合
	collections map[string]CollIdx
}

var index *CollectionIndex

// NewCollectionIndex 创建 集合、索引
func NewCollectionIndex(d *mongo.Database) *CollectionIndex {
	if index == nil {
		index = &CollectionIndex{
			db:          d,
			collections: make(map[string]CollIdx),
		}
	}

	return index
}

// GetAllCollections 获取索引
func (i *CollectionIndex) GetAllCollections(ctx context.Context) ([]string, error) {
	return i.db.ListCollectionNames(ctx, bson.M{})
}

// CreateCollection 创建集合
func (i *CollectionIndex) CreateCollection(ctx context.Context, name string) error {
	return i.db.CreateCollection(ctx, name)
}

// GetCollectionIndex 获取集合的索引
func (i *CollectionIndex) GetCollectionIndex(ctx context.Context, name string) ([]mongo.IndexModel, error) {
	c := i.db.Collection(name)
	// 遍历所有 index
	cur, err := c.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	var result = make([]mongo.IndexModel, 0)
	//var result = make([]bson.M, 0)
	for cur.Next(ctx) {
		var idx bson.M
		if err := cur.Decode(&idx); err != nil {
			return result, err
		}
		result = append(result, bsonMToIndexModel(idx))
	}
	return result, err
}

// DeleteIndex 删除索引
func (i *CollectionIndex) DeleteIndex(ctx context.Context, name, idxName string) (err error) {
	c := i.db.Collection(name)
	if re, err := c.Indexes().DropOne(ctx, idxName); err != nil {
		logger.Ctx(ctx).Error("drop one index re %+v, err %v", re, err)
		return err
	} else {
		return nil
	}
}

// CreateIndex 创建索引
func (i *CollectionIndex) CreateIndex(ctx context.Context, name string, idxes []mongo.IndexModel) (err error) {
	c := i.db.Collection(name)
	if re, err := c.Indexes().CreateMany(ctx, idxes); err != nil {
		logger.Ctx(ctx).Error("create many index re %+v, err %v", re, err)
		return err
	} else {
		return nil
	}
}

// mergeIndex 合并索引，已经存在的索引不变， 新增的创建，多的删除
func (i *CollectionIndex) mergeIndex(ctx context.Context, table CollIdx) (err error) {
	existIdx, err := i.GetCollectionIndex(ctx, table.CollectionName())
	if err != nil {
		return err
	}

	// 获取预期的索引
	indexes := table.Index()

	var toAddIdx = make([]mongo.IndexModel, 0)
	var toDelIdx = make([]mongo.IndexModel, 0)
	// 检查需要新增的 idx
	for _, idx := range indexes {
		var b = idx.Options.Name
		if su_slice.InArray(*b, []string{primaryKeyName, hashedKeyName}) {
			continue // 主键不处理
		}
		for _, exist := range existIdx {
			var a = exist.Options.Name
			// 按照名字检查索引
			if 0 == strings.Compare(*a, *b) {
				// 检查 unique 参数
				if (idx.Options.Unique == nil && exist.Options.Unique == nil) ||
					(idx.Options.Unique != nil && exist.Options.Unique != nil && *idx.Options.Unique == *exist.Options.Unique) {
					goto NEXT // 跳两层
				}
			}
		}
		// 定义了但不存在的索引：加
		toAddIdx = append(toAddIdx, idx)
	NEXT:
		continue
	}

	// 检查需要删除的 idx
	for _, exist := range existIdx {
		var a = exist.Options.Name
		if su_slice.InArray(*a, []string{primaryKeyName, hashedKeyName}) {
			continue // 主键不处理
		}
		for _, idx := range indexes {
			var b = idx.Options.Name
			// 按照名字检查索引
			if 0 == strings.Compare(*a, *b) {
				// 检查 unique 参数
				if (idx.Options.Unique == nil && exist.Options.Unique == nil) ||
					(idx.Options.Unique != nil && exist.Options.Unique != nil && *idx.Options.Unique == *exist.Options.Unique) {
					goto NextDel // 跳两层
				}
			}
		}
		// 存在但没定义的索引：删
		toDelIdx = append(toDelIdx, exist)
	NextDel:
		continue
	}

	// 删除
	for _, idx := range toDelIdx {
		if err = i.DeleteIndex(ctx, table.CollectionName(), *idx.Options.Name); err != nil {
			return err
		}
	}
	// 新增
	if len(toAddIdx) > 0 {
		if err = i.CreateIndex(ctx, table.CollectionName(), toAddIdx); err != nil {
			return err
		}
	}
	return err
}

// bsonMToIndexModel 查询出来的 bson.M 转 mongo.IndexModel
func bsonMToIndexModel(b bson.M) mongo.IndexModel {
	var re = mongo.IndexModel{Options: &options.IndexOptions{}}
	if val, ok := b["v"]; ok {
		var i32 = val.(int32)
		re.Options.Version = &i32
	}
	if val, ok := b["name"]; ok {
		var nStr = val.(string)
		re.Options.Name = &nStr
	}
	//if val, ok := b["ns"]; ok {}
	if val, ok := b["key"]; ok {
		var bm = val.(bson.M)
		re.Keys = bm
	}
	return re
}

// RegisterTable 注册需要建索引的表
func (i *CollectionIndex) RegisterTable(table CollIdx) {
	i.collections[table.CollectionName()] = table
}

// Migrate 自动创建集合与合并索引， 返回所有报错信息
func (i *CollectionIndex) Migrate(ctx context.Context) (msg string) {
	logger.Ctx(ctx).Info("Start Mongo Migrate")
	//获取所有集合
	collections, err := i.GetAllCollections(ctx)
	if err != nil {
		return err.Error()
	}

	collectionMap := make(map[string]uint8)
	for _, collection := range collections {
		collectionMap[collection] = 1
	}

	for _, table := range i.collections {
		tableName := table.CollectionName()
		if _, ok := collectionMap[tableName]; !ok {
			//不存在就创建表
			if err = i.CreateCollection(ctx, tableName); err != nil {
				msg += "\t" + err.Error()
				continue
			}
		}

		// 索引检查与修正
		if err = i.mergeIndex(ctx, table); err != nil {
			msg += "\t" + err.Error()
			continue
		}
	}
	return msg
}
