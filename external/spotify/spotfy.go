package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyService struct {
	client *spotify.Client
}

func NewSpotifyService(clientID, clientSecret string) (*SpotifyService, error) {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		return nil, err
	}
	httpClient := spotifyauth.New().Client(context.Background(), token)
	client := spotify.New(httpClient)
	return &SpotifyService{client: client}, nil
}

func (s *SpotifyService) SearchPlaylistByName(name string) (*spotify.FullPlaylist, error) {
	results, err := s.client.Search(context.Background(), name, spotify.SearchTypePlaylist)
	if err != nil {
		return nil, err
	}
	if results.Playlists == nil || len(results.Playlists.Playlists) == 0 {
		return nil, fmt.Errorf("playlist not found")
	}
	playlistID := results.Playlists.Playlists[0].ID
	fullPlaylist, err := s.client.GetPlaylist(context.Background(), playlistID)
	if err != nil {
		return nil, err
	}
	return fullPlaylist, nil
}
