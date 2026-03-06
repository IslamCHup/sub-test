APP_NAME=subscriptions-service

DB_URL=postgres://admin:12345@localhost:5432/subscriptions?sslmode=disable

.PHONY: run build test migrate-up migrate-down swagger tidy

run:
	go run cmd/app/main.go

build:
	go build -o bin/$(APP_NAME) cmd/app/main.go

test:
	go test ./...

tidy:
	go mod tidy

swagger:
	swag init -g cmd/app/main.go

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down
