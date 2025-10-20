package domain


type TemperatureRequest struct {
	Temperature float64 `json:"temperature"`
}

type TrackInfo struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Link   string `json:"link"`
}

type PlaylistInfo struct {
	Name   string      `json:"name"`
	Tracks []TrackInfo `json:"tracks"`
}

type RecommendationResponse struct {
	BeerStyle string       `json:"beerStyle"`
	Playlist  PlaylistInfo `json:"playlist"`
}
