package models

import (
	"testing"

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
		testCreateAlert(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testCreateAlert(tt, NewController(postgres.NewDefault(), []byte(random.String())))
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
		testDeleteAlert(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testDeleteAlert(tt, NewController(postgres.NewDefault(), []byte(random.String())))
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
		originalCondition := *a.ConditionOP
		var newCondition string
		for c := range tables.ConditionsOPs {
			if c == originalCondition {
				continue
			}
			newCondition = c
			break
		}
		assert.Nil(tt, c.UpdateAlert(u, &tables.Alert{
			Model: tables.Model{
				UUID: a.UUID,
			},
			ConditionOP: &newCondition,
		}))
		var alert tables.Alert
		assert.Nil(tt, c.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.NotEqual(tt, originalCondition, *alert.ConditionOP)
		assert.Equal(tt, newCondition, *alert.ConditionOP)
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
		testUpdateAlert(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testUpdateAlert(tt, NewController(postgres.NewDefault(), []byte(random.String())))
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
		testQueryOneAlert(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testQueryOneAlert(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}
