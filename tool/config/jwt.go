package config

type Jwt struct {
	TokenSecret     string `yaml:"tokenSecret"`
	TokenExpiration int64  `yaml:"tokenExpiration"`
}
