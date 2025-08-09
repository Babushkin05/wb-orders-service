package kafka

import (
	"context"

	"github.com/Babushkin05/wb-orders-service/internal/application"
	"github.com/Babushkin05/wb-orders-service/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type kafkaConfig struct {
	Broker  string `yaml:"broker"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
}

type InboxConsumer interface {
	Start(ctx context.Context)
}

type inboxConsumer struct {
	reader *kafka.Reader
	repo   application.OrdersRepository
}

func NewInboxConsumer(cfg kafkaConfig, repo application.OrdersRepository) InboxConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Broker},
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})

	return &inboxConsumer{
		reader: reader,
		repo:   repo,
	}
}

func (c *inboxConsumer) Start(ctx context.Context) {
	go func() {
		for {
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				logger.Log.Errorf("inbox consumer read error: %v", err)
				continue
			}

			err = c.repo.SaveInboxMessage(ctx, string(m.Key), m.Topic, string(m.Value))
			if err != nil {
				logger.Log.Errorf("failed to save inbox message: %v", err)
			}
		}
	}()
}
