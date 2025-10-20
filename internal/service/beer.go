package service

import (
	"backend-test/internal/domain"
	"backend-test/internal/storage/repository"
)

type BeerService struct {
	beerRepository repository.BeerRepositoryInterface
}

func NewBeerService(beerRepo repository.BeerRepositoryInterface) *BeerService {
	return &BeerService{
		beerRepository: beerRepo,
	}
}

func (bs BeerService) ListAllBeerStyles() ([]domain.BeerStyle, error) {
	beerStyles, err := bs.beerRepository.ListAllBeerStyles()
	if err != nil {
		return []domain.BeerStyle{}, err
	}
	return beerStyles, nil
}

func (bs BeerService) GetBeerStyleByUUID(beerUUID string) (domain.BeerStyle, error) {
	beerStyle, err := bs.beerRepository.GetBeerStyleByUUID(beerUUID)
	if err != nil {
		return domain.BeerStyle{}, err
	}
	return beerStyle, nil
}

func (bs BeerService) UpdateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	updatedBeerStyle, err := bs.beerRepository.UpdateBeerStyle(beerStyle)
	if err != nil {
		return domain.BeerStyle{}, err
	}
	return updatedBeerStyle, nil
}

func (bs BeerService) CreateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	createdBeerStyle, err := bs.beerRepository.CreateBeerStyle(beerStyle)
	if err != nil {
		return domain.BeerStyle{}, err
	}
	return createdBeerStyle, nil
}

func (bs BeerService) DeleteBeerStyle(beerUUID string) error {
	err := bs.beerRepository.DeleteBeerStyle(beerUUID)
	if err != nil {
		return err
	}
	return nil
}
