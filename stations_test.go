package models

import (
	"testing"

	"github.com/biter777/countries"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
)

func testCreateStation(t *testing.T, c *Controller) {
	defer c.Close()
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		u.Confirmed = new(bool)
		*u.Confirmed = true
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		station, err := c.CreateStation(u, s)
		assert.Nil(tt, err)
		assert.Equal(tt, s.Name, station.Name)
	})
	t.Run("Unconfirmed", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		_, err := c.CreateStation(u, s)
		assert.NotNil(tt, err)
	})
}

func TestCreateStation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testCreateStation(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testCreateStation(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}

func testAddSensors(t *testing.T, c *Controller) {
	defer c.Close()
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		s.Sensors = nil
		assert.Nil(tt, c.DB.Create(s).Error)
		var ss []tables.Sensor
		for i := 0; i < 5; i++ {
			ss = append(ss, tables.Sensor{
				StationUUID: s.UUID,
				Type:        random.Name(),
			})
		}
		assert.Nil(tt, c.AddSensors(u, s, ss))
		var sensors []tables.Sensor
		assert.Nil(tt, c.DB.Where("station_uuid = ?", s.UUID).Find(&sensors).Error)
		for index, sensor := range ss {
			assert.Equal(tt, ss[index].Type, sensor.Type)
		}
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		s.Sensors = nil
		assert.Nil(tt, c.DB.Create(s).Error)
		var sensors []tables.Sensor
		for i := 0; i < 5; i++ {
			sensors = append(sensors, tables.Sensor{
				StationUUID: s.UUID,
				Type:        random.Name(),
			})
		}
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		assert.NotNil(tt, c.AddSensors(u2, s, sensors))
	})
}

func TestAddSensors(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testAddSensors(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testAddSensors(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}

func testDeleteSensors(t *testing.T, c *Controller) {
	defer c.Close()
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		s.Sensors = nil
		assert.Nil(tt, c.DB.Create(s).Error)
		var ss []tables.Sensor
		for i := 0; i < 5; i++ {
			ss = append(ss, tables.Sensor{
				StationUUID: s.UUID,
				Type:        random.Name(),
			})
		}
		assert.Nil(tt, c.DB.Create(ss).Error)
		assert.Nil(tt, c.DeleteSensors(u, s, ss))
		var sensors []tables.Sensor
		assert.Nil(tt, c.DB.Where("station_uuid = ?", s.UUID).Find(&sensors).Error)
		assert.Len(tt, sensors, 0)
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		s.Sensors = nil
		assert.Nil(tt, c.DB.Create(s).Error)
		var sensors []tables.Sensor
		for i := 0; i < 5; i++ {
			sensors = append(sensors, tables.Sensor{
				StationUUID: s.UUID,
				Type:        random.Name(),
			})
		}
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		assert.NotNil(tt, c.DeleteSensors(u2, s, sensors))
	})
}

func TestDeleteSensors(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testDeleteSensors(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testDeleteSensors(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}

func testDeleteStation(t *testing.T, c *Controller) {
	defer c.Close()
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		assert.Nil(tt, c.DeleteStation(u, s))
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		assert.NotNil(tt, c.DeleteStation(u2, s))
	})
}

func TestDeleteStation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testDeleteStation(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testDeleteStation(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}

func testUpdateStation(t *testing.T, c *Controller) {
	defer c.Close()
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		// -- Update station
		name := random.Name()
		description := random.Name()
		country := countries.UnitedKingdom
		subdivision := countries.SubdivisionCode("Santander")
		latitude := random.Float(1000.0)
		longitude := random.Float(1000.0)
		assert.Nil(tt, c.UpdateStation(u, &tables.Station{
			Model:       s.Model,
			Name:        &name,
			Description: &description,
			Country:     &country,
			Subdivision: &subdivision,
			Latitude:    &latitude,
			Longitude:   &longitude,
		}))
		//
		var station tables.Station
		assert.Nil(tt, c.DB.Where("uuid = ?", s.UUID).First(&station).Error)
		assert.NotEqual(tt, *s.Name, *station.Name)
		assert.NotEqual(tt, *s.Description, *station.Description)
		assert.NotEqual(tt, *s.Country, *station.Country)
		assert.NotEqual(tt, *s.Subdivision, *station.Subdivision)
		assert.NotEqual(tt, *s.Latitude, *station.Latitude)
		assert.NotEqual(tt, *s.Longitude, *station.Longitude)
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		assert.NotNil(tt, c.UpdateStation(u2, &tables.Station{Model: s.Model}))
	})
}

