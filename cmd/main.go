package main

import (
	"log"

	"github.com/Babushkin05/wb-orders-service/internal/config"
	"github.com/Babushkin05/wb-orders-service/internal/infrastructure/postgres"
	"github.com/Babushkin05/wb-orders-service/internal/infrastructure/redis"
	"github.com/Babushkin05/wb-orders-service/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log.Printf("Loaded config: %+v\n", cfg)

	err := logger.Init(logger.Config{
		Level:  cfg.LoggerConfig.Level,
		Output: cfg.LoggerConfig.Output,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger.Log.Info("Logger initialized successfully")

	DBconn, err := postgres.NewDB(cfg.DataBase)

	db := postgres.NewOrdersRepository(DBconn)
	logger.Log.Info("DB initialized successfully")

	redis, err := redis.NewRedisCache(cfg.RedisConfig, DBconn)
	logger.Log.Info("Redis initialized successfully")

}
