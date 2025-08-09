package main

import (
	"context"
	"log"
	netHttp "net/http"
	"strconv"

	"github.com/Babushkin05/wb-orders-service/internal/application"
	"github.com/Babushkin05/wb-orders-service/internal/config"
	"github.com/Babushkin05/wb-orders-service/internal/infrastructure/http"
	"github.com/Babushkin05/wb-orders-service/internal/infrastructure/kafka"
	"github.com/Babushkin05/wb-orders-service/internal/infrastructure/postgres"
	"github.com/Babushkin05/wb-orders-service/internal/infrastructure/redis"
	"github.com/Babushkin05/wb-orders-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	log.Printf("Loaded config: %+v\n", cfg)

	// Init logger
	err := logger.Init(logger.Config{
		Level:  cfg.LoggerConfig.Level,
		Output: cfg.LoggerConfig.Output,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger.Log.Info("Logger initialized successfully")

	// Init DB
	DBconn, err := postgres.NewDB(cfg.DataBase)
	if err != nil {
		logger.Log.Fatal("Failed to connect to DB: ", err)
	}
	db := postgres.NewOrdersRepository(DBconn)
	logger.Log.Info("DB initialized successfully")

	// Init cache
	redis, err := redis.NewRedisCache(cfg.RedisConfig, DBconn)
	if err != nil {
		logger.Log.Fatal("Failed to connect to Redis: ", err)
	}
	logger.Log.Info("Redis initialized successfully")

	// Init Kafka
	inboxConsumer := kafka.NewInboxConsumer(cfg.KafkaConfig, db)
	inboxProcessor := kafka.NewInboxProcessor(db)

	ctx := context.Background()
	inboxConsumer.Start(ctx)
	logger.Log.Info("Inbox consumer started successfully")

	inboxProcessor.Start(ctx)
	logger.Log.Info("Inbox processor started successfully")

	// Init service
	ordersService := application.NewOrdersService(redis, db)

	// Init HTTP server
	handler := http.NewHandler(ordersService)
	r := gin.Default()
	http.RegisterRoutes(r, handler)

	// Run server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	logger.Log.Info("Starting server")
	if err := r.Run(addr); err != nil && err != netHttp.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
