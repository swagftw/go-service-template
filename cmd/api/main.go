package main

import (
	"context"
	"os"
	"os/signal"

	"go-service-template/infrastructure/db/postgres"
	"go-service-template/infrastructure/router"
	"go-service-template/internal/applog"
	"go-service-template/internal/config"
	"go-service-template/services/healthservice"
	"go-service-template/services/healthservice/healthrepo"
	"go-service-template/transport/http"
)

func main() {
	// try to avoid init function calls as they add implicit behaviour

	// init all the deps here
	cfg, err := config.InitConfig()
	if err != nil {
		return
	}

	err = applog.InitLogger(cfg.Debug)
	if err != nil {
		return
	}

	// init infra
	db, err := postgres.Connect(cfg.DatabaseURL)
	if err != nil {
		return
	}

	// init services
	healthRepo := healthrepo.NewRepo(db)
	healthService := healthservice.NewService(healthRepo)

	// start the server
	rtr := router.NewRouter()
	v1Group := rtr.Group("/api/v1")

	// init routes
	http.InitHealthRouter(v1Group, healthService)

	// start the server
	errChan := make(chan error)
	signChan := make(chan os.Signal)

	signal.Notify(signChan, os.Interrupt, os.Kill)
	go func() {
		err = router.Start(cfg.Port, rtr)
		if err != nil {
			errChan <- err
		}
	}()

	applog.Logger.Info(context.Background(), "starting server", map[string]interface{}{
		"port": cfg.Port,
	})

	select {
	case err = <-errChan:
		applog.Logger.Error(context.Background(), err, "failed to start server", map[string]interface{}{
			"port": cfg.Port,
		})
	case sig := <-signChan:
		applog.Logger.Info(context.Background(), "shutting down", map[string]interface{}{
			"signal": sig.String(),
		})
	}
}
