package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type OpenWeatherMapConfig struct {
	BaseURL string `env:"BASE_URL" envDefault:"https://api.openweathermap.org/data/2.5/"`
	APIKey  string `env:"API_KEY" envDefault:""`
}

type Config struct {
	AppVersion string `env:"APP_VERSION" envDefault:"0.0.0"`

	OpenWeatherMap OpenWeatherMapConfig `envPrefix:"OPEN_WEATHER_MAP_" envDefault:""`
}

const defaultConfigFile = ".env"

func LoadConfig(filenames ...string) (*Config, error) {
	if len(filenames) == 0 {
		filenames = []string{defaultConfigFile}
	}

	err := godotenv.Load(filenames...)
	if err != nil {
		return nil, err
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
