package main

import (
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/vuon9/mcp-play/internal/tools"
)

func main() {
	s := server.NewMCPServer(
		"Calculator Demo",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Add a calculator tool
	calculatorTool := tools.NewCalculateTool()
	s.AddTool(calculatorTool.Definition(), calculatorTool.Handle)

	server := server.NewStreamableHTTPServer(s, server.WithHeartbeatInterval(time.Second*10))
	if err := server.Start(":8080"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
