package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

type calculate struct{}

func NewCalculateTool() *calculate {
	return &calculate{}
}

func (c *calculate) Definition() mcp.Tool {
	return mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithArray("numbers",
			mcp.Required(),
			mcp.Description("Array of numbers for operation"),
		),
	)
}

func (c *calculate) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	default:
		return mcp.NewToolResultError("invalid operation"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
}
