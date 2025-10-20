package service

import (
	"backend-test/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ValidationService struct {
	beerService BeerServiceInterface
}

func NewValidationService(beerService BeerServiceInterface) *ValidationService {
	return &ValidationService{
		beerService: beerService,
	}
}

func (vs *ValidationService) isNoRowsError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set")
}

func (vs *ValidationService) IsNoRowsError(err error) bool {
	return vs.isNoRowsError(err)
}

func (vs *ValidationService) ValidateTemperatureRange(beerStyle domain.BeerStyle) error {
	if beerStyle.TempMin < -90 || beerStyle.TempMin > 60 {
		return fmt.Errorf("minimum temperature (%.1f) must be between -90°C and 60°C", beerStyle.TempMin)
	}

	if beerStyle.TempMax < -90 || beerStyle.TempMax > 60 {
		return fmt.Errorf("maximum temperature (%.1f) must be between -90°C and 60°C", beerStyle.TempMax)
	}

	if beerStyle.TempMin >= beerStyle.TempMax {
		return fmt.Errorf("minimum temperature (%.1f) must be less than maximum temperature (%.1f)",
			beerStyle.TempMin, beerStyle.TempMax)
	}

	return nil
}

func (vs *ValidationService) ValidateTemperatureInput(temperature float64) error {
	if temperature < -90 || temperature > 60 {
		return fmt.Errorf("temperature (%.1f) must be between -90°C and 60°C", temperature)
	}

	return nil
}

func (vs *ValidationService) ValidateUniqueNameForCreate(name string) error {
	beerStyles, err := vs.beerService.ListAllBeerStyles()
	if err != nil {
		if !vs.isNoRowsError(err) {
			return fmt.Errorf("failed to check beer styles: %w", err)
		}
		return nil
	}

	for _, style := range beerStyles {
		if style.Name == name {
			return fmt.Errorf("beer style with name '%s' already exists", name)
		}
	}

	return nil
}

func (vs *ValidationService) ValidateUniqueNameForUpdate(name string, excludeUUID string) error {
	if name == "" {
		return nil
	}

	beerStyles, err := vs.beerService.ListAllBeerStyles()
	if err != nil {
		if !vs.isNoRowsError(err) {
			return fmt.Errorf("failed to check beer styles: %w", err)
		}
		return nil
	}

	for _, style := range beerStyles {
		if style.Name == name && style.UUID != excludeUUID {
			return fmt.Errorf("beer style with name '%s' already exists", name)
		}
	}

	return nil
}

func (vs *ValidationService) ValidateUUID(uuidStr string) error {
	if uuidStr == "" {
		return fmt.Errorf("UUID cannot be empty")
	}

	_, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %s", uuidStr)
	}

	return nil
}
