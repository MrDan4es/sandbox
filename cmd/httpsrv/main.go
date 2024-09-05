package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrdan4es/sandbox/internal/http/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("unexpected", "error", err)
		os.Exit(1)
	}

	slog.Info("successfully stopped server")
}

func run(ctx context.Context) error {
	httpSrv := server.NewFileUploadServer()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutting down http server", "error", err)
		}
	}()

	slog.Info("starting http server on", "addr", httpSrv.Addr, "pid", os.Getpid())

	return httpSrv.ListenAndServe()
}
