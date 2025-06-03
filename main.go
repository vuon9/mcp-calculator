package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/vuon9/mcp-play/internal/tools"
	"github.com/vuon9/mcp-play/pkg/config"
)

func main() {
	cfg, _ := config.LoadConfig()
	slog.Info("Loaded configuration", "version", cfg.AppVersion)

	s := server.NewMCPServer(
		"Calculator Demo",
		cfg.AppVersion,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
		server.WithLogging(),
	)

	// Add a calculator tool
	calculatorTool := tools.NewCalculateTool()
	s.AddTool(calculatorTool.Definition(), calculatorTool.Handle)

	server := server.NewStreamableHTTPServer(s, server.WithHeartbeatInterval(time.Second*10))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server start error", slog.String("err", err.Error()))
			return
		}
	}()

	<-ctx.Done()
	slog.Info("Received shutdown signal, stopping server...")
	if err := server.Shutdown(context.TODO()); err != nil {
		slog.Error("server shutdown returned an error", slog.String("err", err.Error()))
		return
	}

	slog.Info("Server gracefully stopped")
}
