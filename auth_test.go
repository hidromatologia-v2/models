package models

import (
	"testing"

	"github.com/hidromatologia-v2/models/common/random"
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
		c := sqliteController()
		defer c.Close()
		testAuthenticate(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testAuthenticate(tt, c)
	})
}

func testRegister(t *testing.T, c *Controller) {
	t.Run("Succeed", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.Register(u))
	})
	t.Run("RepeatedUser", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.Register(u))
		assert.NotNil(tt, c.Register(u))
	})
}

func TestRegister(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testRegister(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testRegister(tt, c)
	})
}

func testAuthorize(t *testing.T, c *Controller) {

	t.Run("Succeed", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.Register(u))
		authUser, aErr := c.Authenticate(u)
		assert.Nil(tt, aErr)
		token := c.JWT.New(authUser.Claims())
		jwtUser, jwtErr := c.Authorize(token)
		assert.Nil(tt, jwtErr)
		assert.NotNil(tt, jwtUser)
	})
	t.Run("InvalidUser", func(tt *testing.T) {
		u := tables.RandomUser()
		token := c.JWT.New(u.Claims())
		jwtUser, jwtErr := c.Authorize(token)
		assert.NotNil(tt, jwtErr)
		assert.Nil(tt, jwtUser)
	})
	t.Run("InvalidToken", func(tt *testing.T) {
		jwtUser, jwtErr := c.Authorize("")
		assert.NotNil(tt, jwtErr)
		assert.Nil(tt, jwtUser)
	})
}

func TestAuthorize(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testAuthorize(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testAuthorize(tt, c)
	})
}

func testAuthorizeAPIKey(t *testing.T, c *Controller) {
	t.Run("Succeed", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, c.DB.Create(s).Error)
		station, aErr := c.AuthorizeAPIKey(s.APIKey)
		assert.Nil(tt, aErr)
		assert.NotNil(tt, station)
	})
	t.Run("Fail", func(tt *testing.T) {
		station, aErr := c.AuthorizeAPIKey("INVALID")
		assert.NotNil(tt, aErr)
		assert.Nil(tt, station)
	})
}

func TestAuthorizeAPIKey(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testAuthorizeAPIKey(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testAuthorizeAPIKey(tt, c)
	})
}
