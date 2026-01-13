
.PHONY: proto build run

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/*.proto

build:
	go build -o bin/api-gateway ./api-gateway
	go build -o bin/auth-service ./services/auth
	go build -o bin/user-service ./services/user
	go build -o bin/movie-service ./services/movie
	go build -o bin/upload-service ./services/upload
	go build -o bin/streaming-service ./services/streaming
	go build -o bin/recommendation-service ./services/recommendation
	go build -o bin/notification-service ./services/notification

run: build
	./bin/api-gateway &
	./bin/auth-service &
	./bin/user-service &
	./bin/movie-service &
	./bin/upload-service &
	./bin/streaming-service &
	./bin/recommendation-service &
	./bin/notification-service &

docker-up:
	docker-compose up --build
