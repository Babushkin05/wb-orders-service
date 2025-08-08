package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Babushkin05/wb-orders-service/internal/application"
	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewOrdersRepository(db *sqlx.DB) application.OrdersRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Get(orderUID string) (*model.Order, error) {
	ctx := context.Background()

	// Начинаем транзакцию
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Получаем основной заказ
	var order model.Order
	query := `SELECT * FROM orders WHERE order_uid = $1`
	err = tx.GetContext(ctx, &order, query, orderUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Получаем доставку
	var delivery model.Delivery
	query = `SELECT * FROM delivery WHERE order_uid = $1`
	err = tx.GetContext(ctx, &delivery, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	order.Delivery = delivery

	// Получаем платеж
	var payment model.Payment
	query = `SELECT * FROM payment WHERE transaction = $1`
	err = tx.GetContext(ctx, &payment, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	order.Payment = payment

	// Получаем товары
	var items []model.Item
	query = `SELECT * FROM items WHERE order_uid = $1`
	err = tx.SelectContext(ctx, &items, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	order.Items = items

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}

func (r *postgresRepository) Store(order *model.Order) error {
	ctx := context.Background()

	// Начинаем транзакцию
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Сохраняем основной заказ
	query := `INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES (
		:order_uid, :track_number, :entry, :locale, :internal_signature,
		:customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard
	)`
	_, err = tx.NamedExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Сохраняем доставку
	delivery := order.Delivery
	delivery.OrderUID = order.OrderUID
	query = `INSERT INTO delivery (
		order_uid, name, phone, zip, city, address, region, email
	) VALUES (
		:order_uid, :name, :phone, :zip, :city, :address, :region, :email
	)`
	_, err = tx.NamedExecContext(ctx, query, delivery)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	// Сохраняем платеж
	payment := order.Payment
	query = `INSERT INTO payment (
		transaction, request_id, currency, provider, amount, 
		payment_dt, bank, delivery_cost, goods_total, custom_fee
	) VALUES (
		:transaction, :request_id, :currency, :provider, :amount,
		:payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee
	)`
	_, err = tx.NamedExecContext(ctx, query, payment)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	// Сохраняем товары
	for _, item := range order.Items {
		item.OrderUID = order.OrderUID
		query = `INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid, name, 
			sale, size, total_price, nm_id, brand, status
		) VALUES (
			:order_uid, :chrt_id, :track_number, :price, :rid, :name,
			:sale, :size, :total_price, :nm_id, :brand, :status
		)`
		_, err = tx.NamedExecContext(ctx, query, item)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
