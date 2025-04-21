package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain/repository"
)

type tenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) repository.TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

func (r *tenantRepository) Get(ctx context.Context, id string) (*domain.Tenant, error) {
	query := "SELECT id, name, email FROM tenants WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var tenant domain.Tenant
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query tenant: %v", err)
	}

	return &tenant, nil
}
