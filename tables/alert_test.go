package tables

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func testAlert(t *testing.T, db *gorm.DB) {
	t.Run("Valid", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		targetSensor := s.Sensors[0]
		var r []SensorRegistry
		for i := 0; i < 5; i++ {
			r = append(r, SensorRegistry{
				SensorUUID: targetSensor.UUID,
				Value:      float64(i),
			})
		}
		assert.Nil(tt, db.Create(r).Error)
		a := RandomAlert(u, &targetSensor)
		assert.Nil(tt, db.Create(a).Error)
		var alert Alert
		assert.Nil(tt, db.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.Equal(tt, a.Name, alert.Name)
		assert.Equal(tt, a.ConditionOP, alert.ConditionOP)
		assert.Equal(tt, a.Condition, alert.Condition)
		assert.Equal(tt, a.Value, alert.Value)
	})
}

func TestAlert(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		db := sqlite.NewMem()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&User{}, &Station{}, &Sensor{}, &SensorRegistry{}, &Alert{}))
		testAlert(tt, db)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		db := postgres.NewDefault()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&User{}, &Station{}, &Sensor{}, &SensorRegistry{}, &Alert{}))
		testAlert(tt, db)
	})
}
