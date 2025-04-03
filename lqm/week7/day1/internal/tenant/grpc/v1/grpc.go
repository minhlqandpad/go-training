package v1

import (
	"context"

	pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/db"
	tenant_v1 "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
)

type TenantService struct {
	tenant_v1.UnimplementedTenantServiceServer
	DB db.TenantDB
}

func (s *TenantService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.GetTenantResponse, error) {
	tenant, err := s.DB.GetTenantById(req.GetId())
	if err != nil {
		return nil, err
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
