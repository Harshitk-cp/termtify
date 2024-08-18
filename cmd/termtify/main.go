package main

import (
	"github.com/Harshitk-cp/termtify/internal/config"
	"github.com/Harshitk-cp/termtify/internal/spotify"
	"github.com/Harshitk-cp/termtify/internal/ui"
	"github.com/Harshitk-cp/termtify/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.Setup()

	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Could not load config: %v", err)
	}

	spotifyClient, err := spotify.AuthSpotify(cfg.SpotifyClientID, cfg.SpotifyClientSecret)
	if err != nil {
		logrus.Fatalf("Spotify authentication failed: %v", err)
	}

	logrus.Info("Authenticated successfully with Spotify!")

	app := ui.NewApp(spotifyClient)
	if err := app.Run(); err != nil {
		logrus.Fatalf("Failed to run TUI: %v", err)
	}
}
