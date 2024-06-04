package config

type Aws struct {
	AwsAccessId  string `yaml:"awsAccessId"`
	AwsAccessKey string `yaml:"awsAccessKey"`

	S3 []S3Storage `yaml:"s3"` // 单个regin下的配置
}

type S3Storage struct {
	//地区
	Region string `yaml:"region"`
	//桶 名称
	Bucket string `yaml:"bucket"`
	//文件存储位置 s3 key拼接用
	Path string `yaml:"path"`
	//过期时间 单位 小时 (The expire parameter is only used for presigned Amazon S3 API requests)
	Expire int64 `yaml:"expire"`
	//下载域名拼接
	Host string `yaml:"host"`
}
