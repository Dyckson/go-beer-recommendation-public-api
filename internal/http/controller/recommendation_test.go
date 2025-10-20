package controller

import (
	"backend-test/internal/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockRecommendationService struct {
	shouldError bool
	errorMsg    string
	response    domain.RecommendationResponse
}

func (m *mockRecommendationService) GetRecommendationForTemperature(temperature float64) (*domain.RecommendationResponse, error) {
	if m.shouldError {
		return nil, &testError{message: m.errorMsg}
	}
	return &m.response, nil
}

func setupRecommendationTestController() *RecommendationController {
	gin.SetMode(gin.TestMode)

	recommendationService := &mockRecommendationService{
		response: domain.RecommendationResponse{
			BeerStyle: "IPA",
			Playlist: domain.PlaylistInfo{
				Name: "Test Playlist",
				Tracks: []domain.TrackInfo{
					{
						Name:   "Test Song",
						Artist: "Test Artist",
						Link:   "https://spotify.com/track/123",
					},
				},
			},
		},
	}

	validationService := &mockValidationService{}

	return NewRecommendationController(recommendationService, validationService)
}

func TestRecommendationController_NewRecommendationController(t *testing.T) {
	controller := setupRecommendationTestController()

	if controller == nil {
		t.Error("Expected controller not to be nil")
	}
	if controller.RecommendationService == nil {
		t.Error("Expected RecommendationService not to be nil")
	}
	if controller.ValidationService == nil {
		t.Error("Expected ValidationService not to be nil")
	}
}

func TestRecommendationController_SuggestSpotifyPlaylist_Success(t *testing.T) {
	controller := setupRecommendationTestController()

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.SuggestSpotifyPlaylist(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response domain.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.BeerStyle != "IPA" {
		t.Errorf("Expected beer style 'IPA', got '%s'", response.BeerStyle)
	}

	if response.Playlist.Name != "Test Playlist" {
		t.Errorf("Expected playlist name 'Test Playlist', got '%s'", response.Playlist.Name)
	}
}

func TestRecommendationController_SuggestSpotifyPlaylist_InvalidJSON(t *testing.T) {
	controller := setupRecommendationTestController()

	invalidJSON := `{"temperature": "invalid"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.SuggestSpotifyPlaylist(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "invalid request body format" {
		t.Errorf("Expected message 'invalid request body format', got '%s'", response["message"])
	}
}

func TestRecommendationController_SuggestSpotifyPlaylist_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recommendationService := &mockRecommendationService{}
	validationService := &mockValidationService{
		shouldError: true,
		errorMsg:    "temperature out of range",
	}

	controller := NewRecommendationController(recommendationService, validationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 100.0, // Invalid temperature
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.SuggestSpotifyPlaylist(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "temperature out of range" {
		t.Errorf("Expected message 'temperature out of range', got '%s'", response["message"])
	}
}

func TestRecommendationController_SuggestSpotifyPlaylist_NoPlaylistFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recommendationService := &mockRecommendationService{
		shouldError: true,
		errorMsg:    "no playlist found for temperature",
	}
	validationService := &mockValidationService{}

	controller := NewRecommendationController(recommendationService, validationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.SuggestSpotifyPlaylist(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "no playlist found for temperature" {
		t.Errorf("Expected message 'no playlist found for temperature', got '%s'", response["message"])
	}
}

func TestRecommendationController_SuggestSpotifyPlaylist_SpotifyUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recommendationService := &mockRecommendationService{
		shouldError: true,
		errorMsg:    "spotify service unavailable",
	}
	validationService := &mockValidationService{}

	controller := NewRecommendationController(recommendationService, validationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.SuggestSpotifyPlaylist(c)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Spotify service is temporarily unavailable" {
		t.Errorf("Expected message 'Spotify service is temporarily unavailable', got '%s'", response["message"])
	}
}

func TestRecommendationController_SuggestSpotifyPlaylist_BeerStyleError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recommendationService := &mockRecommendationService{
		shouldError: true,
		errorMsg:    "failed to find best beer style",
	}
	validationService := &mockValidationService{}

	controller := NewRecommendationController(recommendationService, validationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.SuggestSpotifyPlaylist(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Unable to determine suitable beer style" {
		t.Errorf("Expected message 'Unable to determine suitable beer style', got '%s'", response["message"])
	}
}

func TestRecommendationController_SuggestSpotifyPlaylist_GenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recommendationService := &mockRecommendationService{
		shouldError: true,
		errorMsg:    "some unexpected error",
	}
	validationService := &mockValidationService{}

	controller := NewRecommendationController(recommendationService, validationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.SuggestSpotifyPlaylist(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Internal server error" {
		t.Errorf("Expected message 'Internal server error', got '%s'", response["message"])
	}
}
