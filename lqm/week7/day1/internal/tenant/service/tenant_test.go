package service

import (
    "context"
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain"
    mocks "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/domain/infrastructure/db/mock_db"
    pb "github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/pb/v1"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func TestTenantService_GetTenantByID(t *testing.T) {
    tests := []struct {
        name       string
        inputId    string
        mockTenant *domain.Tenant
        mockError  error
        wantTenant *pb.Tenant
        wantErr    error // Use error type to check gRPC status codes
   }{
        {
            name:       "tenant found",
            inputId:    "tenant-123",
            mockTenant: &domain.Tenant{ID: "tenant-123", Name: "test", Email: "test@example.com"},
            mockError:  nil,
            wantTenant: &pb.Tenant{Id: "tenant-123", Name: "test", Email: "test@example.com"},
            wantErr:    nil,
        },
        {
            name:       "tenant not found",
            inputId:    "unknown",
            mockTenant: nil, // Mock returns nil for not found
            mockError:  nil,
            wantTenant: nil, // Expect Tenant field to be nil
            wantErr:    nil,
        },
        {
            name:       "database error",
            inputId:    "error-case",
            mockTenant: nil,
            mockError:  errors.New("database failure"),
            wantTenant: nil,
            wantErr:    status.Errorf(codes.Internal, "failed to get tenant: %v", errors.New("database failure")),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set up gomock controller
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            // Create mock
            mockDB := mocks.NewMockTenantDB(ctrl)
            mockDB.EXPECT().GetTenantByID(gomock.Any(), tt.inputId).Return(tt.mockTenant, tt.mockError)

            // Create service with mock
            s := NewTenantService(mockDB)

            // Call GetTenantByID
            resp, err := s.GetTenantByID(context.Background(), &pb.GetTenantByIDRequest{Id: tt.inputId})

            // Assert results
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, status.Code(tt.wantErr), status.Code(err), "error codes should match")
                assert.Contains(t, err.Error(), tt.mockError.Error(), "error message should contain mock error")
            } else {
                assert.NoError(t, err)
            }
            assert.Equal(t, tt.wantTenant, resp.GetTenant(), "tenant should match expected value")
        })
    }
}