package models

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
)

func testAuthenticate(t *testing.T, c *Controller) {
	t.Run("Succeed", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		authUser, aErr := c.Authenticate(u)
		assert.Nil(tt, aErr)
		assert.Equal(tt, u.UUID, authUser.UUID)
	})
	t.Run("InvalidUsername", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		authUser, aErr := c.Authenticate(&tables.User{
			Username: random.String(),
			Password: random.String(),
		})
		assert.NotNil(tt, aErr)
		assert.Nil(tt, authUser)
	})
	t.Run("InvalidPassword", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		authUser, aErr := c.Authenticate(&tables.User{
			Username: u.Username,
			Password: random.String(),
		})
		assert.NotNil(tt, aErr)
		assert.Nil(tt, authUser)
	})
}

func TestAuthenticate(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		testAuthenticate(tt, NewController(sqlite.NewMem(), []byte(random.String())))
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		testAuthenticate(tt, NewController(postgres.NewDefault(), []byte(random.String())))
	})
}
