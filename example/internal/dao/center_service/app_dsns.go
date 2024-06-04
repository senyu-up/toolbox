package center_service

import (
	"context"
	"errors"
	"github.com/senyu-up/toolbox/example/internal/model"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/struct_tool"
	"gorm.io/gorm"
	"reflect"
)

// buildAppDsnsCond
//
//	@Description: 构建查询条件
//	@param db  body any true "-"
//	@param filter  body any true "-"
//	@return *gorm.DB
func buildAppDsnsCond(db *gorm.DB, filter model.CenterServiceAppDsnsFilter) (*gorm.DB, error) {
	if reflect.DeepEqual(filter, model.CenterServiceAppDsnsFilter{}) {
		return db, errors.New("filter is empty") // 如果传了空的过滤条件导致全表更新、删除，
	}
	if 0 < filter.Id {
		db = db.Where("id = ?", filter.Id)
	}
	if 0 < filter.AppId {
		db = db.Where("app_id = ?", filter.AppId)
	}
	if 0 < len(filter.AppKey) {
		db = db.Where("app_key = ?", filter.AppKey)
	}
	if 0 < len(filter.XxxName) {
		db = db.Where("xxx_name = ?", filter.XxxName)
	}
	if nil != filter.IsDelete {
		db = db.Where("is_delete = ?", *filter.IsDelete)
	}
	return db, nil
}

func buildAppDsnsPaginate(db *gorm.DB, page int, limit int) *gorm.DB {
	if limit < 1 {
		limit = 20
	}
	if 1 > page {
		page = 1
	}
	db = db.Limit(limit)
	var theOffset = (page - 1) * limit
	if 0 < theOffset {
		db = db.Offset(theOffset)
	}
	return db
}

func GetAppDsnsByCond(ctx context.Context, tx *gorm.DB, filter model.CenterServiceAppDsnsFilter,
	fields ...string) (res []*model.CenterServiceAppDsns, err error) {

	res = make([]*model.CenterServiceAppDsns, 0)
	dbHan := tx.Table("app_dsns").Model(&model.CenterServiceAppDsns{}).Select(fields).Debug()
	dbHan, _ = buildAppDsnsCond(dbHan, filter)
	err = dbHan.Find(&res).Error
	if err != nil {
		logger.Ctx(ctx).Error("GetAppDsnsByCond err filter(%+v) err(%+v)", filter, err)
	}
	return
}

func GetAppDsnsOneByCond(ctx context.Context, tx *gorm.DB, filter model.CenterServiceAppDsnsFilter,
	fields ...string) (res *model.CenterServiceAppDsns, err error) {

	res = new(model.CenterServiceAppDsns)
	dbHan := tx.Table(res.TableName()).Model(&model.CenterServiceAppDsns{}).Select(fields).Debug()
	dbHan, _ = buildAppDsnsCond(dbHan, filter) // 查询忽略空filter错误
	err = dbHan.Take(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 如果为空，则返回 nil
		}
		logger.Ctx(ctx).Warn("GetAppDsnsOneByCond err filter(%+v) err(%+v)", filter, err)
		return nil, err
	}
	return
}

func GetAppDsnsPaginate(ctx context.Context, tx *gorm.DB, filter model.CenterServiceAppDsnsFilter,
	page int, limit int, fields ...string) (res []*model.CenterServiceAppDsns, total int64, err error) {

	res = make([]*model.CenterServiceAppDsns, 0)
	dbHan := tx.Table("app_dsns").Model(&model.CenterServiceAppDsns{}).Select(fields).Debug()
	dbHan, _ = buildAppDsnsCond(dbHan, filter) // 查询忽略空filter错误
	dbHan.Count(&total)
	dbHan = buildAppDsnsPaginate(dbHan, page, limit)
	err = dbHan.Find(&res).Error
	if err != nil {
		logger.Ctx(ctx).Error("GetAppDsnsPaginate err filter(%+v) err(%+v)", filter, err)
	}
	return
}

func GetAppDsnsMapByCond(ctx context.Context, tx *gorm.DB, filter model.CenterServiceAppDsnsFilter,
	fields ...string) (res map[int32]*model.CenterServiceAppDsns, err error) {

	var dbRe = make([]*model.CenterServiceAppDsns, 0)
	dbHan := tx.Table("app_dsns").Model(&model.CenterServiceAppDsns{}).Select(fields).Debug()
	dbHan, _ = buildAppDsnsCond(dbHan, filter) // 查询忽略空filter错误
	err = dbHan.Find(&dbRe).Error
	if err != nil {
		logger.Ctx(ctx).Error("GetAppDsnsMapByCond err filter(%+v) err(%+v)", filter, err)
	}
	res = make(map[int32]*model.CenterServiceAppDsns, len(dbRe))
	for _, row := range dbRe {
		res[row.Id] = row
	}
	return
}

func UpdateAppDsns(ctx context.Context, tx *gorm.DB, filter model.CenterServiceAppDsnsFilter,
	data *model.CenterServiceAppDsns) (aff int64, err error) {

	var dataMap = struct_tool.StructToMapByTag(*data, "json", nil, true)
	if 0 == len(dataMap) {
		return 0, errors.New("empty update data")
	}
	dbHan := tx.Table("app_dsns").Model(&model.CenterServiceAppDsns{}).Debug()
	dbHan, err = buildAppDsnsCond(dbHan, filter)
	if err != nil { // 更新动作，需要检查条件是否为空
		return 0, errors.New("update filter is empty")
	}
	dbHan.Updates(dataMap)
	err = dbHan.Error
	if err != nil {
		logger.Ctx(ctx).Error("UpdateAppDsns err filter(%+v) params(%+v) err(%+v)", filter, data, err)
	} else {
		aff = dbHan.RowsAffected
		return
	}
	return
}

func DeleteAppDsns(ctx context.Context, tx *gorm.DB, filter model.CenterServiceAppDsnsFilter) (aff int64, err error) {
	dbHan := tx.Table("app_dsns").Model(&model.CenterServiceAppDsns{}).Debug()
	dbHan, err = buildAppDsnsCond(dbHan, filter)
	if err != nil { // 更新动作，需要检查条件是否为空
		return 0, errors.New("delete filter is empty")
	}
	err = dbHan.Delete(&model.CenterServiceAppDsns{}).Error
	if err != nil {
		logger.Ctx(ctx).Error("DeleteAppDsns err filter(%+v) err(%+v)", filter, err)
	} else {
		aff = dbHan.RowsAffected
		return
	}
	return
}

func InsertAppDsns(ctx context.Context, tx *gorm.DB, data *model.CenterServiceAppDsns) (aff int64, err error) {
	dbHan := tx.Table("app_dsns").Model(&model.CenterServiceAppDsns{}).Create(data).Debug()
	err = dbHan.Error
	if err != nil {
		logger.Ctx(ctx).Error("InsertAppDsns err params(%+v) err(%+v)", data, err)
	} else {
		aff = dbHan.RowsAffected
		return
	}
	return
}
