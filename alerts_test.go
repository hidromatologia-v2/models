package models

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
)

func testCreateAlert(t *testing.T, c *Controller) {
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		alert := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.CreateAlert(u, alert))
	})
}

func TestCreateAlert(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
		defer c.Close()
		testCreateAlert(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testCreateAlert(tt, NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String())))
	})
}

func testDeleteAlert(t *testing.T, c *Controller) {
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		alert := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(alert).Error)
		assert.Nil(tt, c.DeleteAlert(u, alert))
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		alert := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(alert).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		assert.NotNil(tt, c.DeleteAlert(u2, alert))
	})
}

func TestDeleteAlert(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
		defer c.Close()
		testDeleteAlert(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testDeleteAlert(tt, NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String())))
	})
}

func testUpdateAlert(t *testing.T, c *Controller) {

	t.Run("Name", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		a := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(a).Error)
		originalName := *a.Name
		newName := random.String()
		assert.Nil(tt, c.UpdateAlert(u, &tables.Alert{
			Model: tables.Model{
				UUID: a.UUID,
			},
			Name: &newName,
		}))
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.NotEqual(tt, originalName, *alert.Name)
		assert.Equal(tt, newName, *alert.Name)
	})
	t.Run("Condition", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		a := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(a).Error)
		originalCondition := *a.Condition
		var newCondition string
		for newCondition = tables.Conditions[random.Int(len(tables.Conditions))]; newCondition == originalCondition; newCondition = tables.Conditions[random.Int(len(tables.Conditions))] {
		}
		assert.Nil(tt, c.UpdateAlert(u, &tables.Alert{
			Model: tables.Model{
				UUID: a.UUID,
			},
			Condition: &newCondition,
		}))
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.NotEqual(tt, originalCondition, *alert.Condition)
		assert.Equal(tt, newCondition, *alert.Condition)
	})
	t.Run("Value", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		a := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(a).Error)
		originalValue := *a.Value
		newValue := originalValue + 10
		assert.Nil(tt, c.UpdateAlert(u, &tables.Alert{
			Model: tables.Model{
				UUID: a.UUID,
			},
			Value: &newValue,
		}))
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.NotEqual(tt, originalValue, *alert.Value)
		assert.Equal(tt, newValue, *alert.Value)
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		a := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(a).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		alert := *a
		*alert.Name = random.String()
		assert.NotNil(tt, c.UpdateAlert(u2, &alert))
	})
}

func TestUpdateAlert(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
		defer c.Close()
		testUpdateAlert(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testUpdateAlert(tt, NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String())))
	})
}

func testQueryOneAlert(t *testing.T, c *Controller) {

	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		a := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(a).Error)
		alert, qErr := c.QueryOneAlert(u, a)
		assert.Nil(tt, qErr)
		assert.Equal(tt, *a.Name, *alert.Name)
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		a := tables.RandomAlert(u, &s.Sensors[0])
		assert.Nil(tt, c.DB.Create(a).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		alert, qErr := c.QueryOneAlert(u2, a)
		assert.NotNil(tt, qErr)
		assert.Nil(tt, alert)
	})
}

func TestQueryOneAlert(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
		defer c.Close()
		testQueryOneAlert(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testQueryOneAlert(tt, NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String())))
	})
}

