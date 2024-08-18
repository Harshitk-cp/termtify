package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Error loading .env file, falling back to system environment variables")
	}

	config := &Config{
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
	}

	if config.SpotifyClientID == "" || config.SpotifyClientSecret == "" || config.SpotifyRedirectURI == "" {
		return nil, fmt.Errorf("missing required Spotify environment variables")
	}

	return config, nil
}
