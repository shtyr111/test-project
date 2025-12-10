package user_client

import (
	"context"
	"test-project/internal/proto"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcUserClient struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
}

func NewGrpcUserClient(addr string) *GrpcUserClient {
	conn, err := grpc.NewClient(addr,
		grpc.WithInsecure(),             // TLS off
		grpc.WithTimeout(5*time.Second), // Dial timeout
	)

	if err != nil {
		log.Fatal("Произошла ошибка при создании grpc-user-client: %w", err)
	}

	return &GrpcUserClient{
		conn:   conn,
		client: proto.NewUserServiceClient(conn),
	}
}

func (c *GrpcUserClient) Close() {
	c.conn.Close()
}

func (c *GrpcUserClient) GetUser(ctx context.Context, getUserRequest *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.GetUser(ctxWithTimeout, getUserRequest)
}
