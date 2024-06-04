package config

type Nsq struct {
	Lookup string `yaml:"lookup"`
	Nsqd   string `yaml:"nsqd"`
}
