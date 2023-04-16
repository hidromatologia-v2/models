package types

import (
	"sort"
	"strings"
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRandomCountry(t *testing.T) {
	assert.NotEqual(t, RandomCountry(), RandomCountry())
}

func testLocation(t *testing.T, db *gorm.DB) {
	t.Run("Query", func(tt *testing.T) {
		c := RandomCountry()
		sort.Slice(c.Regions, func(i, j int) bool {
			return strings.Compare(c.Regions[i].Name, c.Regions[j].Name) < 0
		})
		assert.Nil(tt, db.Create(c).Error)
		var country Country
		assert.Nil(tt, db.Preload("Regions").Where("name = ?", c.Name).First(&country).Error)
		assert.Len(tt, country.Regions, len(c.Regions))
		sort.Slice(country.Regions, func(i, j int) bool {
			return strings.Compare(country.Regions[i].Name, country.Regions[j].Name) < 0
		})
		assert.Equal(tt, c.UUID, country.UUID)
		assert.Equal(tt, c.Name, country.Name)
		for index, region := range c.Regions {
			assert.Equal(tt, region.UUID, country.Regions[index].UUID)
			assert.Equal(tt, region.CountryUUID, country.Regions[index].CountryUUID)
			assert.Equal(tt, region.Name, country.Regions[index].Name)
		}
	})
}

func TestLocation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		db := sqlite.NewMem()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&Country{}, &Region{}))
		testLocation(tt, db)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		db := postgres.NewDefault()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&Country{}, &Region{}))
		testLocation(tt, db)
	})
}
