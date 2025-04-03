package db

import (
	"database/sql"
	"fmt"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
)

// TenantDB defines the database operations for tenants
type TenantDB interface {
	GetTenantById(id string) (*domain.Tenant, error)
}

func (db *MySQLDB) GetTenantById(id string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	query := "SELECT id, name, email FROM tenants WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&tenant.ID, &tenant.Name, &tenant.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query tenant by id: %v", err)
	}
	return &tenant, nil
}
