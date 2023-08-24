package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	memoryStorage "github.com/arsenalvlad/hw12_13_14_15_calendar/internal/storage/memory"
	sqlStorage "github.com/arsenalvlad/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/arsenalvlad/hw12_13_14_15_calendar/internal/server/http"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)
	defer logg.Sync()

	if config.Storage.Type == "postgres" {
		err := migrateAction(logg, config.Storage.Postgres)

		if err != nil {
			logg.Fatal("failed to migrate up database "+err.Error(), zap.Any("psql_setting", config.Storage.Postgres))
		}
	}

	var newStorage app.Storage

	switch config.Storage.Type {
	case "memory":
		newStorage = memoryStorage.New()
	case "postgres":
		newStorage = sqlStorage.New(config.Storage.Postgres.DSN())
	}

	calendar := app.New(logg, newStorage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Address())

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

func migrateAction(logg *logger.Logger, conf Postgres) error { //nolint: cyclop
	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", conf.Migration.Path),
		conf.MigrateDSN(),
	)
	if err != nil {
		return fmt.Errorf("could not connect to migrate db: %w", err)
	}

	logg.Info("migrate new create...")

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not migrate up: %w", err)
	}

	logg.Info("migrate up done...")

	return nil
}
