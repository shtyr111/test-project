package grpc

import (
	"net"
	"test-project/internal/proto"
	"test-project/internal/transport/grpc/handlers/user"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpcUserHandler *user.GrpcUserHandler
}

func NewGrpcServer(grpcUserHandler *user.GrpcUserHandler) *GrpcServer {
	return &GrpcServer{grpcUserHandler: grpcUserHandler}
}

func (gs GrpcServer) RunGrpcServer() {
	log.Info("grpc server started on port 9000")

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// ✅ Регистрируем Prometheus метрики для gRPC
	grpc_prometheus.Register(s)
	proto.RegisterUserServiceServer(s, gs.grpcUserHandler)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed start grpc server: %v", err)
	}

}
