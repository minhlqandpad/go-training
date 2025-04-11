package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
	mocks "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain/repository/mock_repository"
	usecase "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/usecase/tenant"
)

func TestTenantService_GetTenantByID(t *testing.T) {
	tests := []struct {
		name       string
		inputId    string
		mockTenant *domain.Tenant
		mockError  error
		wantTenant *usecase.GetTenantResponse
		wantErr    error // Use error type to check gRPC status codes
	}{
		{
			name:       "tenant found",
			inputId:    "tenant-123",
			mockTenant: &domain.Tenant{ID: "tenant-123", Name: "test", Email: "test@example.com"},
			mockError:  nil,
			wantTenant: &usecase.GetTenantResponse{
				Tenant: domain.Tenant{
					ID:    "tenant-123",
					Name:  "test",
					Email: "test@example.com"},
			},
			wantErr: nil,
		},
		{
			name:       "tenant not found",
			inputId:    "unknown",
			mockTenant: nil, // Mock returns nil for not found
			mockError:  nil,
			wantTenant: nil, // Expect Tenant field to be nil
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mock
			mockDB := mocks.NewMockTenantRepository(ctrl)
			mockDB.EXPECT().Get(gomock.Any(), tt.inputId).Return(tt.mockTenant, tt.mockError)

			// Create service with mock
			s := NewTenantService(mockDB)

			// Call GetTenantByID
			resp, err := s.GetTenant(context.Background(), &usecase.GetTenantRequest{TenantID: tt.inputId})

			// Assert results
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, status.Code(tt.wantErr), status.Code(err), "error codes should match")
				assert.Contains(t, err.Error(), tt.mockError.Error(), "error message should contain mock error")
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantTenant, resp, "tenant should match expected value")
		})
	}
}
