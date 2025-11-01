package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Doc2text struct {
	AccessTokenURL string `validate:"required" env:"ACCESS_TOKEN_URL"`
	ClientID       string `validate:"required" env:"CLIENT_ID_MESSAGE_RESPONDER_OCR"`
	ClientSecret   string `validate:"required" env:"CLIENT_SECRET_MESSAGE_RESPONDER_OCR"`
	GRpcUrl        string `validate:"required" env:"G_RPC_URL"`
}

type Kafka struct {
	BootstrapServersValue string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
	GroupID               string `validate:"required" env:"GROUP_ID_MESSAGE_RESPONDER_OCR"`
	RequestTopicName      string `validate:"required" env:"TOPIC_NAME_OCR_REQUEST"`
	ResponseTopicName     string `validate:"required" env:"TOPIC_NAME_TG_RESPONSE_PREPARER"`
	SaslUsername          string `env:"SASL_USERNAME"`
	SaslPassword          string `env:"SASL_PASSWORD"`
	ClientID              string `env:"CLIENT_ID_MESSAGE_RESPONDER_OCR"`
}
type Config struct {
	Kafka    Kafka    `envPrefix:"KAFKA_"`
	Doc2text Doc2text `envPrefix:"DOC3TEXT_"`
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
