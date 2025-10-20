package repository

import "backend-test/internal/domain"

type BeerRepositoryInterface interface {
	ListAllBeerStyles() ([]domain.BeerStyle, error)
	GetBeerStyleByUUID(beerUUID string) (domain.BeerStyle, error)
	CreateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error)
	UpdateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error)
	DeleteBeerStyle(beerUUID string) error
}
