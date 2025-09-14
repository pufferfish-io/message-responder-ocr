package main

import (
	"context"
	"msg-responder-ocr/internal/config"
	"msg-responder-ocr/internal/logger"
	"msg-responder-ocr/internal/messaging"
	"msg-responder-ocr/internal/processor"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lg, cleanup := logger.NewZapLogger()
	defer cleanup()

	lg.Info("🚀 Starting msg-responder-ocr…")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		lg.Error("❌ Failed to load config: %v", err)
		os.Exit(1)
	}

	producer, err := messaging.NewKafkaProducer(messaging.Option{
		Logger:       lg,
		Broker:       cfg.Kafka.BootstrapServersValue,
		SaslUsername: cfg.Kafka.SaslUsername,
		SaslPassword: cfg.Kafka.SaslPassword,
		ClientID:     cfg.Kafka.ClientID,
	})
	if err != nil {
		lg.Error("❌ Failed to create producer: %v", err)
		os.Exit(1)
	}
	defer producer.Close()

	responder := processor.NewMessageResponder(processor.Option{
		Doc2textURL: cfg.Doc2text.GRpcUrl,
		KafkaTopic:  cfg.Kafka.ResponseTopicName,
		Producer:    producer,
		Logger:      lg,
		Auth: processor.AuthOptions{
			AccessTokenURL: cfg.Doc2text.AccessTokenURL,
			ClientID:       cfg.Doc2text.ClientID,
			ClientSecret:   cfg.Doc2text.ClientSecret,
		},
	})

	consumer, err := messaging.NewKafkaConsumer(messaging.ConsumerOption{
		Logger:       lg,
		Broker:       cfg.Kafka.BootstrapServersValue,
		GroupID:      cfg.Kafka.GroupID,
		Topics:       []string{cfg.Kafka.RequestTopicName},
		Handler:      responder,
		SaslUsername: cfg.Kafka.SaslUsername,
		SaslPassword: cfg.Kafka.SaslPassword,
		ClientID:     cfg.Kafka.ClientID,
	})

	if err != nil {
		lg.Error("❌ Failed to create consumer: %v", err)
		os.Exit(1)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	if err := consumer.Start(ctx); err != nil {
		lg.Error("❌ Consumer error: %v", err)
		os.Exit(1)
	}
}
