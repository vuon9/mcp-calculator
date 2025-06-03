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
	slog.Info("Loaded configuration", "config", slog.AnyValue(cfg))

	s := server.NewMCPServer(
		"MCP Play Server",
		cfg.AppVersion,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
		server.WithLogging(),
	)

	tools := []tools.Tool{
		tools.NewCalculateTool(),
		tools.NewCurrentWeatherTool(cfg),
	}

	for _, tool := range tools {
		s.AddTool(tool.Definition(), tool.Handle)
	}

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
