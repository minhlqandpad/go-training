package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/cmd/server"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/config"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/internal"
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/internal/pkg/worker"
)

func main() {
	curWd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get working directory: %v", err)
	}

	envPath := filepath.Join(curWd, ".env")
	cfg := loadConfig(envPath)

	// init worker pool
	pool, err := worker.NewPool(worker.Config{
		PoolSize:    cfg.PoolSize,
		PoolMin:     cfg.PoolMin,
		MaxJobs:     cfg.MaxJobs,
		NonBlocking: cfg.WorkerPoolNonBlocking,
	})

	if err != nil {
		log.Print("errors: ", err)
		return
	}

	// init async greeter
	service := internal.NewService(pool, cfg.BannedNames)

	// new server
	server := server.NewServer(cfg, pool, service)

	// start server
	server.Start()

	// init signal channel
	sig := make(chan os.Signal, 1)
	// notify the signal channel when receiving interrupt signal
	// syscall.SIGTERM is a signal sent to a process to request its termination
	// os.Interrupt is a signal sent to a process by its controlling terminal when a user wishes to interrupt the process
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// wait for interrupt signal to shutdown the server and release the worker pool
	<-sig

	// shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Stop(ctx)
}

func loadConfig(envPath string) config.Config {
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("[WARN] Error loading env file %v. Loading dafault env", envPath)
	}

	c := config.Config{
		PoolSize:              2,
		PoolMin:               1,
		MaxJobs:               2,
		WorkerPoolNonBlocking: false,
		BannedNames:           map[string]struct{}{},
		HTTPPort:              8080,
	}
	if poolSize := os.Getenv("POOL_SIZE"); poolSize != "" {
		pz, err := strconv.Atoi(poolSize)
		if err != nil {
			log.Fatalf("Invalid pool size: %v", err)
		}
		c.PoolSize = pz
	}
	if poolMin := os.Getenv("POOL_MIN"); poolMin != "" {
		pm, err := strconv.Atoi(poolMin)
		if err != nil {
			log.Fatalf("Invalid pool min: %v", err)
		}
		c.PoolMin = pm
	}
	if maxJobs := os.Getenv("MAX_JOBS"); maxJobs != "" {
		mj, err := strconv.Atoi(maxJobs)
		if err != nil {
			log.Fatalf("Invalid max jobs: %v", err)
		}
		c.MaxJobs = mj
	}
	if nonBlocking := os.Getenv("WORKER_POOL_NON_BLOCKING"); nonBlocking != "" {
		nb, err := strconv.ParseBool(nonBlocking)
		if err != nil {
			log.Fatalf("Invalid worker pool non-blocking: %v", err)
		}
		c.WorkerPoolNonBlocking = nb
	}
	if httpPort := os.Getenv("HTTP_PORT"); httpPort != "" {
		hp, err := strconv.Atoi(httpPort)
		if err != nil {
			log.Fatalf("Invalid http port: %v", err)
		}
		c.HTTPPort = hp
	}
	if bannedNames := os.Getenv("BANNED_NAMES"); bannedNames != "" {
		banned := map[string]struct{}{}
		for _, name := range strings.Split(bannedNames, ",") {
			banned[name] = struct{}{}
		}
		c.BannedNames = banned
	}

	return c
}
