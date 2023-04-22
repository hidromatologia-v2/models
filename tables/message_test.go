package tables

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func testMessage(t *testing.T, db *gorm.DB) {
	t.Run("Valid Email", func(tt *testing.T) {
		m := RandomMessage(Email)
		assert.Nil(tt, db.Create(m).Error)
	})
	t.Run("Invalid Email", func(tt *testing.T) {
		m := RandomMessage(Email)
		m.Recipient = "INVALID"
		assert.NotNil(tt, db.Create(m).Error)
	})
	t.Run("Valid SMS", func(tt *testing.T) {
		m := RandomMessage(SMS)
		assert.Nil(tt, db.Create(m).Error)
	})
	t.Run("Invalid SMS", func(tt *testing.T) {
		m := RandomMessage(SMS)
		m.Recipient = "INVALID"
		assert.NotNil(tt, db.Create(m).Error)
	})
	t.Run("Invalid TYPE", func(tt *testing.T) {
		m := RandomMessage("INVALID")
		m.Type = "INVALID"
		m.Recipient = "INVALID"
		assert.NotNil(tt, db.Create(m).Error)
	})
}

func TestMessage(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		db := sqlite.NewMem()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&Message{}))
		testMessage(tt, db)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		db := postgres.NewDefault()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&Message{}))
		testMessage(tt, db)
	})
}
