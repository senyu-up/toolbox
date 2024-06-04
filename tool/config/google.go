package config

type Google struct {
	ApiKey string `yaml:"apiKey"`
}

type Gcs struct {
	CredentialsJson string
	Bucket          string
}

type GcPubSub struct {
	CredentialsJson string
	ProjectId       string
}

type CloudTask struct {
	CredentialsJson string
	Region          string
}
