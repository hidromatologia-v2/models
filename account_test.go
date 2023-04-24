package models

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
)

func testQueryAccount(t *testing.T, c *Controller) {

	t.Run("Valid", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		user, err := c.QueryAccount(u)
		assert.Nil(tt, err)
		assert.Equal(tt, user.UUID, u.UUID)
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		user, err := c.QueryAccount(&tables.User{})
		assert.NotNil(tt, err)
		assert.ErrorIs(tt, err, ErrUnauthorized)
		assert.Nil(tt, user)
	})
}

func TestQueryAccount(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testQueryAccount(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testQueryAccount(tt, c)
	})
}

func testUpdateAccount(t *testing.T, c *Controller) {
	t.Run("Basic", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		person := gofakeit.NewCrypto().Person()
		err := c.UpdateAccount(
			u,
			&tables.User{
				Name:     &person.FirstName,
				Phone:    &person.Contact.Phone,
				Email:    &person.Contact.Email,
				Password: random.String()[:72],
			},
			u.Password,
		)
		assert.Nil(tt, err)
		var user tables.User
		assert.Nil(tt, c.DB.Where("uuid = ?", u.UUID).First(&user).Error)
		assert.Equal(tt, user.UUID, u.UUID)
		assert.NotEqual(tt, *u.Name, *user.Name)
		assert.NotEqual(tt, *u.Phone, *user.Phone)
		assert.NotEqual(tt, *u.Email, *user.Email)
		assert.NotEqual(tt, u.PasswordHash, user.PasswordHash)
	})
	t.Run("Unauthorized", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		person := gofakeit.NewCrypto().Person()
		err := c.UpdateAccount(
			u,
			&tables.User{
				Name:     &person.FirstName,
				Phone:    &person.Contact.Phone,
				Email:    &person.Contact.Email,
				Password: random.String(),
			},
			"INVALID",
		)
		assert.NotNil(tt, err)
	})
	t.Run("Unauthorized-InvalidUUID", func(tt *testing.T) {
		person := gofakeit.NewCrypto().Person()
		err := c.UpdateAccount(
			tables.RandomUser(),
			&tables.User{
				Name:     &person.FirstName,
				Phone:    &person.Contact.Phone,
				Email:    &person.Contact.Email,
				Password: random.String(),
			},
			"INVALID",
		)
		assert.NotNil(tt, err)
	})
}

func TestUpdateAccount(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testUpdateAccount(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testUpdateAccount(tt, c)
	})
}

func testRequestConfirmation(t *testing.T, c *Controller) {
	t.Run("Basic", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		emailCode, smsCode, rErr := c.RequestConfirmation(u)
		assert.Nil(tt, rErr)
		var emailUser, smsUser tables.User
		assert.Nil(tt, c.Cache.Get(emailCode, &emailUser))
		assert.Nil(tt, c.Cache.Get(smsCode, &smsUser))
		assert.Equal(tt, u.UUID, emailUser.UUID)
		assert.Equal(tt, u.UUID, smsUser.UUID)
	})
	t.Run("Already Confirmed", func(tt *testing.T) {
		u := tables.RandomUser()
		u.Confirmed = new(bool)
		*u.Confirmed = true
		assert.Nil(tt, c.DB.Create(u).Error)
		_, _, rErr := c.RequestConfirmation(u)
		assert.NotNil(tt, rErr)
	})
}

func TestRequestConfirmation(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testRequestConfirmation(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testRequestConfirmation(tt, c)
	})
}

func testConfirmAccount(t *testing.T, c *Controller) {
	t.Run("Basic", func(tt *testing.T) {
		u := tables.RandomUser()
		assert.Nil(tt, c.DB.Create(u).Error)
		emailCode, smsCode, rErr := c.RequestConfirmation(u)
		assert.Nil(tt, rErr)
		var emailUser, smsUser tables.User
		assert.Nil(tt, c.Cache.Get(emailCode, &emailUser))
		assert.Nil(tt, c.Cache.Get(smsCode, &smsUser))
		assert.Equal(tt, u.UUID, emailUser.UUID)
		assert.Equal(tt, u.UUID, smsUser.UUID)
		// Confirm account
		assert.Nil(tt, c.ConfirmAccount(emailCode, smsCode))
		var user tables.User
		assert.Nil(tt, c.DB.Where("uuid = ?", u.UUID).First(&user).Error)
		assert.True(tt, *user.Confirmed)
	})
}

func TestConfirmAccount(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		defer c.Close()
		testConfirmAccount(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		defer c.Close()
		testConfirmAccount(tt, c)
	})
}
