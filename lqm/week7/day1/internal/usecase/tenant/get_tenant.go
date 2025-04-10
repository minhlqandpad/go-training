package usecase

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
)

type GetTenantRequest struct {
	TenantID string
}

type GetTenantResponse struct {
	Tenant domain.Tenant
}

type IGetTenant interface {
	GetTenant(ctx context.Context, req *GetTenantRequest) (*GetTenantResponse, error)
}
