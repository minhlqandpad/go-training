package v1

import (
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain/infrastructure/db"
	admin_v1 "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
)

type AdminService struct {
	admin_v1.UnimplementedAdminServiceServer
	db db.TenantDB
}
