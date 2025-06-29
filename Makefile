# Makefile for grpc-chat-app with Kafka

.PHONY: help kafka-up kafka-down server client proto clean

# 기본 목표
help:
	@echo "Available commands:"
	@echo "  kafka-up    - Start Kafka cluster with Docker Compose"
	@echo "  kafka-down  - Stop Kafka cluster"
	@echo "  server      - Run gRPC chat server"
	@echo "  client      - Run interactive chat client (specify USER=<username>)"
	@echo "  proto       - Generate protobuf code"
	@echo "  clean       - Clean up binaries and logs"
	@echo ""
	@echo "Example usage:"
	@echo "  make kafka-up"
	@echo "  make server"
	@echo "  make client USER=alice"
	@echo ""
	@echo "How to test chat:"
	@echo "  1. Start Kafka: make kafka-up"
	@echo "  2. Start server: make server (in terminal 1)"
	@echo "  3. Start clients: make client USER=alice (in terminal 2)"
	@echo "                   make client USER=bob (in terminal 3)"
	@echo "  4. Type messages in client terminals and see real-time chat!"

# Kafka 클러스터 시작
kafka-up:
	@echo "Starting Kafka cluster..."
	docker-compose up -d
	@echo "Waiting for Kafka to be ready..."
	sleep 10
	@echo "Kafka cluster is ready!"

# Kafka 클러스터 종료
kafka-down:
	@echo "Stopping Kafka cluster..."
	docker-compose down

# gRPC 서버 실행 (메시지 프로세서 내장)
server:
	@echo "Starting gRPC chat server with embedded message processor..."
	go run . server

# 채팅 클라이언트 실행
client:
	@if [ -z "$(USER)" ]; then \
		echo "Please specify USER: make client USER=<username>"; \
		echo "Example: make client USER=alice"; \
		exit 1; \
	fi
	@echo "Starting interactive chat client for user: $(USER)"
	@echo "Type messages and press Enter to send"
	@echo "Use '/quit' or Ctrl+C to exit"
	@echo "========================================"
	go run . client $(USER)

# Protobuf 코드 생성
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/chat.proto
	@echo "Protobuf code generated"

# 정리
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f *.log
	@echo "Cleanup completed" 