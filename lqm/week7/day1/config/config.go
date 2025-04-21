package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	User     string
	Password string
	Host     string
	DBPort   int
	HTTPPort int
	Name     string
	GRPCPort int
}

func LoadConfig(envPath string) (Config, error) {
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			return Config{}, fmt.Errorf("failed to load .env file from %s: %v", envPath, err)
		}
	} else {
		godotenv.Load()
	}

	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_PORT value '%s': %v", portStr, err)
	}
	httpPortStr := os.Getenv("HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_PORT value '%s': %v", httpPortStr, err)
	}
	grpcPortStr := os.Getenv("GRPC_PORT")
	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid GRPC_PORT value '%s': %v", grpcPortStr, err)
	}

	config := Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		DBPort:   port,
		HTTPPort: httpPort,
		Name:     os.Getenv("DB_NAME"),
		GRPCPort: grpcPort,
	}

	if err := config.validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.DBPort, c.Name)
}

func (c Config) validate() error {
	missing := []string{}
	if c.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DBPort == 0 {
		missing = append(missing, "DB_PORT")
	}
	if c.HTTPPort == 0 {
		missing = append(missing, "HTTP_PORT")
	}
	if c.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if c.GRPCPort == 0 {
		missing = append(missing, "GRPC_PORT")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}
	return nil
}
