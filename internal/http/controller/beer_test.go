package controller

import (
	"backend-test/internal/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockBeerService struct {
	beers       []domain.BeerStyle
	shouldError bool
	errorMsg    string
}

func (m *mockBeerService) ListAllBeerStyles() ([]domain.BeerStyle, error) {
	if m.shouldError {
		return nil, &testError{message: m.errorMsg}
	}
	return m.beers, nil
}

func (m *mockBeerService) GetBeerStyleByUUID(beerUUID string) (domain.BeerStyle, error) {
	if m.shouldError {
		return domain.BeerStyle{}, &testError{message: m.errorMsg}
	}
	for _, beer := range m.beers {
		if beer.UUID == beerUUID {
			return beer, nil
		}
	}
	return domain.BeerStyle{}, &testError{message: "not found"}
}

func (m *mockBeerService) CreateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	if m.shouldError {
		return domain.BeerStyle{}, &testError{message: m.errorMsg}
	}
	beerStyle.UUID = "test-uuid"
	beerStyle.CreatedAt = time.Now()
	beerStyle.UpdatedAt = time.Now()
	return beerStyle, nil
}

func (m *mockBeerService) UpdateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	if m.shouldError {
		return domain.BeerStyle{}, &testError{message: m.errorMsg}
	}
	return beerStyle, nil
}

func (m *mockBeerService) DeleteBeerStyle(beerUUID string) error {
	if m.shouldError {
		return &testError{message: m.errorMsg}
	}
	return nil
}

type mockValidationService struct {
	shouldError bool
	errorMsg    string
}

func (m *mockValidationService) ValidateTemperatureRange(beerStyle domain.BeerStyle) error {
	if m.shouldError {
		return &testError{message: m.errorMsg}
	}
	return nil
}

func (m *mockValidationService) ValidateTemperatureInput(temperature float64) error {
	if m.shouldError {
		return &testError{message: m.errorMsg}
	}
	return nil
}

func (m *mockValidationService) ValidateUniqueNameForCreate(name string) error {
	if m.shouldError {
		return &testError{message: m.errorMsg}
	}
	return nil
}

func (m *mockValidationService) ValidateUniqueNameForUpdate(name string, excludeUUID string) error {
	return nil
}

func (m *mockValidationService) IsNoRowsError(err error) bool {
	return false
}

func (m *mockValidationService) ValidateUUID(uuidStr string) error {
	return nil
}

type mockUpdateService struct{}

func (m *mockUpdateService) ApplyBeerStyleUpdates(current *domain.BeerStyle, updates domain.BeerStyleUpdateRequest) bool {
	return false
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

func setupController() *BeerController {
	beerService := &mockBeerService{beers: []domain.BeerStyle{}}
	validationService := &mockValidationService{}
	updateService := &mockUpdateService{}

	return NewBeerController(beerService, validationService, updateService)
}

func TestBeerController_NewBeerController(t *testing.T) {
	controller := setupController()

	if controller == nil {
		t.Error("Controller should not be nil")
	}

	if controller.BeerService == nil {
		t.Error("BeerService should not be nil")
	}

	if controller.ValidationService == nil {
		t.Error("ValidationService should not be nil")
	}

	if controller.UpdateService == nil {
		t.Error("UpdateService should not be nil")
	}
}

func TestBeerController_ListAllBeerStyles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBeers := []domain.BeerStyle{
		{UUID: "1", Name: "IPA", TempMin: 5, TempMax: 8},
		{UUID: "2", Name: "Lager", TempMin: -2, TempMax: 2},
	}

	beerService := &mockBeerService{beers: mockBeers, shouldError: false}
	validationService := &mockValidationService{}
	updateService := &mockUpdateService{}

	controller := NewBeerController(beerService, validationService, updateService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.ListAllBeerStyles(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
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
}

func TestBeerController_ListAllBeerStyles_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	beerService := &mockBeerService{shouldError: true, errorMsg: "database error"}
	validationService := &mockValidationService{}
	updateService := &mockUpdateService{}

	controller := NewBeerController(beerService, validationService, updateService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.ListAllBeerStyles(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestBeerController_CreateBeerStyle_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	beerService := &mockBeerService{shouldError: false}
	validationService := &mockValidationService{shouldError: false}
	updateService := &mockUpdateService{}

	controller := NewBeerController(beerService, validationService, updateService)

	requestBody := map[string]interface{}{
		"name":     "Test Beer",
		"temp_min": 5.0,
		"temp_max": 8.0,
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateBeerStyle(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestBeerController_UpdateBeerStyle_Success(t *testing.T) {
	beerService := &mockBeerService{
		beers: []domain.BeerStyle{
			{
				UUID:      "test-uuid-1",
				Name:      "Test IPA",
				TempMin:   4.0,
				TempMax:   7.0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		shouldError: false,
	}
	validationService := &mockValidationService{}
	updateService := &mockUpdateService{}

	controller := NewBeerController(beerService, validationService, updateService)

	requestBody := map[string]interface{}{
		"name": "Updated Beer",
	}

	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "beerUUID", Value: "test-uuid-1"}}

	controller.UpdateBeerStyle(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBeerController_DeleteBeerStyle_Success(t *testing.T) {
	beerService := &mockBeerService{
		beers: []domain.BeerStyle{
			{
				UUID:      "test-uuid-1",
				Name:      "Test IPA",
				TempMin:   4.0,
				TempMax:   7.0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		shouldError: false,
	}
	validationService := &mockValidationService{}
	updateService := &mockUpdateService{}

	controller := NewBeerController(beerService, validationService, updateService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "beerUUID", Value: "test-uuid-1"}}

	controller.DeleteBeerStyle(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Beer style deleted successfully" {
		t.Errorf("Expected message 'Beer style deleted successfully', got '%s'", response["message"])
	}
}

func TestBeerController_DeleteBeerStyle_ServiceError(t *testing.T) {
	beerService := &mockBeerService{
		shouldError: true,
		errorMsg:    "delete failed",
	}
	validationService := &mockValidationService{}
	updateService := &mockUpdateService{}

	controller := NewBeerController(beerService, validationService, updateService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "beerUUID", Value: "test-uuid-1"}}

	controller.DeleteBeerStyle(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
