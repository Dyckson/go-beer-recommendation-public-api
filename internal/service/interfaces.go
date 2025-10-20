package service

import "backend-test/internal/domain"

type BeerServiceInterface interface {
	ListAllBeerStyles() ([]domain.BeerStyle, error)
	GetBeerStyleByUUID(beerUUID string) (domain.BeerStyle, error)
	CreateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error)
	UpdateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error)
	DeleteBeerStyle(beerUUID string) error
}

type ValidationServiceInterface interface {
	ValidateTemperatureRange(beerStyle domain.BeerStyle) error
	ValidateTemperatureInput(temperature float64) error
	ValidateUniqueNameForCreate(name string) error
	ValidateUniqueNameForUpdate(name string, excludeUUID string) error
	IsNoRowsError(err error) bool
	ValidateUUID(uuidStr string) error
}

type UpdateServiceInterface interface {
	ApplyBeerStyleUpdates(current *domain.BeerStyle, updates domain.BeerStyleUpdateRequest) bool
}

type RecommendationServiceInterface interface {
	GetRecommendationForTemperature(temperature float64) (*domain.RecommendationResponse, error)
}
