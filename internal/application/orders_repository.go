package application

import (
	"context"

	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
)

type OrdersRepository interface {
	Get(orderUID string) (*model.Order, error)
	Store(model *model.Order) error
	SaveInboxMessage(ctx context.Context, messageID, topic, payload string) error
	FetchUnprocessedInboxMessages(ctx context.Context, limit int) ([]model.InboxMessage, error)
	MarkInboxMessageProcessed(ctx context.Context, messageID string) error
}
