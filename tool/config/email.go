package config

// aws ses 邮件服务配置
type EmailConfig struct {
	Region    string `yaml:"region"` // 地区
	AppId     string `yaml:"appid"`  // 应用Id
	AccessID  string `yaml:"AccessID"`
	AccessKey string `yaml:"AccessKey"`
}
