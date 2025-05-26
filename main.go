package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "example.com/grpc-chat-app/gen" // 생성된 Go 코드 패키지 임포트 (여러분의 go_package 경로에 맞게 수정)
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	// 추가: 타임스탬프 사용
)

// Connection: 개별 클라이언트 연결을 나타내는 구조체
type Connection struct {
	pb.UnimplementedBroadcastServer        // 항상 임베드해야 함
	stream                          pb.Broadcast_CreateStreamServer // 클라이언트에게 메시지를 보내기 위한 스트림
	id                              string                            // 연결 ID (예: 사용자 ID)
	active                          bool                              // 연결 활성 상태
	error                           chan error                        // 에러 전파를 위한 채널
}

// Pool: 활성 연결들의 풀(모음)을 관리하는 구조체
type Pool struct {
	pb.UnimplementedBroadcastServer // 항상 임베드해야 함
	Connection                      []*Connection
}

// CreateStream: 클라이언트가 연결을 요청하고 메시지 스트림을 설정
func (p *Pool) CreateStream(pconn *pb.Connect, stream pb.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	// 새로운 연결을 풀에 추가
	p.Connection = append(p.Connection, conn)
	log.Printf("User %s connected", conn.id)

	// 스트림이 활성 상태인 동안 에러 채널을 통해 대기 (에러 발생 시 해당 에러 반환)
	return <-conn.error
}

// BroadcastMessage: 메시지를 받아 모든 활성 연결에 전송
func (p *Pool) BroadcastMessage(ctx context.Context, msg *pb.Message) (*pb.Close, error) {
	wait := sync.WaitGroup{} // 모든 고루틴이 완료될 때까지 기다리기 위한 WaitGroup
	done := make(chan int)   // 모든 메시지 전송이 완료되었음을 알리는 채널

	// 현재 시간을 메시지에 설정 (원문에는 없지만 추가하면 좋음)
	if msg.Timestamp == nil {
		msg.Timestamp = timestamppb.Now()
	}

	// 풀에 있는 모든 연결을 순회
	for _, conn := range p.Connection {
		wait.Add(1) // WaitGroup 카운터 증가

		// 각 연결에 대해 비동기적으로 메시지 전송 (고루틴 사용)
		go func(msg *pb.Message, conn *Connection) {
			defer wait.Done() // 고루틴 완료 시 WaitGroup 카운터 감소

			if conn.active {
				// 클라이언트 스트림으로 메시지 전송
				if err := conn.stream.Send(msg); err != nil {
					log.Printf("Error sending message to %s: %v", conn.id, err)
					// 에러 발생 시 해당 연결은 비활성 처리하고 에러 채널에 에러 전송
					conn.active = false
					conn.error <- err
				} else {
					log.Printf("Sent message to %s from %s", conn.id, msg.Id)
				}
			}
		}(msg, conn)
	}

	// 모든 메시지 전송 고루틴이 완료될 때까지 기다리는 별도의 고루틴
	go func() {
		wait.Wait()  // 모든 Add 호출에 대한 Done 호출이 완료될 때까지 대기
		close(done) // 완료 채널을 닫아 대기 중인 곳에 알림
	}()

	<-done // 모든 메시지 전송이 완료될 때까지 여기서 대기
	return &pb.Close{}, nil // 빈 Close 메시지 반환
}

func main() {
	// 새 gRPC 서버 인스턴스 생성
	grpcServer := grpc.NewServer()

	// 연결 풀 초기화
	var connections []*Connection // 빈 Connection 슬라이스
	pool := &Pool{
		Connection: connections,
	}

	// 생성된 BroadcastServer 인터페이스에 우리 Pool 구현을 등록
	pb.RegisterBroadcastServer(grpcServer, pool)

	// TCP 리스너 생성 (예: 8080 포트)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error creating TCP listener: %v", err)
	}

	fmt.Println("gRPC Server started at port :8080")

	// gRPC 서버가 리스너를 통해 요청을 받기 시작
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error serving gRPC requests: %v", err)
	}
}