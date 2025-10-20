package integration

import (
	"backend-test/internal/domain"
	"backend-test/internal/http/controller"
	"backend-test/internal/http/router"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockBeerService struct {
	beers       []domain.BeerStyle
	shouldError bool
	errorMsg    string
}

func (m *MockBeerService) ListAllBeerStyles() ([]domain.BeerStyle, error) {
	if m.shouldError {
		return nil, &MockError{message: m.errorMsg}
	}
	return m.beers, nil
}

func (m *MockBeerService) GetBeerStyleByUUID(beerUUID string) (domain.BeerStyle, error) {
	if m.shouldError {
		return domain.BeerStyle{}, &MockError{message: m.errorMsg}
	}
	for _, beer := range m.beers {
		if beer.UUID == beerUUID {
			return beer, nil
		}
	}
	return domain.BeerStyle{}, &MockError{message: "beer not found"}
}

func (m *MockBeerService) CreateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	if m.shouldError {
		return domain.BeerStyle{}, &MockError{message: m.errorMsg}
	}
	beerStyle.UUID = "test-uuid-123"
	beerStyle.CreatedAt = time.Now()
	beerStyle.UpdatedAt = time.Now()
	m.beers = append(m.beers, beerStyle)
	return beerStyle, nil
}

func (m *MockBeerService) UpdateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	if m.shouldError {
		return domain.BeerStyle{}, &MockError{message: m.errorMsg}
	}
	return beerStyle, nil
}

func (m *MockBeerService) DeleteBeerStyle(beerUUID string) error {
	if m.shouldError {
		return &MockError{message: m.errorMsg}
	}
	return nil
}

type MockValidationService struct {
	shouldError     bool
	errorMsg        string
	uniqueNameError bool
	tempRangeError  bool
}

func (m *MockValidationService) ValidateTemperatureRange(beerStyle domain.BeerStyle) error {
	if m.tempRangeError {
		return &MockError{message: "temperature range invalid"}
	}
	return nil
}

func (m *MockValidationService) ValidateTemperatureInput(temperature float64) error {
	if m.shouldError {
		return &MockError{message: m.errorMsg}
	}
	return nil
}

func (m *MockValidationService) ValidateUniqueNameForCreate(name string) error {
	if m.uniqueNameError {
		return &MockError{message: m.errorMsg}
	}
	return nil
}

func (m *MockValidationService) ValidateUniqueNameForUpdate(name string, excludeUUID string) error {
	if m.shouldError {
		return &MockError{message: m.errorMsg}
	}
	return nil
}

func (m *MockValidationService) IsNoRowsError(err error) bool {
	return false
}

func (m *MockValidationService) ValidateUUID(uuidStr string) error {
	if m.shouldError {
		return &MockError{message: m.errorMsg}
	}
	return nil
}

type MockUpdateService struct {
	shouldError bool
	errorMsg    string
}

func (m *MockUpdateService) ApplyBeerStyleUpdates(current *domain.BeerStyle, updates domain.BeerStyleUpdateRequest) bool {
	return true
}

type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func setupTestController(mockBeer *MockBeerService, mockValidation *MockValidationService, mockUpdate *MockUpdateService) *controller.BeerController {
	return controller.NewBeerController(mockBeer, mockValidation, mockUpdate)
}

func setupTestRouter(beerController *controller.BeerController) *gin.Engine {
	r := router.NewRouter()
	api := r.Group("/api")

	beer := api.Group("/beer-styles")
	beer.GET("/list", beerController.ListAllBeerStyles)
	beer.POST("/create", beerController.CreateBeerStyle)
	beer.PUT("/edit/:beerUUID", beerController.UpdateBeerStyle)
	beer.DELETE("/:beerUUID", beerController.DeleteBeerStyle)

	return r
}

