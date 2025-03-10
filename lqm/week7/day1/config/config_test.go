package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/tuannguyenandpadcojp/go-training/lqm/utils"
)

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name             string
		prepareVariables func()
		envPath          string
		want             Config
	}{
		{
			name: "load test env",
			prepareVariables: func(){},
			envPath: "test.env",
			want: Config{
				User: "test",
				Password: "test",
				Host: "test",
				DBPort: 1111,
				HTTPPort: 2222,
				Name: "test",
				GRPCPort: 3333,
			},
		},
		// {
		// 	name: "cannot load env",
		// 	prepareVariables: func() {},
		// 	envPath: ".env",
		// },
		{
			name: "load os env",
			prepareVariables: func() {
				_ = os.Setenv("DB_USER", "root")
				_ = os.Setenv("DB_PASSWORD", "root")
				_ = os.Setenv("DB_HOST", "root")
				_ = os.Setenv("DB_PORT", strconv.Itoa(1000))
				_ = os.Setenv("DB_NAME", "root")
				_ = os.Setenv("HTTP_PORT", strconv.Itoa(2000))
				_ = os.Setenv("GRPC_PORT", strconv.Itoa(3000))
			},
			envPath: "",
			want: Config{
				User: "root",
				Password: "root",
				Name: "root",
				Host: "root",
				DBPort: 1000,
				HTTPPort: 2000,
				GRPCPort: 3000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareVariables != nil {
				tt.prepareVariables()
			}
			got, _ := LoadConfig(tt.envPath)
			want := tt.want
			utils.AssertCorrectResult(t, got, want)
		}) 
	}
}
