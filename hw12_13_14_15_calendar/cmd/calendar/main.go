package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"gopkg.in/yaml.v2"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Конфигурация приложения:
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	// Настройка логгера
	logFile, err := filepath.Abs(config.Logger.File)
	if err != nil {
		log.Fatalf("Invalid log file name: %s: %s", config.Logger.File, err)
	}

	logg := logger.New(config.Logger.Level, logFile)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func loadConfig(configPath string) (*Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config %s: %s", configPath, err)
	}

	cfg := &Config{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to read yaml: %s", err)
	}

	return cfg, nil
}
