package tables

import (
	"sort"
	"strings"
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRandomStation(t *testing.T) {
	u := RandomUser()
	assert.NotEqual(t, RandomStation(u), RandomStation(u))
}

func TestRandomSensor(t *testing.T) {
	u := RandomUser()
	s := RandomStation(u)
	assert.NotEqual(t, RandomSensor(s), RandomSensor(s))
}

func testStation(t *testing.T, db *gorm.DB) {
	t.Run("Valid", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		// Check sensors
		var sensors []Sensor
		assert.Nil(tt, db.Where("station_uuid = ?", s.UUID).Find(&sensors).Error)
		sort.Slice(sensors, func(i, j int) bool {
			return strings.Compare(s.Sensors[i].Type, s.Sensors[j].Type) < 0
		})
		assert.Len(tt, sensors, len(s.Sensors))
		for index, sensor := range s.Sensors {
			assert.Equal(tt, sensor.StationUUID, sensors[index].StationUUID)
			assert.Equal(tt, sensor.UUID, sensors[index].UUID)
			assert.Equal(tt, sensor.Type, sensors[index].Type)
		}
		// Check preloading
		var station Station
		assert.Nil(tt, db.Preload("Sensors").Where("uuid = ?", s.UUID).First(&station).Error)
		assert.Equal(tt, s.UUID, station.UUID)
		assert.Equal(tt, s.UserUUID, station.UserUUID)
		assert.Len(tt, station.Sensors, len(s.Sensors))
		for index, sensor := range s.Sensors {
			assert.Equal(tt, sensor.UUID, station.Sensors[index].UUID)
			assert.Equal(tt, sensor.Type, station.Sensors[index].Type)
		}
	})
	t.Run("InvalidCountry", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		s.CountryName = "WATEMALA"
		assert.NotNil(tt, db.Create(s).Error)
	})
	t.Run("InvalidSubdivision", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		s.SubdivisionName = "WATEMALA"
		assert.NotNil(tt, db.Create(s).Error)
	})
	t.Run("NoSubdivision", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		s.SubdivisionName = s.CountryName
		assert.Nil(tt, db.Create(s).Error)
	})
	t.Run("Registries", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		var r []SensorRegistry
		for i := 0; i < 5; i++ {
			r = append(r, SensorRegistry{
				SensorUUID: s.Sensors[0].UUID,
				Value:      float64(i),
			})
		}
		assert.Nil(tt, db.Create(r).Error)
		var registries []SensorRegistry
		assert.Nil(tt, db.Where("sensor_uuid = ?", s.Sensors[0].UUID).Find(&registries).Error)
		assert.Len(tt, registries, len(r))
		for index, registry := range r {
			assert.Equal(tt, registry.SensorUUID, registries[index].SensorUUID)
			assert.Equal(tt, registry.Value, registries[index].Value)
		}
	})
}

func TestStation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		db := sqlite.NewMem()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&User{}, &Station{}, &Sensor{}, &SensorRegistry{}))
		testStation(tt, db)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		db := postgres.NewDefault()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&User{}, &Station{}, &Sensor{}, &SensorRegistry{}))
		testStation(tt, db)
	})
}
