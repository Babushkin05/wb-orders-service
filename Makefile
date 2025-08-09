ifneq (,$(wildcard .env))
	include .env
	export
endif

DB_PORT ?= 5432
SERVER_PORT ?= 8080
DB_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR=./migrations

# Применить все миграции
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

# Откатить последнюю миграцию
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1

# Сбросить все миграции (удалить все таблицы)
migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" drop -f

docker-migrate-up:
	docker run --rm -v $(pwd)/migrations:/migrations --network=orders-service_default migrate/migrate -path=/migrations -database "postgres://postgres:postgres@postgres:5432/orders?sslmode=disable" up

docker-migrate-down:
	docker run --rm -v $(pwd)/migrations:/migrations --network=orders-service_default migrate/migrate -path=/migrations -database "postgres://postgres:postgres@postgres:5432/orders?sslmode=disable" down 1

docker-migrate-reset:
	docker run --rm -v $(pwd)/migrations:/migrations --network=orders-service_default migrate/migrate -path=/migrations -database "postgres://postgres:postgres@postgres:5432/orders?sslmode=disable" drop -f

# Создать новую миграцию: make migrate-new name=create_users
migrate-new:
ifndef name
	$(error "Usage: make migrate-new name=create_something")
endif
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)


up:
	docker-compose up --build

order-generator:
	go run cmd/order-generator/main.go