func testQueryManyAlert(t *testing.T, c *Controller) {

	t.Run("NoFilter", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		for i := 0; i < 100; i++ {
			a := tables.RandomAlert(u, &s.Sensors[0])
			assert.Nil(tt, c.DB.Create(a).Error)
		}
		results, qErr := c.QueryManyAlert(u, &Filter[tables.Alert]{PageSize: 100})
		assert.Nil(tt, qErr)
		assert.Equal(tt, 100, results.Count)
	})
	t.Run("SensorUUID", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		for i := 0; i < 100; i++ {
			a := tables.RandomAlert(u, &s.Sensors[0])
			assert.Nil(tt, c.DB.Create(a).Error)
		}
		results, qErr := c.QueryManyAlert(u, &Filter[tables.Alert]{PageSize: 100, Target: tables.Alert{SensorUUID: s.Sensors[0].UUID}})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Name", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		for i := 0; i < 100; i++ {
			a := tables.RandomAlert(u, &s.Sensors[0])
			assert.Nil(tt, c.DB.Create(a).Error)
		}
		name := "%a%"
		results, qErr := c.QueryManyAlert(u, &Filter[tables.Alert]{PageSize: 100, Target: tables.Alert{Name: &name}})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Condition", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		for i := 0; i < 1000; i++ {
			a := tables.RandomAlert(u, &s.Sensors[0])
			assert.Nil(tt, c.DB.Create(a).Error)
		}
		condition := ">"
		results, qErr := c.QueryManyAlert(u, &Filter[tables.Alert]{PageSize: 100, Target: tables.Alert{Condition: &condition}})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
	t.Run("Enabled", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		for i := 0; i < 1000; i++ {
			a := tables.RandomAlert(u, &s.Sensors[0])
			assert.Nil(tt, c.DB.Create(a).Error)
		}
		enabled := false
		results, qErr := c.QueryManyAlert(u, &Filter[tables.Alert]{PageSize: 1000, Target: tables.Alert{Enabled: &enabled}})
		assert.Nil(tt, qErr)
		assert.Greater(tt, results.Count, 0)
	})
}

func TestQueryManyAlert(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
		defer c.Close()
		testQueryManyAlert(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testQueryManyAlert(tt, NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String())))
	})
}

func testCheckAlert(t *testing.T, c *Controller) {

	t.Run("Lt", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		ss := s.Sensors[0]
		a := tables.RandomAlert(u, &ss)
		a.Enabled = &tables.True
		*a.Condition = tables.Lt
		*a.Value = 10
		assert.Nil(tt, c.DB.Create(a).Error)
		users, err := c.CheckAlert(&tables.SensorRegistry{
			SensorUUID: ss.UUID,
			Value:      9,
		})
		assert.Nil(tt, err)
		assert.Len(tt, users, 1)
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.False(tt, *alert.Enabled)

	})
	t.Run("Gt", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		ss := s.Sensors[0]
		a := tables.RandomAlert(u, &ss)
		a.Enabled = &tables.True
		*a.Condition = tables.Gt
		*a.Value = 10
		assert.Nil(tt, c.DB.Create(a).Error)
		users, err := c.CheckAlert(&tables.SensorRegistry{
			SensorUUID: ss.UUID,
			Value:      11,
		})
		assert.Nil(tt, err)
		assert.Len(tt, users, 1)
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.False(tt, *alert.Enabled)
	})
	t.Run("Le", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		ss := s.Sensors[0]
		a := tables.RandomAlert(u, &ss)
		a.Enabled = &tables.True
		*a.Condition = tables.Le
		*a.Value = 10
		assert.Nil(tt, c.DB.Create(a).Error)
		users, err := c.CheckAlert(&tables.SensorRegistry{
			SensorUUID: ss.UUID,
			Value:      10,
		})
		assert.Nil(tt, err)
		assert.Len(tt, users, 1)
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.False(tt, *alert.Enabled)
	})
	t.Run("Ge", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		ss := s.Sensors[0]
		a := tables.RandomAlert(u, &ss)
		a.Enabled = &tables.True
		*a.Condition = tables.Ge
		*a.Value = 10
		assert.Nil(tt, c.DB.Create(a).Error)
		users, err := c.CheckAlert(&tables.SensorRegistry{
			SensorUUID: ss.UUID,
			Value:      10,
		})
		assert.Nil(tt, err)
		assert.Len(tt, users, 1)
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.False(tt, *alert.Enabled)
	})
}

func TestCheckAlert(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
		defer c.Close()
		testCheckAlert(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testCheckAlert(tt, NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String())))
	})
}
