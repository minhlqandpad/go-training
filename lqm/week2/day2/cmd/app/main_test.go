package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/config"
)

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name             string
		prepareVariables func()
		envPath          string
		want             config.Config
	}{
		{
			name:             "default config",
			prepareVariables: func() {},
			envPath:          "",
			want: config.Config{
				PoolSize:              2,
				PoolMin:               1,
				MaxJobs:               2,
				BannedNames:           map[string]struct{}{},
				WorkerPoolNonBlocking: false,
				HTTPPort:              8080,
			},
		},
		{
			name:             "custom env config",
			prepareVariables: func() {},
			envPath:          "test.env",
			want: config.Config{
				PoolSize: 5,
				MaxJobs:  10,
				PoolMin:  3,
				BannedNames: map[string]struct{}{
					"Alice": {},
					"Bob":   {},
				},
				WorkerPoolNonBlocking: true,
				HTTPPort:              8080,
			},
		},
		{
			name: "os env config",
			prepareVariables: func() {
				_ = os.Setenv("POOL_SIZE", "10")
				_ = os.Setenv("POOL_MIN", "5")
				_ = os.Setenv("MAX_JOBS", "20")
				_ = os.Setenv("WORKER_POOL_NON_BLOCKING", "true")
				_ = os.Setenv("BANNED_NAMES", "")
				_ = os.Setenv("HTTP_PORT", "8081")
			},
			envPath: "",
			want: config.Config{
				PoolSize:              10,
				PoolMin:               5,
				MaxJobs:               20,
				BannedNames:           map[string]struct{}{},
				WorkerPoolNonBlocking: true,
				HTTPPort:              8081,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareVariables != nil {
				tt.prepareVariables()
			}
			got := loadConfig(tt.envPath)
			want := tt.want
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("loadConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
