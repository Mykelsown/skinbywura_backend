// Package main boots the API application.
// It loads configuration, migrates the database, and serves HTTP until told to stop.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mykelsown/skinbywura_backend.git/config"
	"github.com/Mykelsown/skinbywura_backend.git/db"
	"github.com/Mykelsown/skinbywura_backend.git/internal/mailer"
	"github.com/Mykelsown/skinbywura_backend.git/internal/migrate"
	"github.com/Mykelsown/skinbywura_backend.git/internal/server"
	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	// A .env file is a local-development convenience. On Render/Docker the real
	// environment is already set, so a missing file is normal and not an error.
	_ = godotenv.Load()

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connection retries briefly because managed databases can lag the app at first boot.
	dbStore, err := store.Connection(cfg.DB)
	if err != nil {
		return err
	}
	defer dbStore.Pool.Close()
	slog.Info("database connection verified")

	if cfg.AutoMigrate {
		if err := migrate.Up(ctx, dbStore.Pool, db.Migrations, "migrations"); err != nil {
			return err
		}
	}

	svc := service.New(dbStore, service.Options{
		Mailer:      mailer.New(cfg.ResendAPIKey, cfg.MailFrom),
		FrontendURL: cfg.FrontendURL,
		StoreEmail:  cfg.StoreEmail,
		StoreName:   cfg.StoreName,
		PaymentNote: cfg.PaymentNote,
	})
	go svc.CleanupLoop(ctx)

	return server.New(cfg, svc).Run(ctx)
}
