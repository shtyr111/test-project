package user

import (
	"context"
	"test-project/internal/grpc_client/user_client"
	"test-project/internal/proto"
	"test-project/internal/service/user"

	log "github.com/sirupsen/logrus"
)

type GrpcUserHandler struct {
	userService    *user.UserService
	grpcUserClient *user_client.GrpcUserClient
	proto.UnimplementedUserServiceServer
}

func NewGrpcUserHandler(userService *user.UserService, grpcUserClient *user_client.GrpcUserClient) *GrpcUserHandler {
	return &GrpcUserHandler{userService: userService, grpcUserClient: grpcUserClient}
}

func (g GrpcUserHandler) GetUser(ctx context.Context, getUserRequest *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	grpcResponse, err := g.grpcUserClient.GetUser(ctx, getUserRequest)

	if err != nil {
		log.Error(err)
	}

	return grpcResponse, err
}

func (h *GrpcUserHandler) mustEmbedUnimplementedUserServiceServer() {
	// Пустой метод — просто удовлетворяет интерфейс
}
