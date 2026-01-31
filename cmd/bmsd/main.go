// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Server entry point.
// This file defines the bmsd main function, which resolves configuration,
// initializes structured logging with server defaults, emits startup
// diagnostics (redacted), and exits on configuration errors.

package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SandorMiskey/bms-core/internal/config"
	"github.com/SandorMiskey/bms-core/internal/errtext"
	"github.com/SandorMiskey/bms-core/internal/health"
	"github.com/SandorMiskey/bms-core/internal/logging"
)

// Main entry point. {{{

func main() {
	configPath := flag.String("config", "", "path to config.toml")
	flag.Parse()

	configResult, path, warnings, err := config.ResolveConfigDiagnostics(*configPath, config.ConfigOverlay{}, config.ConfigOverlay{})
	logger, format := initLogger(configResult, logging.ComponentServer)

	if err != nil {
		var validationErrors config.ValidationErrors
		if errors.As(err, &validationErrors) {
			logging.LogConfigDiagnostics(logger, format, configResult, path, warnings)
			logger.Error(errtext.ErrConfigValidationFailed, "error", err)
		} else {
			logger.Error(errtext.ErrConfigResolutionFailed, "error", err)
		}
		os.Exit(1)
	}

	logging.LogConfigDiagnostics(logger, format, configResult, path, warnings)

	healthState := health.NewState()
	healthServer := startHealthServer(logger, configResult.REST.Address, healthState)
	healthState.SetReady(true)

	waitForShutdown(logger, healthServer)
}

func initLogger(cfg config.Config, component logging.Component) (*slog.Logger, config.LogFormat) {
	defaults := logging.LoggerDefaults{
		Fields: logging.DefaultFields{
			Component:   component,
			Environment: string(cfg.Server.Environment),
			ServerID:    cfg.Server.ID,
		},
		Format: config.LogFormatJSON,
		Level:  config.LogLevelInfo,
	}

	logger, err := logging.NewLogger(cfg.Logging, defaults)
	if err != nil {
		fallback := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
		fallback.Error(errtext.ErrLoggerInitFailed, "error", err)
		logger = fallback
	}

	format := cfg.Logging.Format
	if format == "" || err != nil {
		format = defaults.Format
	}

	return logger, format
}

func startHealthServer(logger *slog.Logger, address string, state *health.State) *http.Server {
	if address == "" {
		logger.Warn("health server disabled", "reason", "rest address is empty")
		return nil
	}

	server := &http.Server{
		Addr:              address,
		Handler:           health.NewMux(state),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(errtext.ErrHealthServerServeFailed, "error", err, "address", address)
		}
	}()

	return server
}

func waitForShutdown(logger *slog.Logger, server *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	if server == nil {
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(errtext.ErrHealthServerShutdownFailed, "error", err)
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
