package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func PlayTrack(client *spotify.Client, trackID spotify.ID) error {
	devices, err := client.PlayerDevices(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get devices: %v", err)
	}

	if len(devices) == 0 {
		return fmt.Errorf("no active Spotify devices found")
	}

	err = client.PlayOpt(context.Background(), &spotify.PlayOptions{
		URIs: []spotify.URI{spotify.URI("spotify:track:" + string(trackID))},
	})

	if err != nil {
		return fmt.Errorf("failed to play track: %v", err)
	}

	return nil
}
