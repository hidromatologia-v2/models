package tables

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestRandomAdmin(t *testing.T) {
	assert.NotEqual(t, RandomAdmin(), RandomAdmin())
}

func testAdmin(t *testing.T, db *gorm.DB) {
	t.Run("ValidUser", func(tt *testing.T) {
		a := RandomAdmin()
		assert.Nil(tt, db.Create(a).Error)
	})
	t.Run("RepeatedUser", func(tt *testing.T) {
		a := RandomAdmin()
		assert.Nil(tt, db.Create(a).Error)
		assert.NotNil(tt, db.Create(a).Error)
	})
	t.Run("Claims", func(tt *testing.T) {
		a := RandomAdmin()
		assert.Nil(tt, db.Create(a).Error)
		assert.NotNil(tt, a.Claims())
	})
	t.Run("FromClaims-Valid", func(tt *testing.T) {
		a := RandomAdmin()
		var u2 User
		assert.Nil(tt, u2.FromClaims(a.Claims()))
		assert.Equal(tt, a.UUID, u2.UUID)
	})
	t.Run("FromClaims-NoUUID", func(tt *testing.T) {
		a := RandomAdmin()
		var u2 User
		claims := a.Claims()
		delete(claims, "uuid")
		assert.NotNil(tt, u2.FromClaims(claims))
	})
	t.Run("Authenticate-Succeed", func(tt *testing.T) {
		a := RandomAdmin()
		a.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(a.Password), DefaultPasswordCost)
		assert.True(tt, a.Authenticate(a.Password))
	})
	t.Run("Authenticate-Fail", func(tt *testing.T) {
		a := RandomAdmin()
		a.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(a.Password), DefaultPasswordCost)
		assert.False(tt, a.Authenticate("wrong"))
	})
	t.Run("BeforeSave-ValidUser", func(tt *testing.T) {
		a := RandomAdmin()
		assert.Nil(tt, db.Create(a).Error)
	})
	t.Run("BeforeSave-NoPassword", func(tt *testing.T) {
		a := RandomAdmin()
		a.Password = ""
		assert.NotNil(tt, db.Create(a).Error)
	})
}

func TestAdmin(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		db := sqlite.NewMem()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&Admin{}))
		testAdmin(tt, db)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		db := postgres.NewDefault()
		conn, _ := db.DB()
		defer conn.Close()
		assert.Nil(tt, db.AutoMigrate(&Admin{}))
		testAdmin(tt, db)
	})
}
