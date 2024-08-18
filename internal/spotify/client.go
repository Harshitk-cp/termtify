package spotify

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func GetCurrentUser(client *spotify.Client) (*spotify.PrivateUser, error) {
	return client.CurrentUser(context.Background())
}

func GetUserPlaylists(client *spotify.Client) (*spotify.SimplePlaylistPage, error) {
	return client.CurrentUsersPlaylists(context.Background())
}

func GetPlaylistTracks(client *spotify.Client, playlistID spotify.ID) (*spotify.PlaylistItemPage, error) {
	return client.GetPlaylistItems(context.Background(), playlistID)
}

func Search(client *spotify.Client, query string) (*spotify.SearchResult, error) {
	return client.Search(context.Background(), query, spotify.SearchTypeTrack|spotify.SearchTypeAlbum)
}
