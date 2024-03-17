package main

import (
	"github.com/Logity-App/sso/internal/app"
	config "github.com/Logity-App/sso/internal/pkg/config"
	"github.com/Logity-App/sso/internal/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		panic(err)
	}

	log := logger.SetupLogger(cfg.App.Env)

	log.Info("str", slog.String("env", cfg.App.Env))

	application := app.New(log, cfg)

	go func() {
		application.GRPCSrv.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()

	log.Info("Application stop")
}
