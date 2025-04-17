package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/cmd/server"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/config"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week7/day1/internal/infrastructure/db/mysql"
)

func main() {
	// Load .env file
	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("failed to read env variables")
	}

	// Construct MySQL DSN from environment variables
	dsn := config.DSN()

	// Run migrations
	m, err := migrate.New("file://database/migrations", "mysql://"+dsn)
	if err != nil {
		log.Fatalf("failed to initialize migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	log.Println("Building gRPC")
	log.Println("Migrations applied successfully")

	// Connect to MySQL
	mysqlDB, err := mysql.NewMySQLDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer mysqlDB.Close()

	// Start gRPC server
	s := server.NewServer(config, mysqlDB)

	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// init signal channel
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	// shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Stop(ctx)
}
