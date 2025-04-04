package usecase

import "context"

type CreateTenantRequest struct{}
type CreateTenantResponse struct{}

type TenantCreator interface {
	Create(ctx context.Context, req CreateTenantRequest) (CreateTenantResponse, error)
}
