package usecase

import "context"

type GetTenantRequest struct{}
type GetTenantResponse struct{}

type ITenantGetter interface {
	Get(ctx context.Context, req GetTenantRequest) (GetTenantResponse, error)
}