package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
	pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
	usecase "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/usecase/tenant"
)

type TenantServiceServer struct {
	pb.UnimplementedTenantServiceServer
	GetTenantService usecase.IGetTenant
}

func NewTenantService(uc usecase.IGetTenant) *TenantServiceServer {
	return &TenantServiceServer{
		GetTenantService: uc,
	}
}

func toPbTenant(tenant domain.Tenant) pb.Tenant {
	return pb.Tenant{
		Id:    tenant.ID,
		Name:  tenant.Name,
		Email: tenant.Email,
	}
}

func (s *TenantServiceServer) GetTenantByID(ctx context.Context, req *pb.GetTenantByIDRequest) (*pb.GetTenantResponse, error) {
	usecaseReq := &usecase.GetTenantRequest{
		TenantID: req.Id,
	}

	resp, err := s.GetTenantService.GetTenant(ctx, usecaseReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tenant: %v", err)
	}

	respTenant := toPbTenant(resp.Tenant)
	return &pb.GetTenantResponse{
		Tenant: &respTenant,
	}, nil
}
