package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/config"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/db"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/tenant/service"
	"google.golang.org/grpc"
)

type Server struct {
	config     config.Config
	grpcServer *grpc.Server
}

func NewServer(cfg config.Config, db *db.MySQLDB) *Server {
	grpcServer := grpc.NewServer()

	// Create tenant service
	tenantService := service.NewTenantService(db)

	// Register the service implementation
	tenantv1.RegisterTenantServiceServer(grpcServer, tenantService)

	return &Server{
		config:     cfg,
		grpcServer: grpcServer,
	}
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %d", s.config.GRPCPort)
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.GracefulStop()
	log.Printf("gRPC server: shutdown")
}
