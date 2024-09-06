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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mrdan4es/sandbox/api/fileuploadpb/v1"
	"github.com/mrdan4es/sandbox/internal/http/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error("unexpected", "error", err)
		os.Exit(1)
	}

	slog.Info("successfully stopped http server")
}

func run(ctx context.Context) error {
	conn, err := grpc.NewClient(
		":12345",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	httpSrv := server.NewFileUploadServer(pb.NewFileUploadServiceClient(conn))

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutting down http server", "error", err)
		}
	}()

	slog.Info("starting http server on", "addr", httpSrv.Addr, "pid", os.Getpid())

	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
