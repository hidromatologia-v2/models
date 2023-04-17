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
		u2 := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u2).Error)
		assert.Nil(tt, c.DB.Create(alert).Error)
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
