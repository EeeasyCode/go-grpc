# 🚀 Go + gRPC + Kafka 실시간 채팅 시스템

Go, gRPC, Apache Kafka를 활용한 확장 가능한 실시간 채팅 시스템입니다.

## ✨ 주요 특징

- **gRPC 스트리밍**: 양방향 실시간 통신
- **Apache Kafka**: 분산 메시지 브로커로 확장성 보장
- **고가용성**: 3개 브로커 Kafka 클러스터
- **동시성**: Go의 Goroutine을 활용한 병렬 처리
- **분산 아키텍처**: 여러 서버 인스턴스 지원

## 🏗️ 시스템 아키텍처

```
gRPC Clients ←→ gRPC Server ←→ Kafka Cluster (3 Brokers)
                     ↓
             Message Processor (내장)
```

**핵심 컴포넌트:**

- **gRPC Server**: 클라이언트 연결 관리 & 메시지 게이트웨이
- **Kafka Cluster**: 분산 메시지 브로커 (3 브로커, 고가용성)
- **Message Processor**: Kafka 메시지 소비 & 분산 처리 (서버에 내장)

## 🚀 빠른 시작

### 1. 사전 요구사항

- Go 1.19+
- Docker & Docker Compose

### 2. 설치 & 실행

```bash
# 1. 의존성 설치
go mod tidy

# 2. Kafka 클러스터 시작
make kafka-up

# 3. gRPC 서버 시작 (터미널 1)
make server

# 4. 클라이언트 실행 (터미널 2, 3, ...)
make client USER=alice
make client USER=bob
```

### 3. 채팅 테스트

1. 각 클라이언트 터미널에서 메시지 입력
2. 모든 연결된 사용자에게 실시간 전달 확인
3. `/quit` 또는 `Ctrl+C`로 종료

## 📝 사용 방법

### 기본 명령어

```bash
# 도움말
make help

# Kafka 관리
make kafka-up    # Kafka 클러스터 시작
make kafka-down  # Kafka 클러스터 종료

# 서버 실행
make server      # gRPC 서버 (Message Processor 내장)

# 클라이언트 실행
make client USER=<사용자명>

# Protobuf 코드 생성
make proto

# 정리
make clean
```

### 직접 실행

```bash
# 서버 모드
go run . server

# 클라이언트 모드
go run . client <사용자명>

# 컨슈머 모드 (별도 실행시)
go run . consumer
```

## 🔧 구성 요소

### Kafka 토픽

- **`chatting`**: 채팅 메시지 (2 파티션, 2 복제본)
- **`user-connections`**: 사용자 연결 이벤트 (3 파티션, 2 복제본)

### 포트 설정

- gRPC Server: `8081`
- Kafka Brokers: `9092`, `9093`, `9094`
- Kafka UI: `8080` (http://localhost:8080)

## 📊 모니터링

### Kafka UI

- **URL**: http://localhost:8080
- **기능**: 토픽, 메시지, 컨슈머 그룹 실시간 모니터링

### 로그 확인

```bash
# 성공적인 연결
User alice connected to server grpc-server-xxxxx

# 메시지 전송 성공
Message sent to Kafka - Topic: chatting, Partition: 1, Offset: 127

# 메시지 수신 성공
Sent message to bob from alice
```

## 🚨 문제 해결

### 일반적인 문제

**1. Kafka 연결 실패**

```bash
# Kafka 상태 확인
docker-compose ps

# 재시작
make kafka-down && make kafka-up
```

**2. 포트 충돌**

```bash
# 포트 사용 프로세스 확인 및 종료
lsof -ti:8081 | xargs kill -9
```

**3. 메시지 전달 안됨**

- 서버가 Message Processor를 내장하고 있으므로 별도 consumer 실행 불필요
- 서버만 실행하면 자동으로 메시지 처리

### 디버깅 팁

```bash
# Docker 로그 확인
docker-compose logs kafka-01

# 서버 상세 로그 확인
go run . server 2>&1 | tee server.log
```

## 🔄 확장성

### 수평 확장

- **여러 gRPC 서버**: 로드 밸런서와 함께 사용
- **Kafka 파티션 증가**: 처리량 향상
- **Consumer Group**: 자동 부하 분산

### 성능 최적화

- **배치 처리**: Kafka Producer 배치 설정
- **연결 풀링**: 효율적인 연결 관리
- **비동기 처리**: Goroutine 기반 병렬 처리

## 📚 기술 스택

- **언어**: Go 1.19+
- **통신**: gRPC (Protocol Buffers)
- **메시지 브로커**: Apache Kafka
- **컨테이너**: Docker & Docker Compose
- **라이브러리**: Sarama (Kafka Go Client)

## 🔗 참고 자료

- [gRPC 공식 문서](https://grpc.io/docs/)
- [Apache Kafka 문서](https://kafka.apache.org/documentation/)
- [Sarama - Go Kafka Client](https://github.com/IBM/sarama)

---

**💡 Tip**: Kafka UI (http://localhost:8080)에서 실시간으로 메시지 흐름을 모니터링할 수 있습니다!
