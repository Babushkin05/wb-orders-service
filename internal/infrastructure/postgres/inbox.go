package postgres

import (
	"context"
	"time"

	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
)

func (r *postgresRepository) SaveInboxMessage(ctx context.Context, messageID, topic, payload string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO inbox (message_id, topic, payload, created_at, processed)
		VALUES ($1, $2, $3, $4, false)
		ON CONFLICT (message_id) DO NOTHING
	`, messageID, topic, payload, time.Now())

	return err
}

func (r *postgresRepository) FetchUnprocessedInboxMessages(ctx context.Context, limit int) ([]model.InboxMessage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT message_id, topic, payload
		FROM inbox
		WHERE processed = false
		ORDER BY created_at
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []model.InboxMessage
	for rows.Next() {
		var m model.InboxMessage
		if err := rows.Scan(&m.ID, &m.Topic, &m.Payload); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}

func (r *postgresRepository) MarkInboxMessageProcessed(ctx context.Context, messageID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE inbox SET processed = true WHERE message_id = $1
	`, messageID)
	return err
}
