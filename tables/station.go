package tables

import (
	"fmt"
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
		User        User                       `json:"user" gorm:"foreignKey:UserUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		UserUUID    uuid.UUID                  `json:"userUUID" gorm:"uniqueIndex:idx_unique_station;not null;"`
		Name        *string                    `json:"name" gorm:"uniqueIndex:idx_unique_station;not null;"`
		Description *string                    `json:"description"`
		Country     *countries.CountryCode     `json:"country" gorm:"not null;"`
		Subdivision *countries.SubdivisionCode `json:"subdivision" gorm:"not null"`
		Latitude    *float64                   `json:"latitude" gorm:"not null;"`
		Longitude   *float64                   `json:"longitude" gorm:"not null;"`
		APIKey      string                     `json:"-" gorm:"not null;"`
		APIKeyJSON  *string                    `json:"apiKey" gorm:"-"`
		Sensors     []Sensor                   `json:"sensors" gorm:"constraint:OnDelete:CASCADE;"`
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
		SensorType *string   `json:"sensorType,omitempty" gorm:"-"`
		Value      float64   `json:"value"`
	}
)

func (s *Station) BeforeCreate(tx *gorm.DB) error {
	s.APIKey = random.String()
	sort.Slice(s.Sensors, func(i, j int) bool {
		return strings.Compare(s.Sensors[i].Type, s.Sensors[j].Type) < 0
	})
	if s.Country == nil {
		return fmt.Errorf("no country provided")
	}
	if *s.Country == countries.Unknown {
		return fmt.Errorf("invalid country provided")
	}
	subdivisionSet := false
	for _, subdivision := range s.Country.Subdivisions() {
		if s.Subdivision.String() != subdivision.String() {
			continue
		}
		subdivisionSet = true
		break
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
	name := random.Name()
	country := countries.Colombia
	subdivision := countries.Colombia.Subdivisions()[0]
	description := gofakeit.NewCrypto().JobDescriptor()
	latitude := random.Float(10000.0)
	longitude := random.Float(10000.0)
	station := &Station{
		UserUUID:    user.UUID,
		Name:        &name,
		Country:     &country,
		Subdivision: &subdivision,
		Description: &description,
		Latitude:    &latitude,
		Longitude:   &longitude,
	}
	for i := 0; i < 5; i++ {
		station.Sensors = append(station.Sensors, *RandomSensor(station))
	}
	return station
}
