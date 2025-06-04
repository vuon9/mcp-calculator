package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// NewCalculateTool creates a new instance of the Calculate tool.
func Test_NewCalculateTool(t *testing.T) {
	c := NewCalculateTool()
	if c == nil {
		t.Fatal("NewCalculateTool() returned nil")
	}
	if c.Definition().Name != "calculate" {
		t.Errorf("NewCalculateTool() returned tool with unexpected name: %s", c.Definition().Name)
	}
}

func Test_calculate_Handle(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		request mcp.CallToolRequest
		want    *mcp.CallToolResult
		wantErr bool
	}{
		{
			name:    "missing operation parameter",
			request: mcp.CallToolRequest{},
			want:    mcp.NewToolResultError(errors.New("required argument \"operation\" not found").Error()),
			wantErr: false,
		},
		{
			name: "missing numbers parameter",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "add",
					},
				},
			},
			want:    mcp.NewToolResultError(errors.New("required argument \"numbers\" not found").Error()),
			wantErr: false,
		},
		{
			name: "invalid operation",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "invalid",
						"numbers":   []float64{1, 2, 3},
					},
				},
			},
			want:    mcp.NewToolResultError("invalid operation"),
			wantErr: false,
		},
		{
			name: "add operation",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "add",
						"numbers":   []float64{1, 2, 3},
					},
				},
			},
			want:    mcp.NewToolResultText("6.00"),
			wantErr: false,
		},
		{
			name: "subtract operation",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "subtract",
						"numbers":   []float64{3, 2, 1},
					},
				},
			},
			want:    mcp.NewToolResultText("0.00"),
			wantErr: false,
		},
		{
			name: "subtract operation with no numbers",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "subtract",
						"numbers":   []float64{},
					},
				},
			},
			want:    mcp.NewToolResultError("At least one number is required for subtraction"),
			wantErr: false,
		},
		{
			name: "multiply operation",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "multiply",
						"numbers":   []float64{3, 3, 1},
					},
				},
			},
			want:    mcp.NewToolResultText("9.00"),
			wantErr: false,
		},
		{
			name: "divide operation",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "divide",
						"numbers":   []float64{9, 3, 1},
					},
				},
			},
			want:    mcp.NewToolResultText("3.00"),
			wantErr: false,
		},
		{
			name: "divide operation without numbers",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "divide",
						"numbers":   []float64{},
					},
				},
			},
			want:    mcp.NewToolResultError("At least one number is required for division"),
			wantErr: false,
		},
		{
			name: "divide operation with zero",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"operation": "divide",
						"numbers":   []float64{10, 0, 2},
					},
				},
			},
			want:    mcp.NewToolResultError("Division by zero is not allowed"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCalculateTool()
			got, gotErr := c.Handle(context.Background(), tt.request)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Handle() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Handle() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
