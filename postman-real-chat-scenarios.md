# 실제 채팅 서비스 Postman 테스트 시나리오

이 문서는 실제 채팅 서비스를 운영한다고 가정하고 작성된 포괄적인 Postman 테스트 시나리오입니다.

## 🎯 테스트 목표

- **사용자 경험 검증**: 실제 사용자가 겪을 수 있는 모든 상황 테스트
- **성능 및 안정성**: 동시 접속, 대용량 메시지 처리
- **비즈니스 로직**: 채팅방, 알림, 메시지 상태 관리
- **예외 상황 처리**: 네트워크 장애, 서버 오류 등

## 📋 테스트 환경 설정

### 사전 준비

```bash
# 1. 전체 환경 시작
make kafka-up
sleep 15
make server  # 터미널 1
make consumer  # 터미널 2

# 2. Postman Collection 생성
# "Real Chat Service Tests" 컬렉션 생성
```

### 테스트 데이터

```javascript
// Postman Environment Variables
{
  "server_url": "localhost:8081",
  "current_user": "alice",
  "chat_partner": "bob",
  "group_name": "dev-team",
  "test_message_count": 10
}
```

## 🧪 상세 테스트 시나리오

### 시나리오 1: 신규 사용자 온보딩

#### 1.1 첫 접속 사용자

```json
// CreateStream: 신규 사용자 첫 연결
{
  "user": {
    "id": "{{$randomFirstName}}_{{$randomInt}}",
    "name": "{{$randomFullName}}"
  },
  "active": true
}

// 검증 포인트:
// - 연결 성공
// - 서버 로그에 사용자 등록 확인
// - Kafka user-connections 토픽에 이벤트 발행
```

#### 1.2 환영 메시지 전송

```json
// BroadcastMessage: 시스템 환영 메시지
{
  "id": "system",
  "content": "{{current_user}}님, 채팅 서비스에 오신 것을 환영합니다! 🎉",
  "timestamp": null
}
```

### 시나리오 2: 1:1 채팅

#### 2.1 두 사용자 동시 접속

```json
// User A 연결 (Postman Tab 1)
{
  "user": {
    "id": "alice",
    "name": "Alice Johnson"
  },
  "active": true
}

// User B 연결 (Postman Tab 2)
{
  "user": {
    "id": "bob",
    "name": "Bob Smith"
  },
  "active": true
}
```

#### 2.2 실시간 대화 시뮬레이션

```json
// Alice가 첫 메시지 전송
{
  "id": "alice",
  "content": "안녕하세요 Bob! 프로젝트 관련해서 얘기하고 싶은 게 있어요.",
  "timestamp": null
}

// 3초 대기 후 Bob 응답
{
  "id": "bob",
  "content": "안녕하세요 Alice! 네, 말씀하세요.",
  "timestamp": null
}

// Alice 상세 설명
{
  "id": "alice",
  "content": "새로운 채팅 시스템 도입을 검토 중인데요, 어떻게 생각하시나요?",
  "timestamp": null
}
```

### 시나리오 3: 그룹 채팅

#### 3.1 다중 사용자 접속 (5명)

```json
// 각각 별도 Postman 탭에서 실행

// User 1: Project Manager
{
  "user": {
    "id": "pm_sarah",
    "name": "Sarah (PM)"
  },
  "active": true
}

// User 2: Frontend Developer
{
  "user": {
    "id": "fe_mike",
    "name": "Mike (Frontend)"
  },
  "active": true
}

// User 3: Backend Developer
{
  "user": {
    "id": "be_jane",
    "name": "Jane (Backend)"
  },
  "active": true
}

// User 4: Designer
{
  "user": {
    "id": "ui_tom",
    "name": "Tom (Designer)"
  },
  "active": true
}

// User 5: QA Engineer
{
  "user": {
    "id": "qa_lisa",
    "name": "Lisa (QA)"
  },
  "active": true
}
```

#### 3.2 팀 회의 시뮬레이션

