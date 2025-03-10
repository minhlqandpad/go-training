package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/gen/go/tenant/v1"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/db"
	mocks "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/db/mock_db"
)

func TestTenantService_GetTenant(t *testing.T) {
	tests := []struct {
		name       string
		inputName  string
		mockTenant *db.Tenant
		mockError  error
		wantTenant *pb.Tenant
		wantErr    bool
	}{
		{
			name:       "tenant found",
			inputName:  "test",
			mockTenant: &db.Tenant{ID: "tenant-123", Name: "test", Email: "test@example.com"},
			mockError:  nil,
			wantTenant: &pb.Tenant{Id: "tenant-123", Name: "test", Email: "test@example.com"},
			wantErr:    false,
		},
		{
			name:       "tenant not found",
			inputName:  "unknown",
			mockTenant: nil,
			mockError:  nil,
			wantTenant: nil,
			wantErr:    false,
		},
		{
			name:       "database error",
			inputName:  "error-case",
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
			mockDB.EXPECT().GetTenantByName(tt.inputName).Return(tt.mockTenant, tt.mockError)

			// Create service with mock
			s := NewTenantService(mockDB)

			// Call GetTenant
			resp, err := s.GetTenant(context.Background(), &pb.GetTenantRequest{Name: tt.inputName})

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
