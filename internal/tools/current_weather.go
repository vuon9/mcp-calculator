package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vuon9/mcp-play/pkg/config"
	"github.com/vuon9/mcp-play/pkg/httpclient"
)

type Unit string

const (
	Metric   Unit = "metric"
	Imperial Unit = "imperial"
	Standard Unit = "standard"
)

func (u Unit) String() string {
	return string(u)
}

func (u Unit) ToOutputUnit() string {
	switch u {
	case Imperial:
		return "°F"
	case Standard:
		return "K"
	default:
		return "°C"
	}
}

type currentWeather struct {
	httpClient *httpclient.OpenWeatherMapClient
}

func NewCurrentWeatherTool(cfg *config.Config) *currentWeather {
	return &currentWeather{
		httpClient: httpclient.NewOpenWeatherMapClient(
			cfg.OpenWeatherMap.BaseURL,
			cfg.OpenWeatherMap.APIKey,
		),
	}
}

func (w *currentWeather) Definition() mcp.Tool {
	return mcp.NewTool("current_weather",
		mcp.WithString("lat",
			mcp.Description("Latitude of the location (optional, e.g., '37.7749')"),
		),
		mcp.WithString("long",
			mcp.Description("Longitude of the location (optional, e.g., '-122.4194')"),
		),
		mcp.WithString("units",
			mcp.Description("Units for temperature (default: 'metric', options: 'metric', 'imperial', 'standard')"),
			mcp.Enum(string(Metric), string(Imperial), string(Standard)),
			mcp.DefaultString(string(Metric)),
		),
	)
}

func (w *currentWeather) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	lat := request.GetString("lat", "0")
	long := request.GetString("long", "0")
	unit := Unit(request.GetString("units", string(Metric)))

	if lat == "" || long == "" {
		return mcp.NewToolResultError("Latitude and longitude are required"), nil
	}

	weatherResponse, err := w.httpClient.GetCurrentWeather(ctx, httpclient.WithQueryParams(url.Values{
		"lat":   []string{lat},
		"lon":   []string{long},
		"units": []string{unit.String()},
	}))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch weather data: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f%s", weatherResponse.Main.Temp, unit.ToOutputUnit())), nil
}
