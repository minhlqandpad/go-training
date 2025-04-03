package v1

import (
	admin_v1 "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/db"
)

type AdminService struct {
	admin_v1.UnimplementedAdminServiceServer
	db db.TenantDB
}