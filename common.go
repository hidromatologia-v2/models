package models

import "fmt"

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
)
