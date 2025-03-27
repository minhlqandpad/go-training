package internal

import (
	"context"
	"testing"
	"time"

	"github.com/tuannguyenandpadcojp/go-training/lqm/utils"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/config"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/internal/pkg/worker"
)

func Test_async_greeting_service(t *testing.T) {
	tests := []struct {
		name        string
		inputNames  []string
		wantErr     bool
		banned      map[string]struct{}
		nonBlocking bool
	}{
		{
			name:       "Single valid input name",
			inputNames: []string{"Minh"},
			wantErr:    false,
			banned:     map[string]struct{}{},
		},
		{
			name:       "Multiple valid input names",
			inputNames: []string{"Alice", "Bob"},
			wantErr:    false,
			banned:     map[string]struct{}{},
		},
		{
			name:       "Single banned input name",
			inputNames: []string{"Alice"},
			wantErr:    false,
			banned: map[string]struct{}{
				"Alice": {},
			},
		},
		{
			name:       "Multiple banned input names",
			inputNames: []string{"Alice", "Bob"},
			wantErr:    false,
			banned: map[string]struct{}{
				"Alice": {},
				"Bob":   {},
			},
		},
		{
			name:       "Multiple names but one banned name",
			inputNames: []string{"Alice", "Bob"},
			wantErr:    false,
			banned: map[string]struct{}{
				"Alice": {},
			},
		},
		{
			name:        "Number of names exceed max jobs and pool nonblocking",
			inputNames:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "DDoS"},
			wantErr:     true,
			banned:      map[string]struct{}{},
			nonBlocking: true,
		},
		{
			name:        "Number of names exceed max jobs and pool blocking",
			inputNames:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			wantErr:     false,
			banned:      map[string]struct{}{},
			nonBlocking: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				PoolSize:              4,
				PoolMin:               2,
				MaxJobs:               6,
				WorkerPoolNonBlocking: tt.nonBlocking,
				BannedNames:           tt.banned,
			}

			pool, err := worker.NewPool(worker.Config{
				MaxJobs:     cfg.MaxJobs,
				PoolSize:    cfg.PoolSize,
				PoolMin:     cfg.PoolMin,
				NonBlocking: cfg.WorkerPoolNonBlocking,
			})

			pool.Start(context.Background())

			if err != nil {
				t.Errorf("Error creating pool")
			}

			service := NewService(pool, cfg.BannedNames)
			time.Sleep(10 * time.Second)
			err = service.Greeting(context.Background(), tt.inputNames)
			utils.AssertError(t, err, tt.wantErr)

			pool.Release()
		})
	}
}
