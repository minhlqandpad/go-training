// internal/tenant/grpc/v1/tenant_service.go
package v1

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain/infrastructure/db"
	pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TenantService struct {
	pb.UnimplementedTenantServiceServer
	DB db.TenantDB
}

// GetTenantByID implements the gRPC method
func (s *TenantService) GetTenantByID(ctx context.Context, req *pb.GetTenantByIDRequest) (*pb.GetTenantResponse, error) {
	tenant, err := s.DB.GetTenantByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tenant: %v", err)
	}
	if tenant == nil {
		return &pb.GetTenantResponse{Tenant: nil}, nil
	}
	return &pb.GetTenantResponse{
		Tenant: &pb.Tenant{
			Id:    tenant.ID,
			Name:  tenant.Name,
			Email: tenant.Email,
		},
	}, nil
}