func TestUpdateStation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testUpdateStation(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testUpdateStation(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}

func testQueryStation(t *testing.T, c *Controller) {
	defer c.Close()
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		station, qErr := c.QueryStation(s)
		assert.Nil(tt, qErr)
		assert.Equal(tt, s.UUID, station.UUID)
	})
}

func TestQueryStation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testQueryStation(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testQueryStation(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}

func testQueryManyStation(t *testing.T, c *Controller) {
	defer c.Close()
	t.Run("NoFilter", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		for i := 0; i < 100; i++ {
			s := tables.RandomStation(u)
			assert.Nil(tt, c.DB.Create(s).Error)
		}
		results, qErr := c.QueryManyStation(&Filter[tables.Station]{PageSize: 100})
		assert.Nil(tt, qErr)
		assert.Equal(tt, 100, results.Count)
	})
	t.Run("Name", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		for i := 0; i < 100; i++ {
			s := tables.RandomStation(u)
			assert.Nil(tt, c.DB.Create(s).Error)
		}
		name := "%a%"
		results, qErr := c.QueryManyStation(&Filter[tables.Station]{
			PageSize: 100,
			Target: tables.Station{
				Name: &name,
			},
		})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Description", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		for i := 0; i < 100; i++ {
			s := tables.RandomStation(u)
			assert.Nil(tt, c.DB.Create(s).Error)
		}
		description := "%a%"
		results, qErr := c.QueryManyStation(&Filter[tables.Station]{
			PageSize: 100,
			Target: tables.Station{
				Description: &description,
			},
		})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Country", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		var country countries.CountryCode
		for i := 0; i < 100; i++ {
			s := tables.RandomStation(u)
			if country == countries.Unknown {
				country = *s.Country
			}
			assert.Nil(tt, c.DB.Create(s).Error)
		}
		results, qErr := c.QueryManyStation(&Filter[tables.Station]{
			PageSize: 100,
			Target: tables.Station{
				Country: &country,
			},
		})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Subdivision", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		var subdivision countries.SubdivisionCode
		for i := 0; i < 100; i++ {
			s := tables.RandomStation(u)
			if subdivision == "" {
				subdivision = *s.Subdivision
			}
			assert.Nil(tt, c.DB.Create(s).Error)
		}
		results, qErr := c.QueryManyStation(&Filter[tables.Station]{
			PageSize: 100,
			Target: tables.Station{
				Subdivision: &subdivision,
			},
		})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Latitude", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		var latitude float64
		for i := 0; i < 100; i++ {
			s := tables.RandomStation(u)
			if latitude == 0.0 {
				latitude = *s.Latitude
			}
			assert.Nil(tt, c.DB.Create(s).Error)
		}
		results, qErr := c.QueryManyStation(&Filter[tables.Station]{
			PageSize: 100,
			Target: tables.Station{
				Latitude: &latitude,
			},
		})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Longitude", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		var longitude float64
		for i := 0; i < 100; i++ {
			s := tables.RandomStation(u)
			if longitude == 0.0 {
				longitude = *s.Longitude
			}
			assert.Nil(tt, c.DB.Create(s).Error)
		}
		results, qErr := c.QueryManyStation(&Filter[tables.Station]{
			PageSize: 100,
			Target: tables.Station{
				Longitude: &longitude,
			},
		})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
}

func TestQueryManyStation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testQueryManyStation(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testQueryManyStation(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}
