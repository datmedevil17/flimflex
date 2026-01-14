.PHONY: proto proto-auth proto-user proto-movie proto-upload proto-streaming proto-recommendation proto-notification

# Generate all proto files
proto: proto-auth proto-user proto-movie proto-upload proto-streaming proto-recommendation

# Generate individual proto files
proto-auth:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/auth.proto

proto-user:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto

proto-movie:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/movie/movie.proto

proto-upload:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/upload/upload.proto

proto-streaming:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/streaming/streaming.proto

proto-recommendation:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/recommendation/recommendation.proto

# Docker commands
build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

restart:
	docker compose restart

clean:
	docker compose down -v

# Service specific commands
auth-build:
	cd services/auth && docker build -t netflix-auth-service -f docker/Dockerfile .

user-build:
	cd services/user && docker build -t netflix-user-service -f docker/Dockerfile .

movie-build:
	cd services/movie && docker build -t netflix-movie-service -f docker/Dockerfile .

upload-build:
	cd services/upload && docker build -t netflix-upload-service -f docker/Dockerfile .

streaming-build:
	cd services/streaming && docker build -t netflix-streaming-service -f docker/Dockerfile .

recommendation-build:
	cd services/recommendation && docker build -t netflix-recommendation-service -f docker/Dockerfile .



gateway-build:
	cd api-gateway && docker build -t netflix-api-gateway -f docker/Dockerfile .

# Development commands
dev-auth:
	cd services/auth && go run cmd/main.go

dev-user:
	cd services/user && go run cmd/main.go

dev-movie:
	cd services/movie && go run cmd/main.go

dev-upload:
	cd services/upload && go run cmd/main.go

dev-streaming:
	cd services/streaming && go run cmd/main.go

dev-recommendation:
	cd services/recommendation && go run cmd/main.go



dev-gateway:
	cd api-gateway && go run cmd/main.go

# Database commands
db-shell:
	docker exec -it netflix-postgres psql -U postgres

db-reset:
	docker-compose down -v
	docker-compose up -d postgres
	@echo "Databases recreated from init script"

# Test commands
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Install dependencies
install-protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Lint
lint:
	golangci-lint run ./...

# Format
fmt:
	go fmt ./...

# Help
help:
	@echo "Available commands:"
	@echo "  make proto                 - Generate all proto files"
	@echo "  make build                 - Build all Docker images"
	@echo "  make up                    - Start all services"
	@echo "  make down                  - Stop all services"
	@echo "  make logs                  - View logs"
	@echo "  make restart               - Restart all services"
	@echo "  make clean                 - Stop and remove all containers and volumes"
	@echo "  make dev-<service>         - Run service locally"
	@echo "  make test                  - Run tests"
	@echo "  make lint                  - Run linter"
	@echo "  make fmt                   - Format code"