func TestBeerAPI_ListAllBeerStyles_Success(t *testing.T) {
	mockBeers := []domain.BeerStyle{
		{
			UUID:      "uuid-1",
			Name:      "IPA",
			TempMin:   5.0,
			TempMax:   8.0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:      "uuid-2",
			Name:      "Lager",
			TempMin:   -2.0,
			TempMax:   2.0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockBeerService := &MockBeerService{
		beers:       mockBeers,
		shouldError: false,
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}
	mockUpdateService := &MockUpdateService{shouldError: false}

	beerController := setupTestController(mockBeerService, mockValidationService, mockUpdateService)
	testRouter := setupTestRouter(beerController)

	req, _ := http.NewRequest("GET", "/api/beer-styles/list", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		BeerStyles []domain.BeerStyle `json:"beerStyles"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.BeerStyles) != 2 {
		t.Errorf("Expected 2 beer styles, got %d", len(response.BeerStyles))
	}

	if response.BeerStyles[0].Name != "IPA" {
		t.Errorf("Expected first beer to be 'IPA', got '%s'", response.BeerStyles[0].Name)
	}

	if response.BeerStyles[1].Name != "Lager" {
		t.Errorf("Expected second beer to be 'Lager', got '%s'", response.BeerStyles[1].Name)
	}
}

func TestBeerAPI_ListAllBeerStyles_DatabaseError(t *testing.T) {
	mockBeerService := &MockBeerService{
		beers:       nil,
		shouldError: true,
		errorMsg:    "database connection failed",
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}
	mockUpdateService := &MockUpdateService{shouldError: false}

	beerController := setupTestController(mockBeerService, mockValidationService, mockUpdateService)
	testRouter := setupTestRouter(beerController)

	req, _ := http.NewRequest("GET", "/api/beer-styles/list", nil)
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

	if response["message"] != "internal error" {
		t.Errorf("Expected message 'internal error', got '%s'", response["message"])
	}
}

func TestBeerAPI_ListAllBeerStyles_EmptyResult(t *testing.T) {
	mockBeerService := &MockBeerService{
		beers:       []domain.BeerStyle{},
		shouldError: false,
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}
	mockUpdateService := &MockUpdateService{shouldError: false}

	beerController := setupTestController(mockBeerService, mockValidationService, mockUpdateService)
	testRouter := setupTestRouter(beerController)

	req, _ := http.NewRequest("GET", "/api/beer-styles/list", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		BeerStyles []domain.BeerStyle `json:"beerStyles"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.BeerStyles) != 0 {
		t.Errorf("Expected empty beer styles array, got %d items", len(response.BeerStyles))
	}
}

func TestBeerAPI_CreateBeerStyle_Success(t *testing.T) {
	mockBeerService := &MockBeerService{
		beers:       []domain.BeerStyle{},
		shouldError: false,
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}
	mockUpdateService := &MockUpdateService{shouldError: false}

	beerController := setupTestController(mockBeerService, mockValidationService, mockUpdateService)
	testRouter := setupTestRouter(beerController)

	requestBody := map[string]interface{}{
		"name":     "Test IPA",
		"temp_min": 5.0,
		"temp_max": 8.0,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/beer-styles/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response struct {
		Data domain.BeerStyle `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Data.Name != "Test IPA" {
		t.Errorf("Expected beer name 'Test IPA', got '%s'", response.Data.Name)
	}

	if response.Data.UUID == "" {
		t.Error("Expected UUID to be generated")
	}
}

func TestBeerAPI_CreateBeerStyle_ValidationError(t *testing.T) {
	mockBeerService := &MockBeerService{
		beers:       []domain.BeerStyle{},
		shouldError: false,
	}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  true,
		errorMsg:        "temperature range invalid",
	}
	mockUpdateService := &MockUpdateService{shouldError: false}

	beerController := setupTestController(mockBeerService, mockValidationService, mockUpdateService)
	testRouter := setupTestRouter(beerController)

	requestBody := map[string]interface{}{
		"name":     "Invalid Beer",
		"temp_min": 10.0,
		"temp_max": 5.0, // Invalid: min > max
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/beer-styles/create", bytes.NewBuffer(jsonBody))
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

	if response["message"] != "temperature range invalid" {
		t.Errorf("Expected message 'temperature range invalid', got '%s'", response["message"])
	}
}

func TestBeerAPI_CreateBeerStyle_InvalidJSON(t *testing.T) {
	mockBeerService := &MockBeerService{beers: []domain.BeerStyle{}, shouldError: false}
	mockValidationService := &MockValidationService{
		shouldError:     false,
		uniqueNameError: false,
		tempRangeError:  false,
	}
	mockUpdateService := &MockUpdateService{shouldError: false}

	beerController := setupTestController(mockBeerService, mockValidationService, mockUpdateService)
	testRouter := setupTestRouter(beerController)

	invalidJSON := `{"name": "Test Beer", "temp_min": invalid}`

	req, _ := http.NewRequest("POST", "/api/beer-styles/create", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
