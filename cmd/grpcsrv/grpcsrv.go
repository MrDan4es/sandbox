package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	pb "github.com/mrdan4es/sandbox/api/fileuploadpb/v1"
	"github.com/mrdan4es/sandbox/internal/grpc/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error("unexpected", "error", err)
		os.Exit(1)
	}

	slog.Info("successfully stopped grpc server")
}

func run(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}

	serverRegister := grpc.NewServer()
	go func() {
		<-ctx.Done()

		slog.Info("gracefully stopping gRPC server")
		serverRegister.GracefulStop()
	}()

	pb.RegisterFileUploadServiceServer(serverRegister, server.New())

	metrics := server.NewMetrics()
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := metrics.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutting down http server", "error", err)
		}
	}()
	go func() {
		slog.Info("starting http metric server on", "addr", metrics.Addr)

		if err := metrics.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server startup failed", "error", err)
		}
	}()

	slog.Info("gRPC server started on port :12345")
	if err = serverRegister.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}

	return nil
}
