package kafka

import (
	"context"
	"time"

	"github.com/Babushkin05/wb-orders-service/internal/application"
	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
	"github.com/Babushkin05/wb-orders-service/pkg/logger"
)

type InboxProcessor interface {
	Start(ctx context.Context)
	processBatch(ctx context.Context) error
}

type inboxProcessor struct {
	repo application.OrdersRepository
}

func NewInboxProcessor(repo application.OrdersRepository) InboxProcessor {
	return &inboxProcessor{
		repo: repo,
	}
}

func (p *inboxProcessor) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				if err := p.processBatch(ctx); err != nil {
					logger.Log.Errorf("inbox processor error: %v", err)
				}
			case <-ctx.Done():
				logger.Log.Info("inbox processor stopped")
				return
			}
		}
	}()
}

func (p *inboxProcessor) processBatch(ctx context.Context) error {
	msgs, err := p.repo.FetchUnprocessedInboxMessages(ctx, 10)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		order, err := model.UnmarshalOrder([]byte(msg.Payload))
		if err != nil {
			logger.Log.Errorf("failed to unmarshal order: %v", err)
			if err := p.repo.MarkInboxMessageProcessed(ctx, msg.ID); err != nil {
				logger.Log.Errorf(
					"failed to mark inbox message as processed: %v",
					err,
				)
				return err
			}
		}

		if err := p.repo.Store(order); err != nil {
			logger.Log.Errorf("failed to store order: %v", err)
			if err := p.repo.MarkInboxMessageProcessed(ctx, msg.ID); err != nil {
				logger.Log.Errorf(
					"failed to mark inbox message as processed: %v",
					err,
				)
				return err
			}
		}
	}

	return nil
}
