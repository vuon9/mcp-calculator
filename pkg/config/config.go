package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppVersion           string `env:"APP_VERSION" envDefault:"0.0.0"`
	OpenWeatherMapApiKey string `env:"OPEN_WEATHER_MAP_API_KEY" envDefault:""`
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
