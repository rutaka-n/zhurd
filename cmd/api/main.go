package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"zhurd/internal/adapters/httpapi"
	"zhurd/internal/config"
)

func main() {
	cfgFile, err := os.Open("./share/config.json.example")
	if err != nil {
		panic(err)
	}
	cfg, err := config.Load(cfgFile)
	if err != nil {
		panic(err)
	}

	apiRouter, err := httpapi.New()
	if err != nil {
		panic(err)
	}

	var logDest io.Writer
	if cfg.Server.Logger.Destination == "stdout" {
		logDest = os.Stdout
	} else {
		logFile, err := os.OpenFile(cfg.Server.Logger.Destination, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		logDest = logFile
	}

	var logger *slog.Logger
	switch cfg.Server.Logger.Format {
	case "json":
		logger = slog.New(slog.NewJSONHandler(logDest, nil))
	case "text":
		logger = slog.New(slog.NewJSONHandler(logDest, nil))
	default:
		panic("unknown logger format, supported formats are: 'json', 'text'")
	}
	slog.SetDefault(logger)

	srv := &http.Server{
		Addr: cfg.Server.Addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      apiRouter, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("fail to run API server", "error", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Server.GracefulTimeoutSec))
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	slog.Info("shutting down")
	os.Exit(0)
}