```json
// PM이 회의 시작
{
  "id": "pm_sarah",
  "content": "👋 안녕하세요 팀! 오늘 스프린트 리뷰 시간입니다.",
  "timestamp": null
}

// 순차적 응답 (각 사용자별 3초 간격)
{
  "id": "fe_mike",
  "content": "프론트엔드 작업 90% 완료되었습니다!",
  "timestamp": null
}

{
  "id": "be_jane",
  "content": "백엔드 API 개발 완료, 테스트 진행 중입니다.",
  "timestamp": null
}

{
  "id": "ui_tom",
  "content": "새로운 UI 컴포넌트 디자인 완료했어요. 검토 부탁드립니다.",
  "timestamp": null
}

{
  "id": "qa_lisa",
  "content": "현재까지 발견된 버그 3개, 모두 낮은 우선순위입니다.",
  "timestamp": null
}
```

### 시나리오 4: 메시지 타입 다양성 테스트

#### 4.1 일반 텍스트 메시지

```json
{
  "id": "alice",
  "content": "안녕하세요! 오늘 날씨가 정말 좋네요.",
  "timestamp": null
}
```

#### 4.2 이모지 및 특수문자

```json
{
  "id": "bob",
  "content": "정말요? 😊🌞 저도 산책하고 싶어요! ☀️🚶‍♂️",
  "timestamp": null
}
```

#### 4.3 긴 메시지 (소설 한 단락)

```json
{
  "id": "alice",
  "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium.",
  "timestamp": null
}
```

#### 4.4 코드 블록 시뮬레이션

````json
{
  "id": "fe_mike",
  "content": "```javascript\nfunction sendMessage(content) {\n  return fetch('/api/messages', {\n    method: 'POST',\n    body: JSON.stringify({ content })\n  });\n}\n```",
  "timestamp": null
}
````

#### 4.5 다국어 메시지

```json
{
  "id": "international_user",
  "content": "Hello! 안녕하세요! こんにちは! Bonjour! Hola! Привет! مرحبا! 你好!",
  "timestamp": null
}
```

### 시나리오 5: 고급 사용 사례

#### 5.1 링크 및 멘션 시뮬레이션

```json
{
  "id": "alice",
  "content": "@bob 이 링크 확인해보세요: https://github.com/example/chat-app. 정말 유용한 자료인 것 같아요!",
  "timestamp": null
}
```

#### 5.2 시스템 공지사항

```json
{
  "id": "system_admin",
  "content": "🔔 [시스템 공지] 오늘 오후 3시에 10분간 서비스 점검이 있을 예정입니다. 이용에 불편을 드려 죄송합니다.",
  "timestamp": null
}
```

#### 5.3 파일 업로드 시뮬레이션

```json
{
  "id": "ui_tom",
  "content": "📎 파일 첨부: design_mockup_v2.pdf (2.3MB)\n프로젝트 최신 디자인 목업입니다. 피드백 부탁드려요!",
  "timestamp": null
}
```

### 시나리오 6: 부하 및 성능 테스트

#### 6.1 연속 메시지 전송 (Postman Runner 사용)

```json
// 반복 설정: 20회
{
  "id": "stress_test_user",
  "content": "부하 테스트 메시지 #{{$iteration}}",
  "timestamp": null
}

// Pre-request Script:
setTimeout(() => {}, 100); // 100ms 간격
```

#### 6.2 대용량 메시지

```json
{
  "id": "large_data_user",
  "content": "{{$randomLoremParagraphs(10)}}",
  "timestamp": null
}
```

### 시나리오 7: 예외 상황 테스트

#### 7.1 빈 메시지 전송

```json
{
  "id": "test_user",
  "content": "",
  "timestamp": null
}
```

#### 7.2 특수한 ID 패턴

```json
{
  "id": "user@email.com",
  "content": "이메일 형태의 사용자 ID 테스트",
  "timestamp": null
}

{
  "id": "user-with-special-chars!@#$%",
  "content": "특수문자가 포함된 사용자 ID 테스트",
  "timestamp": null
}
```

#### 7.3 잘못된 형식의 요청

```json
// 잘못된 필드명
{
  "user_id": "wrong_field",
  "message": "잘못된 필드명 테스트",
  "time": null
}

// 필수 필드 누락
{
  "content": "ID 없는 메시지 테스트"
}
```

### 시나리오 8: 실시간성 검증

#### 8.1 메시지 순서 보장 테스트

