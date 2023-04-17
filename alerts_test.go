package models

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
)

func testCreateAlert(t *testing.T, c *Controller) {
	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(t, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(t, c.DB.Create(s).Error)
		alert := &tables.Alert{
			Name:        fmt.Sprintf("%s %s %s %s", gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word()),
			SensorUUID:  s.Sensors[0].UUID,
			ConditionOP: ">",
			Value:       10,
		}
		assert.Nil(t, c.CreateAlert(u, alert))
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
