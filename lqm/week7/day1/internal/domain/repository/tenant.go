package repository

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
)

type TenantRepository interface {
	Get(ctx context.Context, id string) (*domain.Tenant, error)
}
