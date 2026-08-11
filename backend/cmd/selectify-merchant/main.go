package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"alwis.dev/selectify/internal/helper"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
	"alwis.dev/selectify/merchant-pkg/app"
	"alwis.dev/selectify/merchant-pkg/routers"
)

func main() {
	ctx := context.Background()
	logger.Init()

	// Set application prefix for environment variable prefixing
	program.AppPrefix = "MRC"

	err := helper.LoadEnv(ctx)
	if err != nil {
		panic(err)
	}

	app.NewAppEnvironment()

	h := routers.CreateHandler()

	s := &http.Server{
		Addr:    ":3004",
		Handler: h,
	}
	go func() {
		logger.Infof(ctx, "Server starting on %s", s.Addr)
		if err = s.ListenAndServe(); err != nil {
			logger.Fatal(ctx, err, "Server failed to start")
		}
	}()

	logger.Info(ctx, "Server started successfully")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "Shutting down server...")
	if err = s.Shutdown(ctx); err != nil {
		_ = logger.Errorf(ctx, err, "Server forced to shutdown")
	}
	//app.Close()

	logger.Info(ctx, "Server exited")
}
