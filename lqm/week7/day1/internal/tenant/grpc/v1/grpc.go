package v1

import (
	"context"

	pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/tenant/service"
	usecase "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/usecase/tenant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TenantServiceServer struct {
	pb.UnimplementedTenantServiceServer
	TenantService service.TenantService
}

func NewTenantService(s *service.TenantService) *TenantServiceServer{
	return &TenantServiceServer{
		TenantService: *s,
	}
}

func toGRPCResponse(resp *usecase.GetTenantResponse) *pb.GetTenantResponse {
	return &pb.GetTenantResponse{
		Tenant: &pb.Tenant{
			Id:    resp.Tenant.ID,
			Name:  resp.Tenant.Name,
			Email: resp.Tenant.Email,
		},
	}
}

func toUsecaseRequest(req *pb.GetTenantByIDRequest) *usecase.GetTenantRequest {
	return &usecase.GetTenantRequest{
		TenantID: req.Id,
	}
}

func (s *TenantServiceServer) GetTenantByID(ctx context.Context, req *pb.GetTenantByIDRequest) (*pb.GetTenantResponse, error) {
	usecaseReq := toUsecaseRequest(req)
	resp, err := s.TenantService.GetTenant(ctx, usecaseReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tenant: %v", err)
	}
	return toGRPCResponse(resp), nil
}
