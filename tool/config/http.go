package config

// FiberConfig
// http 框架 fiber 配置
type FiberConfig struct {
	Name                    string `yaml:"-"`                       // http 服务名
	Addr                    string `yaml:"addr"`                    // http 监听地址, 例：0.0.0.0:80
	CaseSensitive           bool   `yaml:"caseSensitive"`           // 路由是否大小写敏感
	Timeout                 int64  `yaml:"timeout"`                 // 请求处理时长, 单位 毫秒（ms）
	BodyLimit               int    `yaml:"bodyLimit"`               // 返回的数据大小限制， 单位：字节， 默认 : 4 * 1024 * 1024（4MB）
	EnablePrintRoutes       bool   `yaml:"enablePrintRoutes"`       // 启动时是否打印路由信息
	EnableTrustedProxyCheck bool   `yaml:"enableTrustedProxyCheck"` // 是否启用代理检查
	//Pprof             bool   `yaml:"pprof"`             // 是否开启pprof, 不能在业务服务上开启fiber路由，暴露pprof通过healthChecker来实现
}

// GinConfig http框架 gin配置
type GinConfig struct {
	Name                    string `yaml:"-"`                       // http 服务名
	Addr                    string `yaml:"addr"`                    // http 监听地址, 例：0.0.0.0:80
	CaseSensitive           bool   `yaml:"caseSensitive"`           // 路由是否大小写敏感
	Timeout                 int64  `yaml:"timeout"`                 // 请求处理时长, 单位 毫秒（ms）
	BodyLimit               int    `yaml:"bodyLimit"`               // 返回的数据大小限制， 单位：字节， 默认 : 4 * 1024 * 1024（4MB）
	EnablePrintRoutes       bool   `yaml:"enablePrintRoutes"`       // 启动时是否打印路由信息
	EnableTrustedProxyCheck bool   `yaml:"enableTrustedProxyCheck"` // 是否启用代理检查
}

type HealthCheck struct {
	Pprof      bool   `yaml:"pprof"`      // 是否开启pprof
	Addr       string `yaml:"addr"`       // 监听地址
	Port       uint32 `yaml:"port"`       // 端口
	DisableLog bool   `yaml:"disableLog"` // 是否禁用日志, 禁用后访问 /system/health 时不会打印日志
}
