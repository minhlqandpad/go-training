package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/db/mock_db"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
	pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
)

func TestTenantService_GetTenant(t *testing.T) {
	tests := []struct {
		name       string
		inputId    string
		mockTenant *domain.Tenant
		mockError  error
		wantTenant *pb.Tenant
		wantErr    bool
	}{
		{
			name:       "tenant found",
			inputId:    "tenant-123",
			mockTenant: &domain.Tenant{ID: "tenant-123", Name: "test", Email: "test@example.com"},
			mockError:  nil,
			wantTenant: &pb.Tenant{Id: "tenant-123", Name: "test", Email: "test@example.com"},
			wantErr:    false,
		},
		{
			name:       "tenant not found",
			inputId:    "unknown",
			mockTenant: &domain.Tenant{ID: "tenant-123", Name: "test", Email: "test@example.com"},
			mockError:  nil,
			wantTenant: nil,
			wantErr:    false,
		},
		{
			name:       "database error",
			inputId:    "error-case",
			mockTenant: nil,
			mockError:  errors.New("database failure"),
			wantTenant: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mock
			mockDB := mocks.NewMockTenantDB(ctrl)
			mockDB.EXPECT().GetTenantById(tt.inputId).Return(tt.mockTenant, tt.mockError)

			// Create service with mock
			s := NewTenantService(mockDB)

			// Call GetTenant
			resp, err := s.GetTenant(context.Background(), &pb.GetTenantRequest{Id: tt.inputId})

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError, err) // Check exact error
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantTenant, resp.GetTenant())
		})
	}
}
