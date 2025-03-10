package service

import (
	"context"

	pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/gen/go/tenant/v1"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/db"
)

type TenantService struct {
	pb.UnimplementedTenantServiceServer
	db db.TenantDB
}

func NewTenantService(db db.TenantDB) *TenantService {
	return &TenantService{
		db: db,
	}
}

func (s *TenantService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.GetTenantResponse, error) {
	tenant, err := s.db.GetTenantByName(req.GetName())
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
