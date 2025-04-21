package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain/repository"
	usecase "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/usecase/tenant"
)

type TenantService struct {
	tenant repository.TenantRepository
}

var _ usecase.IGetTenant = (*TenantService)(nil)

func NewTenantService(tenantRepo repository.TenantRepository) *TenantService {
	return &TenantService{
		tenant: tenantRepo,
	}
}

func (s *TenantService) GetTenant(ctx context.Context, req *usecase.GetTenantRequest) (*usecase.GetTenantResponse, error) {
	tenant, err := s.tenant.Get(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, nil
	}
	return &usecase.GetTenantResponse{Tenant: *tenant}, nil
}
