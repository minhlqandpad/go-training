package db

import (
	"database/sql"
	"fmt"
)

// TenantDB defines the database operations for tenants
type TenantDB interface {
	GetTenantByName(name string) (*Tenant, error)
	// Add more methods later (e.g., GetTenantByID)
}

type Tenant struct {
	ID    string
	Name  string
	Email string
}

func (db *MySQLDB) GetTenantByName(name string) (*Tenant, error) {
	var tenant Tenant
	query := "SELECT id, name, email FROM tenants WHERE name = ?"
	err := db.QueryRow(query, name).Scan(&tenant.ID, &tenant.Name, &tenant.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query tenant by name: %v", err)
	}
	return &tenant, nil
}
