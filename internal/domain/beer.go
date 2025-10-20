package domain

import "time"

type BeerStyle struct {
	UUID      string    `json:"uuid" ksql:"uuid"`
	Name      string    `json:"name" binding:"required" ksql:"name"`
	TempMin   float64   `json:"temp_min" ksql:"temp_min"`
	TempMax   float64   `json:"temp_max" ksql:"temp_max"`
	CreatedAt time.Time `json:"created_at" ksql:"created_at"`
	UpdatedAt time.Time `json:"updated_at" ksql:"updated_at"`
}

type BeerStyleUpdateRequest struct {
	Name    *string  `json:"name,omitempty"`
	TempMin *float64 `json:"temp_min,omitempty"`
	TempMax *float64 `json:"temp_max,omitempty"`
}
