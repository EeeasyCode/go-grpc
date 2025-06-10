package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "example.com/grpc-chat-app/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ChatClient: gRPC 채팅 클라이언트
type ChatClient struct {
	client pb.BroadcastClient
	conn   *grpc.ClientConn
	userID string
}

// NewChatClient: 새로운 채팅 클라이언트 생성
func NewChatClient(serverAddr, userID string) (*ChatClient, error) {
	// gRPC 서버에 연결
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := pb.NewBroadcastClient(conn)

	return &ChatClient{
		client: client,
		conn:   conn,
		userID: userID,
	}, nil
}

// Close: 클라이언트 연결 종료
func (c *ChatClient) Close() error {
	return c.conn.Close()
}

// ConnectAndListen: 서버에 연결하고 메시지 수신 대기
func (c *ChatClient) ConnectAndListen(ctx context.Context) error {
	// 서버에 연결 요청
	connectMsg := &pb.Connect{
		User: &pb.User{
			Id:   c.userID,
			Name: fmt.Sprintf("User-%s", c.userID),
		},
		Active: true,
	}

	// 스트림 생성
	stream, err := c.client.CreateStream(ctx, connectMsg)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	log.Printf("Connected to server as user: %s", c.userID)

	// 메시지 수신 루프
	for {
		select {
		case <-ctx.Done():
			log.Printf("Client %s disconnecting...", c.userID)
			return ctx.Err()
		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Server closed the stream for user: %s", c.userID)
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to receive message: %w", err)
			}

			// 수신한 메시지 출력
			timestamp := "unknown"
			if msg.Timestamp != nil {
				timestamp = msg.Timestamp.AsTime().Format("15:04:05")
			}
			
			fmt.Printf("[%s] %s: %s\n", timestamp, msg.Id, msg.Content)
		}
	}
}

// SendMessage: 메시지 전송
func (c *ChatClient) SendMessage(ctx context.Context, content string) error {
	msg := &pb.Message{
		Id:        c.userID,
		Content:   content,
		Timestamp: timestamppb.Now(),
	}

	_, err := c.client.BroadcastMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Message sent: %s", content)
	return nil
}

// runClient: 클라이언트 실행 함수
func runClient(userID string) {
	serverAddr := "localhost:8081"

	// 클라이언트 생성
	client, err := NewChatClient(serverAddr, userID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 연결 및 메시지 수신을 별도 고루틴에서 실행
	go func() {
		if err := client.ConnectAndListen(ctx); err != nil {
			log.Printf("Listen error: %v", err)
		}
	}()

	// 잠시 대기 후 테스트 메시지 전송
	time.Sleep(2 * time.Second)

	// 몇 개의 테스트 메시지 전송
	messages := []string{
		fmt.Sprintf("Hello from %s!", userID),
		fmt.Sprintf("This is %s testing the chat", userID),
		fmt.Sprintf("%s says goodbye!", userID),
	}

	for i, msg := range messages {
		if err := client.SendMessage(ctx, msg); err != nil {
			log.Printf("Failed to send message %d: %v", i+1, err)
		}
		time.Sleep(3 * time.Second) // 메시지 간 간격
	}

	// 추가로 10초 대기하여 다른 클라이언트의 메시지를 수신할 수 있도록 함
	time.Sleep(10 * time.Second)
	
	log.Printf("Client %s finished", userID)
}

// main 함수 - 클라이언트를 독립적으로 실행하기 위해 주석 해제됨
func init() {
	// main.go와 함께 빌드될 때는 이 함수가 실행되지 않도록 함
}

// runClientMain: 독립 실행 시 사용
func runClientMain() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <user_id>")
		fmt.Println("Example: go run client.go alice")
		return
	}
	
	userID := os.Args[1]
	runClient(userID)
} 