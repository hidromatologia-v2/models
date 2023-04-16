package types

import (
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
)

type (
	Country struct {
		Model
		Name    string   `json:"name" gorm:"unique;not null"`
		Regions []Region `json:"regions"`
	}
	Region struct {
		Model
		Country     Country   `json:"country,omitempty" gorm:"foreignKey:CountryUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CountryUUID uuid.UUID `json:"countryUUID" gorm:"uniqueIndex:idx_unique_region;not null;"`
		Name        string    `json:"name" gorm:"uniqueIndex:idx_unique_region;not null;"`
	}
)

func RandomCountry() *Country {
	return &Country{
		Name: gofakeit.Country(),
		Regions: []Region{
			{
				Name: gofakeit.State(),
			},
			{
				Name: gofakeit.State(),
			},
			{
				Name: gofakeit.State(),
			},
			{
				Name: gofakeit.State(),
			},
		},
	}
}
