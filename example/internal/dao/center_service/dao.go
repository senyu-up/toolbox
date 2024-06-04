package center_service

import "gorm.io/gorm"

// 数据库 center_service 实例

var (
	centerDb *gorm.DB
)

func Init(db *gorm.DB) {
	centerDb = db
}

func GetDB() *gorm.DB {
	return centerDb
}
