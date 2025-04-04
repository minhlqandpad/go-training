// service/service.go
package service

import (
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain/infrastructure/db"
	v1 "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/tenant/grpc/v1"
)

func NewTenantService(db db.TenantDB) *v1.TenantService {
	return &v1.TenantService{
		DB: db,
	}
}
