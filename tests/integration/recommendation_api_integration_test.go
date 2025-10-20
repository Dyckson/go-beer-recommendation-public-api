package integration

import (
	"backend-test/internal/domain"
	"backend-test/internal/http/controller"
	"backend-test/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockRecommendationService struct {
	shouldError bool
	errorMsg    string
	response    domain.RecommendationResponse
}

func (m *MockRecommendationService) GetRecommendationForTemperature(temperature float64) (*domain.RecommendationResponse, error) {
	if m.shouldError {
		return nil, &MockError{message: m.errorMsg}
	}
	return &m.response, nil
}

func setupRecommendationTestRouter(recommendationService service.RecommendationServiceInterface, validationService service.ValidationServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)

	recommendationController := controller.NewRecommendationController(recommendationService, validationService)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		recommendations := api.Group("/recommendations")
		{
			recommendations.POST("/suggest", recommendationController.SuggestSpotifyPlaylist)
		}
	}

	return r
}

func TestRecommendationAPI_SuggestSpotifyPlaylist_Success(t *testing.T) {
	mockRecommendationService := &MockRecommendationService{
		shouldError: false,
		response: domain.RecommendationResponse{
			BeerStyle: "IPA",
			Playlist: domain.PlaylistInfo{
				Name: "Rock Playlist for IPA",
				Tracks: []domain.TrackInfo{
					{
						Name:   "Bohemian Rhapsody",
						Artist: "Queen",
						Link:   "https://spotify.com/track/123",
					},
					{
						Name:   "Stairway to Heaven",
						Artist: "Led Zeppelin",
						Link:   "https://spotify.com/track/456",
					},
				},
			},
		},
	}

	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}

	testRouter := setupRecommendationTestRouter(mockRecommendationService, mockValidationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/recommendations/suggest", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response domain.RecommendationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.BeerStyle != "IPA" {
		t.Errorf("Expected beer style 'IPA', got '%s'", response.BeerStyle)
	}

	if response.Playlist.Name != "Rock Playlist for IPA" {
		t.Errorf("Expected playlist name 'Rock Playlist for IPA', got '%s'", response.Playlist.Name)
	}

	if len(response.Playlist.Tracks) != 2 {
		t.Errorf("Expected 2 tracks, got %d", len(response.Playlist.Tracks))
	}

	if response.Playlist.Tracks[0].Name != "Bohemian Rhapsody" {
		t.Errorf("Expected first track 'Bohemian Rhapsody', got '%s'", response.Playlist.Tracks[0].Name)
	}
}

func TestRecommendationAPI_SuggestSpotifyPlaylist_InvalidJSON(t *testing.T) {
	mockRecommendationService := &MockRecommendationService{shouldError: false}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}

	testRouter := setupRecommendationTestRouter(mockRecommendationService, mockValidationService)

	invalidJSON := `{"temperature": "not-a-number"}`

	req, _ := http.NewRequest("POST", "/api/recommendations/suggest", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
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

func TestRecommendationAPI_SuggestSpotifyPlaylist_ValidationError(t *testing.T) {
	mockRecommendationService := &MockRecommendationService{shouldError: false}
	mockValidationService := &MockValidationService{
		shouldError:     true,
		uniqueNameError: false,
		tempRangeError:  false,
		errorMsg:        "temperature must be between -10 and 50 degrees",
	}

	testRouter := setupRecommendationTestRouter(mockRecommendationService, mockValidationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 100.0, // Invalid temperature
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/recommendations/suggest", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "temperature must be between -10 and 50 degrees" {
		t.Errorf("Expected message 'temperature must be between -10 and 50 degrees', got '%s'", response["message"])
	}
}

func TestRecommendationAPI_SuggestSpotifyPlaylist_NoPlaylistFound(t *testing.T) {
	mockRecommendationService := &MockRecommendationService{
		shouldError: true,
		errorMsg:    "no playlist found for temperature",
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}

	testRouter := setupRecommendationTestRouter(mockRecommendationService, mockValidationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 15.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/recommendations/suggest", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
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

func TestRecommendationAPI_SuggestSpotifyPlaylist_SpotifyUnavailable(t *testing.T) {
	mockRecommendationService := &MockRecommendationService{
		shouldError: true,
		errorMsg:    "spotify service unavailable",
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}

	testRouter := setupRecommendationTestRouter(mockRecommendationService, mockValidationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/recommendations/suggest", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, w.Code)
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

func TestRecommendationAPI_SuggestSpotifyPlaylist_BeerStyleError(t *testing.T) {
	mockRecommendationService := &MockRecommendationService{
		shouldError: true,
		errorMsg:    "failed to find best beer style",
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}

	testRouter := setupRecommendationTestRouter(mockRecommendationService, mockValidationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/recommendations/suggest", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
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

func TestRecommendationAPI_SuggestSpotifyPlaylist_GenericError(t *testing.T) {
	mockRecommendationService := &MockRecommendationService{
		shouldError: true,
		errorMsg:    "unexpected database error",
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}

	testRouter := setupRecommendationTestRouter(mockRecommendationService, mockValidationService)

	requestBody := domain.TemperatureRequest{
		Temperature: 6.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/recommendations/suggest", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
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
