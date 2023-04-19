package tables

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/biter777/countries"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/hidromatologia-v2/models/common/random"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	Station struct {
		Model
		User            User                      `json:"user" gorm:"foreignKey:UserUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		UserUUID        uuid.UUID                 `json:"userUUID" gorm:"uniqueIndex:idx_unique_station;not null;"`
		Name            string                    `json:"name" gorm:"uniqueIndex:idx_unique_station;not null;"`
		Description     string                    `json:"description"`
		CountryName     string                    `json:"countryName" gorm:"-"`
		SubdivisionName string                    `json:"subdivisionName" gorm:"-"`
		Country         countries.CountryCode     `json:"-" gorm:"not null;"`
		Subdivision     countries.SubdivisionCode `json:"-" gorm:"not null"`
		Latitude        float64                   `json:"latitude" gorm:"not null;"`
		Longitude       float64                   `json:"longitude" gorm:"not null;"`
		APIKey          string                    `json:"apiKey" gorm:"not null;"`
		Sensors         []Sensor                  `json:"sensors" gorm:"constraint:OnDelete:CASCADE;"`
	}
	Sensor struct {
		Model
		Station     Station          `json:"station" gorm:"foreignKey:StationUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		StationUUID uuid.UUID        `json:"stationUUID" gorm:"uniqueIndex:idx_unique_sensor;not null;"`
		Type        string           `json:"type" gorm:"uniqueIndex:idx_unique_sensor;not null;"`
		Registries  []SensorRegistry `json:"registries" gorm:"constraint:OnDelete:CASCADE;"`
	}
	SensorRegistry struct {
		Model
		Sensor     Sensor    `json:"sensor" gorm:"foreignKey:SensorUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		SensorUUID uuid.UUID `json:"sensorUUID" gorm:"not null;"`
		Value      float64   `json:"value" gorm:"constraint:OnDelete:CASCADE;"`
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
	s.APIKey = random.String()
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
		Type:        random.Name(),
	}
}

func RandomStation(user *User) *Station {
	latitude, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	longitude, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	station := &Station{
		UserUUID:        user.UUID,
		Name:            random.Name(),
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
