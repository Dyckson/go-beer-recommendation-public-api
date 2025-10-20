package service

import (
	"backend-test/external/spotify"
	config "backend-test/internal/cmd/server"
	"backend-test/internal/domain"
	"fmt"
	"log"

	spotifyapi "github.com/zmb3/spotify/v2"
)

type RecommendationService struct {
	beerService    BeerServiceInterface
	spotifyService *spotify.SpotifyService
}

func NewRecommendationService(beerService BeerServiceInterface, spotifyService *spotify.SpotifyService) *RecommendationService {
	return &RecommendationService{
		beerService:    beerService,
		spotifyService: spotifyService,
	}
}

func (rs *RecommendationService) FindBestBeerStyleForTemperature(temperature float64) (*domain.BeerStyle, error) {
	allBeerStyles, err := rs.beerService.ListAllBeerStyles()
	if err != nil {
		return nil, fmt.Errorf("failed to get beer styles: %w", err)
	}

	if len(allBeerStyles) == 0 {
		return nil, fmt.Errorf("no beer styles found")
	}

	type candidateStyle struct {
		style    *domain.BeerStyle
		distance float64
		average  float64
	}

	var candidates []candidateStyle
	var minDistance float64 = 999999

	for i := range allBeerStyles {
		beerStyle := &allBeerStyles[i]
		average := (beerStyle.TempMin + beerStyle.TempMax) / 2
		distance := abs(temperature - average)

		if distance < minDistance {
			minDistance = distance
			candidates = []candidateStyle{{
				style:    beerStyle,
				distance: distance,
				average:  average,
			}}
		} else if distance == minDistance {
			candidates = append(candidates, candidateStyle{
				style:    beerStyle,
				distance: distance,
				average:  average,
			})
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no suitable beer style found for temperature %.1f°C", temperature)
	}

	if len(candidates) == 1 {
		return candidates[0].style, nil
	}

	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].style.Name > candidates[j].style.Name {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	bestMatch := candidates[0].style
	return bestMatch, nil
}
func (rs *RecommendationService) GetRecommendationForTemperature(temperature float64) (*domain.RecommendationResponse, error) {
	beerStyle, err := rs.FindBestBeerStyleForTemperature(temperature)
	if err != nil {
		return nil, err
	}

	var tracks []domain.TrackInfo
	var playlistName string

	spotifyService := config.GetSpotifyService()

	if spotifyService != nil {
		playlist, err := spotifyService.SearchPlaylistByName(beerStyle.Name)
		if err != nil {
			log.Printf("Failed to find Spotify playlist for %s: %v", beerStyle.Name, err)
			return nil, fmt.Errorf("no playlist found for beer style '%s'", beerStyle.Name)
		}

		playlistName = playlist.Name
		tracks = rs.convertSpotifyTracks(playlist)

		if len(tracks) == 0 {
			return nil, fmt.Errorf("playlist '%s' found but contains no valid tracks", playlistName)
		}
	} else {
		log.Println("Spotify service not available")
		return nil, fmt.Errorf("spotify service unavailable")
	}

	response := &domain.RecommendationResponse{
		BeerStyle: beerStyle.Name,
		Playlist: domain.PlaylistInfo{
			Name:   playlistName,
			Tracks: tracks,
		},
	}

	return response, nil
}

func (rs *RecommendationService) convertSpotifyTracks(playlist *spotifyapi.FullPlaylist) []domain.TrackInfo {
	tracks := make([]domain.TrackInfo, 0)
	if len(playlist.Tracks.Tracks) > 0 {
		maxTracks := len(playlist.Tracks.Tracks)
		if maxTracks > 10 {
			maxTracks = 10
		}

		for i := 0; i < maxTracks; i++ {
			track := playlist.Tracks.Tracks[i].Track
			if track.Name != "" {
				spotifyLink := fmt.Sprintf("https://open.spotify.com/track/%s", track.ID)

				artistName := "Unknown Artist"
				if len(track.Artists) > 0 {
					artistName = track.Artists[0].Name
				}

				tracks = append(tracks, domain.TrackInfo{
					Name:   track.Name,
					Artist: artistName,
					Link:   spotifyLink,
				})
			}
		}
	}
	return tracks
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
