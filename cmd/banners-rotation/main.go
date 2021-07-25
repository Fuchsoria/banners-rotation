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

	"github.com/Fuchsoria/banners-rotation/internal/app"
	"github.com/Fuchsoria/banners-rotation/internal/bandit"
	"github.com/Fuchsoria/banners-rotation/internal/logger"
	gw "github.com/Fuchsoria/banners-rotation/internal/server/grpc"
	sqlstorage "github.com/Fuchsoria/banners-rotation/internal/storage/sql"
	"github.com/Fuchsoria/banners-rotation/internal/version"
	_ "github.com/lib/pq"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/banners-rotation/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()

		return
	}

	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New(config.Logger.Level, config.Logger.File)

	ctx, cancel := context.WithCancel(context.Background())

	storage, err := initStorage(ctx, config)
	if err != nil {
		logg.Error(err.Error())

		log.Fatal(err)
	}

	brApp := app.New(logg, storage)

	server, err := gw.NewServer(brApp, config.HTTP.Host, config.HTTP.Port, config.HTTP.GrpcPort)
	if err != nil {
		logg.Error(err.Error())
	}

	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("banners rotation service is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func initStorage(ctx context.Context, config Config) (*sqlstorage.Storage, error) {
	bandit := bandit.New()

	storage, err := sqlstorage.New(ctx, config.DB.ConnectionString, bandit)
	if err != nil {
		return nil, fmt.Errorf("can't create new storage instance, %w", err)
	}

	err = storage.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't connect to storage, %w", err)
	}

	err = storage.ClearSessionClicks()
	if err != nil {
		return nil, fmt.Errorf("cannot clear session clicks table, %w", err)
	}

	err = storage.ClearSessionViews()
	if err != nil {
		return nil, fmt.Errorf("cannot clear session views table, %w", err)
	}

	return storage, nil
}
