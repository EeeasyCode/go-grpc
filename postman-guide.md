# Postman으로 gRPC 채팅 애플리케이션 테스트하기

이 가이드는 Postman을 사용하여 gRPC 채팅 애플리케이션을 테스트하는 방법을 설명합니다.

## 🛠️ 사전 준비

### 1. Postman 설치 및 버전 확인

- **Postman 최신 버전** (v9.0 이상) 설치 필요
- gRPC 지원이 포함된 버전이어야 함

### 2. 서버 실행

```bash
# 1. Kafka 클러스터 시작
make kafka-up

# 2. gRPC 서버 시작 (Reflection 활성화됨)
make server

# 3. Consumer 시작 (선택사항 - 메시지 처리 확인용)
make consumer
```

## 📡 Postman에서 gRPC 설정

### 방법 1: Server Reflection 사용 (추천)

1. **새 gRPC 요청 생성**

   - Postman에서 "New" → "gRPC Request" 선택

2. **서버 URL 입력**

   ```
   localhost:8081
   ```

3. **Service Discovery**

   - "Use Server Reflection" 옵션 활성화
   - "Refresh" 버튼 클릭하여 서비스 목록 가져오기

4. **사용 가능한 서비스**
   - `chat.Broadcast` 서비스가 나타남
   - 두 개의 메서드 확인:
     - `CreateStream`
     - `BroadcastMessage`

### 방법 2: Proto 파일 직접 업로드

1. **Proto 파일 선택**

   - "Select proto file" 옵션 선택
   - `proto/chat.proto` 파일 업로드

2. **Import 설정**
   - Import paths 설정 (필요시)

## 🧪 테스트 시나리오

### 시나리오 1: 메시지 브로드캐스트 테스트

#### 1. BroadcastMessage 호출

**설정:**

- **Service**: `chat.Broadcast`
- **Method**: `BroadcastMessage`
- **Server URL**: `localhost:8081`

**요청 본문 (JSON):**

```json
{
  "id": "alice",
  "content": "Hello from Postman!",
  "timestamp": {
    "seconds": 1640995200,
    "nanos": 0
  }
}
```

**응답 예상:**

```json
{}
```

#### 2. 로그 확인

서버 콘솔에서 다음과 같은 로그 확인:

```
Message sent to Kafka - Topic: chatting, Partition: 1, Offset: 0
```

### 시나리오 2: 클라이언트 연결 스트림 테스트

#### 1. CreateStream 호출

**설정:**

- **Service**: `chat.Broadcast`
- **Method**: `CreateStream`
- **Server URL**: `localhost:8081`

**요청 본문 (JSON):**

```json
{
  "user": {
    "id": "postman-user",
    "name": "Postman Test User"
  },
  "active": true
}
```

#### 2. 스트림 응답 확인

- 연결이 성공하면 스트림이 열림
- 다른 클라이언트가 메시지를 보내면 실시간으로 수신

## 📋 테스트 케이스

### 테스트 케이스 1: 기본 메시지 전송

```json
// 요청: BroadcastMessage
{
  "id": "test-user-1",
  "content": "First test message",
  "timestamp": null
}

// 기대 결과: 성공 응답 + Kafka에 메시지 발행
```

### 테스트 케이스 2: 특수 문자 포함 메시지

```json
// 요청: BroadcastMessage
{
  "id": "test-user-2",
  "content": "Special chars: 한글, emoji 😀, symbols @#$%",
  "timestamp": null
}

// 기대 결과: 특수 문자가 정상적으로 처리됨
```

### 테스트 케이스 3: 긴 메시지 전송

```json
// 요청: BroadcastMessage
{
  "id": "test-user-3",
  "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
  "timestamp": null
}

// 기대 결과: 긴 메시지도 정상 처리
```

### 테스트 케이스 4: 동시 연결 테스트

1. **첫 번째 CreateStream 연결**

```json
{
  "user": {
    "id": "user-1",
    "name": "First User"
  },
  "active": true
}
```

2. **두 번째 CreateStream 연결 (새 탭)**

```json
{
  "user": {
    "id": "user-2",
    "name": "Second User"
  },
  "active": true
}
```

3. **메시지 전송 테스트**

```json
{
  "id": "broadcast-test",
  "content": "Message to all connected users",
  "timestamp": null
}
```

## 🔍 모니터링 및 디버깅

### 1. Kafka UI 모니터링

- **URL**: http://localhost:8080
- **확인 사항**:
  - `chatting` 토픽의 메시지 수신 여부
  - `user-connections` 토픽의 연결 이벤트

### 2. 서버 로그 확인

```bash
# 서버 실행 시 다음 로그들을 확인:
User postman-user connected to server grpc-server-1234567890
Message sent to Kafka - Topic: chatting, Partition: 0, Offset: 1
Published user connection event: postman-user (connected: true)
```

### 3. Consumer 로그 확인

```bash
# Consumer 실행 시 다음 로그들을 확인:
Received chat message: Topic=chatting, Partition=0, Offset=1
Processing chat message: ID=alice, Content=Hello from Postman!, From=alice
```

## ⚠️ 일반적인 문제 해결

### 1. "Service Unavailable" 오류

```bash
# 서버 상태 확인
curl -v http://localhost:8081
# 또는
telnet localhost 8081
```

**해결방법:**

- 서버가 실행 중인지 확인
- 포트 8081이 사용 가능한지 확인

### 2. "Failed to connect" 오류

**확인사항:**

- gRPC 서버가 실행 중인지
- Reflection이 활성화되어 있는지
- 방화벽 설정

### 3. "No services found" 오류

**해결방법:**

- Server Reflection 대신 Proto 파일 업로드 시도
- 서버 재시작 후 재시도

### 4. 스트림 연결이 즉시 종료되는 경우

**원인:**

- 클라이언트가 연결을 유지하지 않음
- 서버 측 에러 발생

**해결방법:**

- 서버 로그에서 에러 메시지 확인
- Kafka 연결 상태 확인

## 🎯 고급 테스트 시나리오

### 1. 부하 테스트

- Postman Runner 사용
- 여러 BroadcastMessage 요청을 연속 실행
- 응답 시간 및 성공률 모니터링

### 2. 장시간 연결 테스트

- CreateStream 연결을 장시간 유지
- 주기적으로 메시지 전송하여 연결 안정성 확인

### 3. 오류 조건 테스트

```json
// 잘못된 요청 형식
{
  "invalid_field": "test"
}

// 빈 메시지
{
  "id": "",
  "content": "",
  "timestamp": null
}
```

## 📝 성공 기준

### BroadcastMessage 테스트 성공

- [x] HTTP 200 응답 수신
- [x] 빈 응답 본문 (`{}`) 수신
- [x] 서버 로그에 Kafka 발행 성공 메시지
- [x] Kafka UI에서 메시지 확인 가능

### CreateStream 테스트 성공

- [x] 스트림 연결 성공
- [x] 서버 로그에 사용자 연결 메시지
- [x] 다른 클라이언트 메시지 실시간 수신
- [x] 연결 해제 시 정상적인 정리

이제 Postman을 사용하여 gRPC 채팅 애플리케이션을 완전히 테스트할 수 있습니다!
