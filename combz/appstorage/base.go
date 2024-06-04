package appstorage

type Storage interface {
	// 获取 实例 的 key
	getInstStatusKey(DSNCategory, string) string
	// 通过 app_key 检查是否连接
	checkConn(appKey string) bool
	// 添加通过 app_key 添加 实例连接
	addDSN(appKey string) error
	// 查询/调用rpc获取所有 dsn并存储到对向内
	allDSN() error
	// 移除 app_key 对应的 实例连接
	removeApp(appKey string) bool
	RemoveDB(app string) bool
}
