package config

type ImageAudit struct {
	SecretId  string `yaml:"secretid"`
	SecretKey string `yaml:"secretkey"`

	Bucket string `yaml:"bucket"` // 桶
	Regin  string `yaml:"regin"`  // 地区

	BizType string `yaml:"biztype"` // 业务类型
}
