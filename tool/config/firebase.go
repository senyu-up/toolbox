package config

type Firebase struct {
	// Credentials CredentialsFile 二选一
	// CredentialsFile 优先级高于 CredentialsJson
	// CredentialsJson 需要经过base64编码
	CredentialsJson string
	CredentialsFile string
	ProjectId       string
	AccountId       string
	Bucket          string
	DatabaseURL     string
}
