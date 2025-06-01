package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Calculator Demo",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Add a calculator tool
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithArray("numbers",
			mcp.Required(),
			mcp.Description("Array of numbers for operations like sum or average"),
		),
	)

	// Add the calculator handler
	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Using helper functions for type-safe argument access
		op, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		numbers, err := request.RequireFloatSlice("numbers")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var result float64
		switch op {
		case "add":
			for _, num := range numbers {
				result += num
			}
		case "subtract":
			if len(numbers) == 0 {
				return mcp.NewToolResultError("At least one number is required for subtraction"), nil
			}
			result = numbers[0]
			for _, num := range numbers[1:] {
				result -= num
			}
		case "multiply":
			result = 1
			for _, num := range numbers {
				result *= num
			}
		case "divide":
			if len(numbers) == 0 {
				return mcp.NewToolResultError("At least one number is required for division"), nil
			}

			result = numbers[0]
			for _, num := range numbers[1:] {
				if num == 0 {
					return mcp.NewToolResultError("Division by zero is not allowed"), nil
				}
				result /= num
			}
		}

		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	server := server.NewStreamableHTTPServer(s, server.WithHeartbeatInterval(time.Second*10))
	if err := server.Start(":8080"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