```json
// 5초 간격으로 순차 전송
{
  "id": "sequence_test",
  "content": "메시지 순서 테스트 1/5",
  "timestamp": null
}
// ... 2, 3, 4, 5번째 메시지
```

#### 8.2 동시 전송 테스트

```json
// 여러 사용자가 동시에 메시지 전송 (Postman Runner로 동시 실행)
{
  "id": "user_{{$randomInt}}",
  "content": "동시 전송 테스트 - {{$timestamp}}",
  "timestamp": null
}
```

### 시나리오 9: 비즈니스 로직 테스트

#### 9.1 고객 지원 채팅 시뮬레이션

```json
// 고객
{
  "id": "customer_john",
  "content": "안녕하세요, 결제 관련 문의가 있습니다.",
  "timestamp": null
}

// 상담원
{
  "id": "support_agent_1",
  "content": "안녕하세요! 어떤 결제 문제가 있으신가요? 자세히 설명해 주시면 도와드리겠습니다.",
  "timestamp": null
}

// 고객 상세 설명
{
  "id": "customer_john",
  "content": "어제 결제했는데 승인은 됐지만 서비스가 활성화되지 않았어요. 주문번호는 ORD-20231210-001입니다.",
  "timestamp": null
}
```

#### 9.2 교육 플랫폼 채팅 시뮬레이션

```json
// 강사
{
  "id": "instructor_kim",
  "content": "📚 오늘의 수업: gRPC와 Kafka를 활용한 실시간 채팅 시스템 구축",
  "timestamp": null
}

// 학생들 질문
{
  "id": "student_01",
  "content": "선생님, gRPC와 WebSocket의 차이점이 궁금합니다.",
  "timestamp": null
}

{
  "id": "student_02",
  "content": "Kafka Consumer Group은 어떻게 동작하나요?",
  "timestamp": null
}
```

## 📊 성능 지표 및 검증 기준

### 응답 시간 기준

- **메시지 전송**: < 100ms
- **연결 설정**: < 500ms
- **동시 접속 100명**: 응답 시간 < 200ms

### 처리량 기준

- **초당 메시지 처리**: 1000+ messages/sec
- **동시 접속자**: 10,000+ users
- **메시지 크기**: 최대 64KB

### 가용성 기준

- **서비스 가동률**: 99.9%
- **메시지 손실률**: < 0.01%
- **순서 보장률**: 99.99%

## 🔍 모니터링 체크리스트

### Kafka 모니터링 (http://localhost:8080)

- [ ] 메시지 누적 없이 실시간 처리
- [ ] 파티션 간 균등 분배
- [ ] Consumer 지연(lag) < 100ms
- [ ] 토픽 용량 증가 추이

### 서버 로그 확인

- [ ] 모든 연결/해제 이벤트 기록
- [ ] 에러 없이 메시지 처리
- [ ] Kafka 발행 성공률 100%
- [ ] 메모리 사용량 안정성

### 비즈니스 지표

- [ ] 메시지 전달 성공률
- [ ] 사용자 만족도 (응답성)
- [ ] 기능별 사용률 분석

## 🎯 실전 배포 시나리오

### A/B 테스트 시뮬레이션

```json
// 기존 버전 사용자
{
  "id": "user_v1",
  "content": "[V1] 기존 채팅 기능 테스트",
  "timestamp": null
}

// 신규 버전 사용자
{
  "id": "user_v2",
  "content": "[V2] 신규 채팅 기능 테스트 - 더 빠른 전송!",
  "timestamp": null
}
```

### 지역별 사용자 시뮬레이션

```json
// 서울 사용자
{
  "id": "seoul_user",
  "content": "🇰🇷 서울에서 보내는 메시지 (KST 기준)",
  "timestamp": null
}

// 뉴욕 사용자
{
  "id": "ny_user",
  "content": "🇺🇸 Message from New York (EST timezone)",
  "timestamp": null
}
```

## 📝 테스트 실행 순서

1. **환경 설정** (5분)
2. **기본 기능 테스트** (15분)
3. **시나리오별 테스트** (30분)
4. **부하 테스트** (15분)
5. **예외 상황 테스트** (10분)
6. **성능 지표 수집** (5분)

**총 소요 시간: 약 80분**

이 테스트 시나리오를 통해 실제 채팅 서비스의 모든 측면을 검증할 수 있습니다!
