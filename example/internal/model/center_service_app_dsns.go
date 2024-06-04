package model

// center_service.app_dsns
type CenterServiceAppDsns struct {
	Id        int32  `json:"id" gorm:"id"`                  // 应用ID
	AppKey    string `json:"app_key" gorm:"app_key"`        // 应用唯一标识
	Dsn       string `json:"dsn" gorm:"dsn"`                // 主数据库DSN
	DsnSlave  string `json:"dsn_slave" gorm:"dsn_slave"`    // 从数据库DSN
	CreatedAt int64  `json:"created_at" gorm:"created_at"`  // 创建时间
	UpdatedAt int64  `json:"updated_at" gorm:"updated_at"`  // 更新时间
	XxxName   string `json:"xxx_name" gorm:"xxx_name"`      // 游戏名称
	IsDelete  int8   `json:"is_delete" gorm:"is_delete"`    // 是否删除，0-否，1-是
	Gateway   string `json:"gateway" gorm:"gateway"`        // 网关地址
	AppSecret string `json:"app_secret"  gorm:"app_secret"` // 应用密钥
	Icon      string `json:"icon"        gorm:"icon"`       // 应用图标
	AppId     int32  `json:"app_id"      gorm:"app_id"`     // 对应的应用ID
}

// TableName 设置表名
func (CenterServiceAppDsns) TableName() string {
	return "app_dsns"
}

type CenterServiceAppDsnsFilter struct {
	Id       int32  // 通过应用ID精准查询
	AppKey   string // 通过 app_key 查
	XxxName  string // 游戏名称
	IsDelete *int8  // 查询是否删除的，为 nil 则表示不限制，否则按照指定的值查询
	AppId    int32  // 对应的应用ID
}
