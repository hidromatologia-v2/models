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
		a := RandomAlert(u, &targetSensor)
		assert.Nil(tt, db.Create(a).Error)
		var alert Alert
		assert.Nil(tt, db.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.Equal(tt, a.Name, alert.Name)
		assert.Equal(tt, a.ConditionOP, alert.ConditionOP)
		assert.Equal(tt, a.Condition, alert.Condition)
		assert.Equal(tt, a.Value, alert.Value)
	})
	t.Run("NoName", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		targetSensor := s.Sensors[0]
		a := RandomAlert(u, &targetSensor)
		a.Name = nil
		assert.NotNil(tt, db.Create(a).Error)
	})
	t.Run("NoConditionOP", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		targetSensor := s.Sensors[0]
		a := RandomAlert(u, &targetSensor)
		a.ConditionOP = nil
		assert.NotNil(tt, db.Create(a).Error)
	})
	t.Run("NoValue", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		targetSensor := s.Sensors[0]
		a := RandomAlert(u, &targetSensor)
		a.Value = nil
		assert.NotNil(tt, db.Create(a).Error)
	})
	t.Run("UpdateCondition", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		targetSensor := s.Sensors[0]
		a := RandomAlert(u, &targetSensor)
		assert.Nil(tt, db.Create(a).Error)
		*a.ConditionOP = ">="
		assert.Nil(tt, db.Where("uuid = ?", a.UUID).Where("user_uuid = ?", u.UUID).Limit(1).Updates(a).Error)
	})
	t.Run("UpdateWithNilCondition", func(tt *testing.T) {
		u := RandomUser()
		assert.Nil(tt, db.Create(u).Error)
		s := RandomStation(u)
		assert.Nil(tt, db.Create(s).Error)
		targetSensor := s.Sensors[0]
		a := RandomAlert(u, &targetSensor)
		assert.Nil(tt, db.Create(a).Error)
		a.ConditionOP = nil
		assert.Nil(tt, db.Where("uuid = ?", a.UUID).Where("user_uuid = ?", u.UUID).Limit(1).Updates(a).Error)
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
