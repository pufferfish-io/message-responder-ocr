package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Doc2text struct {
	AccessTokenURL string `validate:"required" env:"ACCESS_TOKEN_URL"`
	ClientID       string `validate:"required" env:"CLIENT_ID"`
	ClientSecret   string `validate:"required" env:"CLIENT_SECRET"`
	GRpcUrl        string `validate:"required" env:"G_RPC_URL"`
}

type Kafka struct {
	BootstrapServersValue string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
	GroupID               string `validate:"required" env:"GROUP_ID"`
	RequestTopicName      string `validate:"required" env:"REQUEST_TOPIC_NAME"`
	ResponseTopicName     string `validate:"required" env:"RESPONSE_TOPIC_NAME"`
	SaslUsername          string `env:"SASL_USERNAME"`
	SaslPassword          string `env:"SASL_PASSWORD"`
	ClientID              string `env:"CLIENT_ID"`
}
type Config struct {
	Kafka    Kafka    `envPrefix:"MSG_RESP_OCR_KAFKA_"`
	Doc2text Doc2text `envPrefix:"MSG_RESP_OCR_DOC3TEXT_"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}
	return &c, nil
}
