package types

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/biter777/countries"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	Station struct {
		Model
		Responsible     User                      `json:"responsible" gorm:"foreignKey:ResponsibleUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ResponsibleUUID uuid.UUID                 `json:"responsibleUUID" gorm:"uniqueIndex:idx_unique_station;not null;"`
		Name            string                    `json:"name" gorm:"uniqueIndex:idx_unique_station;not null;"`
		CountryName     string                    `json:"countryName" gorm:"-"`
		SubdivisionName string                    `json:"subdivisionName" gorm:"-"`
		Country         countries.CountryCode     `json:"-" gorm:"not null;"`
		Subdivision     countries.SubdivisionCode `json:"-" gorm:"not null"`
		Description     string                    `json:"description"`
		Latitude        float64                   `json:"latitude" gorm:"not null;"`
		Longitude       float64                   `json:"longitude" gorm:"not null;"`
		Sensors         []Sensor                  `json:"sensors"`
	}
	Sensor struct {
		Model
		Station     Station   `json:"station" gorm:"foreignKey:StationUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		StationUUID uuid.UUID `json:"stationUUID" gorm:"uniqueIndex:idx_unique_sensor;not null;"`
		Type        string    `json:"type" gorm:"uniqueIndex:idx_unique_sensor;not null;"`
	}
)

func (s *Station) AfterFind(tx *gorm.DB) error {
	sort.Slice(s.Sensors, func(i, j int) bool {
		return strings.Compare(s.Sensors[i].Type, s.Sensors[j].Type) < 0
	})
	s.CountryName = countries.ByNumeric(int(s.Country)).String()
	s.SubdivisionName = countries.SubdivisionUnknown.String()
	for _, subdivision := range s.Country.Subdivisions() {
		if subdivision == s.Subdivision {
			s.SubdivisionName = subdivision.String()
		}
	}
	return nil
}

func (s *Station) BeforeSave(tx *gorm.DB) error {
	bErr := s.Model.BeforeSave(tx)
	if bErr != nil {
		return bErr
	}
	sort.Slice(s.Sensors, func(i, j int) bool {
		return strings.Compare(s.Sensors[i].Type, s.Sensors[j].Type) < 0
	})
	s.Country = countries.ByName(s.CountryName)
	if s.Country == countries.Unknown {
		return fmt.Errorf("invalid country provided")
	}
	subdivisionSet := false
	for _, subdivision := range s.Country.Subdivisions() {
		if subdivision.String() == s.SubdivisionName {
			subdivisionSet = true
			s.Subdivision = subdivision.Info().Code
		}
	}
	if !subdivisionSet && s.CountryName == s.SubdivisionName {
		s.Subdivision = countries.SubdivisionUnknown
		subdivisionSet = true
	}
	if !subdivisionSet {
		return fmt.Errorf("invalid subdivision")
	}
	return nil
}

func RandomSensor(station *Station) *Sensor {
	return &Sensor{
		StationUUID: station.UUID,
		Type:        fmt.Sprintf("%s %s %s %s", gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word()),
	}
}

func RandomStation(user *User) *Station {
	latitude, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	longitude, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	station := &Station{
		ResponsibleUUID: user.UUID,
		Name:            fmt.Sprintf("%s %s %s %s", gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word()),
		CountryName:     countries.Colombia.String(),
		SubdivisionName: countries.Colombia.Subdivisions()[0].String(),
		Description:     gofakeit.NewCrypto().JobDescriptor(),
		Latitude:        float64(latitude.Int64()),
		Longitude:       float64(longitude.Int64()),
	}
	for i := 0; i < 5; i++ {
		station.Sensors = append(station.Sensors, *RandomSensor(station))
	}
	return station
}
