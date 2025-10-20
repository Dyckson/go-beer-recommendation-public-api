package repository

import (
	"backend-test/internal/domain"
	postgres "backend-test/internal/storage/database"
	"context"
	"time"
)

type BeerRepository struct{}

func (u BeerRepository) ListAllBeerStyles() ([]domain.BeerStyle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := postgres.GetDB()

	var beerStyles []domain.BeerStyle
	err := db.Query(ctx, &beerStyles, u.getAllBeerStylesQuery())
	if err != nil {
		return nil, err
	}

	return beerStyles, nil
}

func (u BeerRepository) CreateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := postgres.GetDB()

	var createdBeerStyle domain.BeerStyle
	err := db.QueryOne(ctx, &createdBeerStyle, u.createBeerStyleQuery(), beerStyle.Name, beerStyle.TempMin, beerStyle.TempMax)
	if err != nil {
		return domain.BeerStyle{}, err
	}

	return createdBeerStyle, nil
}

func (u BeerRepository) GetBeerStyleByUUID(beerUUID string) (domain.BeerStyle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := postgres.GetDB()

	var beerStyle domain.BeerStyle
	err := db.QueryOne(ctx, &beerStyle, u.getBeerStyleByUUIDQuery(), beerUUID)
	if err != nil {
		return domain.BeerStyle{}, err
	}

	return beerStyle, nil
}

func (u BeerRepository) UpdateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := postgres.GetDB()

	var updatedBeerStyle domain.BeerStyle
	err := db.QueryOne(ctx, &updatedBeerStyle, u.updateBeerStyleQuery(),
		beerStyle.Name, beerStyle.TempMin, beerStyle.TempMax, beerStyle.UUID)
	if err != nil {
		return domain.BeerStyle{}, err
	}

	return updatedBeerStyle, nil
}

func (u BeerRepository) DeleteBeerStyle(beerUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := postgres.GetDB()

	_, err := db.Exec(ctx, u.deleteBeerStyleQuery(), beerUUID)
	if err != nil {
		return err
	}

	return nil
}

func (BeerRepository) getAllBeerStylesQuery() string {
	return `
		SELECT uuid, name, temp_min, temp_max, created_at, updated_at
		FROM beer_styles
	`
}

func (BeerRepository) getBeerStyleByUUIDQuery() string {
	return `
		SELECT uuid, name, temp_min, temp_max, created_at, updated_at
		FROM beer_styles
		WHERE uuid = $1
	`
}

func (BeerRepository) createBeerStyleQuery() string {
	return `
		INSERT INTO beer_styles (name, temp_min, temp_max)
		VALUES ($1, $2, $3)
		RETURNING uuid, name, temp_min, temp_max, created_at, updated_at;
	`
}

func (BeerRepository) updateBeerStyleQuery() string {
	return `
		UPDATE beer_styles
		SET name = $1,
		temp_min = $2,
		temp_max = $3,
		updated_at = NOW()
		WHERE uuid = $4
		RETURNING uuid, name, temp_min, temp_max, created_at, updated_at;
	`
}

func (BeerRepository) deleteBeerStyleQuery() string {
	return `
		DELETE FROM beer_styles
		WHERE uuid = $1;
	`
}
