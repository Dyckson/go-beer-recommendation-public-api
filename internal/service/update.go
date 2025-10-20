package service

import (
	"backend-test/internal/domain"
)

type UpdateService struct{}

func NewUpdateService() *UpdateService {
	return &UpdateService{}
}

func (us *UpdateService) ApplyBeerStyleUpdates(current *domain.BeerStyle, updates domain.BeerStyleUpdateRequest) bool {
	changed := false

	if updates.Name != nil && *updates.Name != "" && *updates.Name != current.Name {
		current.Name = *updates.Name
		changed = true
	}

	if updates.TempMin != nil && *updates.TempMin != current.TempMin {
		current.TempMin = *updates.TempMin
		changed = true
	}

	if updates.TempMax != nil && *updates.TempMax != current.TempMax {
		current.TempMax = *updates.TempMax
		changed = true
	}

	return changed
}

func (us *UpdateService) GetChangedFields(original domain.BeerStyle, updates domain.BeerStyleUpdateRequest) []string {
	var changedFields []string

	if updates.Name != nil && *updates.Name != "" && *updates.Name != original.Name {
		changedFields = append(changedFields, "Name")
	}
	if updates.TempMin != nil && *updates.TempMin != original.TempMin {
		changedFields = append(changedFields, "TempMin")
	}
	if updates.TempMax != nil && *updates.TempMax != original.TempMax {
		changedFields = append(changedFields, "TempMax")
	}

	return changedFields
}
