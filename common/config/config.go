package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	APIKey       string `env:"APIKEY" envDefault:"demo"`
	PrivateKey   string `env:"PRIVATE_KEY" envDefault:"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"`
	GinMode      string `env:"GIN_MODE" envDefault:"debug"`
	Environment  string `env:"ENVIRONMENT" envDefault:"Dev"`
	KafkaHosts   string `env:"KAFKA_HOSTS" envDefault:"sync-streaming.chainbase.online:9093"`
	KafkaUser    string `env:"KAFKA_USER" envDefault:"testuser"`
	KafkaPass    string `env:"KAFKA_PASS" envDefault:"testuser"`
	ConsumeGroup string `env:"CONSUME_GROUP" `
}

func GetConf() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
	}
	return cfg
}
