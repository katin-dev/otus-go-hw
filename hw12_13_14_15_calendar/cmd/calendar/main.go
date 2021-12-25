package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"gopkg.in/yaml.v2"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Init: App Config
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	logg, err := logger.New(config.Logger.File, config.Logger.Level, config.Logger.Formatter)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Init: Storage
	var storage app.Storage
	switch config.Storage.Type {
	case StorageMem:
		storage = memorystorage.New()
	case StorageSQL:
		storage = sqlstorage.New(ctx, config.Storage.Dsn)
	default:
		log.Fatalf("Unknown storage type: %s\n", config.Storage.Type)
	}
	defer cancel()

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.HTTP.Host, config.HTTP.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
	}
}

func loadConfig(configPath string) (*Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", configPath, err)
	}

	cfg := NewConfig()
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml: %w", err)
	}

	return &cfg, nil
}
