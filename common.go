package models

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

type (
	Filter[T any] struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
		Target   T   `json:"target"`
	}
	Results[T any] struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
		Count    int `json:"count"`
		Entries  []T `json:"entries"`
	}
	HistoricalFilter struct {
		SensorUUID uuid.UUID  `json:"sensorUUID"`
		From       *time.Time `json:"from"`
		To         *time.Time `json:"to"`
	}
)
