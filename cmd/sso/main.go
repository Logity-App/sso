package main

import (
	"fmt"
	"github.com/Logity-App/sso/internal/app"
	config "github.com/Logity-App/sso/internal/pkg/config"
	"github.com/Logity-App/sso/internal/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	log := logger.SetupLogger(cfg.App.Env)

	//log.Info("str", slog.String("env", cfg.App.Env))

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
