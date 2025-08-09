package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Babushkin05/wb-orders-service/internal/application"
	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type redisCache struct {
	client *redis.Client
	db     *sqlx.DB
	ttl    time.Duration
}

type RedisConfig struct {
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Password string        `yaml:"password"`
	DB       int           `yaml:"db"`
	TTL      time.Duration `yaml:"ttl"`
}

func NewRedisCache(cfg RedisConfig, db *sqlx.DB) (application.Cacher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	cache := &redisCache{
		client: client,
		db:     db,
		ttl:    cfg.TTL,
	}

	if err := cache.WarmUp(); err != nil {
		return nil, fmt.Errorf("cache warmup failed: %w", err)
	}

	return cache, nil
}

func (r *redisCache) WarmUp() error {
	ctx := context.Background()

	// Получаем только 1000 последних заказов
	var orders []model.Order
	query := `SELECT * FROM orders ORDER BY date_created DESC LIMIT 1000`
	err := r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return fmt.Errorf("db query error: %w", err)
	}

	for _, order := range orders {
		if err := r.Cache(&order); err != nil {
			return fmt.Errorf("caching error: %w", err)
		}
	}

	return nil
}

func (r *redisCache) Cache(order *model.Order) error {
	ctx := context.Background()

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := r.client.Set(ctx, order.OrderUID.String(), data, r.ttl).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

func (r *redisCache) GetOrderFromCache(orderUID string) (model.Order, error) {
	ctx := context.Background()

	data, err := r.client.Get(ctx, orderUID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return model.Order{}, nil
		}
		return model.Order{}, fmt.Errorf("redis get error: %w", err)
	}

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return model.Order{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return order, nil
}
