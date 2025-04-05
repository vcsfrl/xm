package cmd

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/vcsfrl/xm/internal/api"
	"github.com/vcsfrl/xm/internal/api/handler"
	"github.com/vcsfrl/xm/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Run application.
func run() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	restApi := api.NewRestApi(ctx, logger, appConfig, db)
	// run api
	go func() {
		restApi.Run()
		stop()
	}()

	// Shut down app.
	shutdown := func() {
		if err := restApi.Close(); err != nil {
			logger.Error().Err(err).Msg("Close api.")
		}

		logger.Info().Msg("Close db.")
		dbInstance, err := db.DB()
		if err != nil {
			logger.Error().Err(err).Msg("Get db.")
		}
		if err := dbInstance.Close(); err != nil {
			logger.Error().Err(err).Msg("Close db.")
		}

		logger.Info().Msg("Close log.")
		logger.Info().Msg("Exit application.")
		_ = loggerOutput.Close()
		os.Exit(0)
	}

	defer func() {
		// handle panic
		if err := recover(); err != nil { //catch
			logger.Error().Msgf("Panic in application: %v", err)
			shutdown()
		}
	}()

	go func() {
		sig := <-ctx.Done()
		logger.Info().Msgf("Caught a stop signal: %+v.", sig)
		shutdown()
	}()

	runDebug(appConfig, logger)

	// Wait...
	select {}
}

// Start Debug Endpoints.
//
// /debug/pprof
// /debug/vars
//
// Not concerned with shutting this down when the application is shutdown.
func runDebug(cfg *config.Config, logger zerolog.Logger) {
	go func() {

		traceHost := fmt.Sprintf("0.0.0.0:%s", cfg.TracePort)
		logger.Info().Str("port", traceHost).Msg("Debug endpoint started.")

		if err := http.ListenAndServe(traceHost, handler.NewDebug(logger)); err != nil {
			logger.Error().Err(err).Msg("Debug endpoint start.")
		}
	}()
}
