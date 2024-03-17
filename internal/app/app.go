package app

import (
	grpcapp "github.com/Logity-App/sso/internal/app/grpc"
	"github.com/Logity-App/sso/internal/pkg/config"
	"github.com/Logity-App/sso/internal/services/auth"
	"github.com/Logity-App/sso/internal/storage/postgres"

	"os"

	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	dbClient, err := postgres.New(&cfg.Database)

	if err != nil {
		log.Error("error init client db: %s", err)
		os.Exit(0) // TODO
	}

	authService := auth.New(log, dbClient, cfg.App.TokenTTL)

	grpcApp := grpcapp.New(log, cfg.GRPC.Port, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
