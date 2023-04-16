package types

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/sqlite"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

/*
TestBaseModel basic unit test to reduce the coverage footprint
*/
func testBaseModel_BeforeSafe(t *testing.T, db *gorm.DB) {
	db.AutoMigrate(&Model{})
	t.Run("Null UUID", func(t *testing.T) {
		assert.Nil(t, db.Create(&Model{}).Error)
	})
	t.Run("Set UUID", func(t *testing.T) {
		assert.Nil(t, db.Create(&Model{UUID: uuid.NewV4()}).Error)
	})
}

func TestBaseModel_BeforeSafe(t *testing.T) {
	t.Run("SQLite", func(t *testing.T) {
		db := sqlite.NewMem()
		conn, _ := db.DB()
		defer conn.Close()
		testBaseModel_BeforeSafe(t, db)
	})
	t.Run("Postgres", func(t *testing.T) {
		db := postgres.NewDefault()
		conn, _ := db.DB()
		defer conn.Close()
		testBaseModel_BeforeSafe(t, db)
	})
}
