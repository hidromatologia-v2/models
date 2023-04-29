package models

import (
	"testing"

	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
	"github.com/wneessen/go-mail"
)

func testSendMail(t *testing.T, c *Controller) {
	t.Run("Valid", func(tt *testing.T) {
		m := tables.RandomMessage(tables.Email)
		assert.Nil(tt, c.SendMessage(m))
		var message tables.Message
		assert.Nil(tt, c.DB.Where("uuid = ?", m.UUID).First(&message).Error)
		assert.Equal(tt, m.UUID, message.UUID)
		assert.Equal(tt, m.Recipient, message.Recipient)
	})
	t.Run("Invalid Type", func(tt *testing.T) {
		m := tables.RandomMessage(tables.Email)
		m.Type = "Invalid"
		assert.NotNil(tt, c.SendMessage(m))
	})
	t.Run("Invalid mail From", func(tt *testing.T) {
		real := c.MailFrom
		c.MailFrom = "INVALID"
		defer func() {
			c.MailFrom = real
		}()
		m := tables.RandomMessage(tables.Email)
		assert.NotNil(tt, c.SendMessage(m))
	})
	t.Run("Invalid mail To", func(tt *testing.T) {
		m := tables.RandomMessage(tables.Email)
		m.Recipient = "INVALID"
		assert.NotNil(tt, c.SendMessage(m))
	})
	t.Run("Invalid Client", func(tt *testing.T) {
		real := c.MailHost
		c.MailHost = "9.9.9.9.9.9.9.9.9.9.9.9.9.9.9.9.9.9.9.9.9.9"
		defer func() {
			c.MailHost = real
		}()
		m := tables.RandomMessage(tables.Email)
		assert.NotNil(tt, c.SendMessage(m))
	})
}

func TestSendMail(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		c.MailFrom = "sulcud@mail.com"
		c.MailHost = "127.0.0.1"
		c.MailOptions = []mail.Option{
			mail.WithPort(1025), mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(""), mail.WithPassword(""),
			mail.WithTLSPolicy(mail.NoTLS),
		}
		defer c.Close()
		testSendMail(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		c.MailFrom = "sulcud@mail.com"
		c.MailHost = "127.0.0.1"
		c.MailOptions = []mail.Option{
			mail.WithPort(1025), mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(""), mail.WithPassword(""),
			mail.WithTLSPolicy(mail.NoTLS),
		}
		defer c.Close()
		testSendMail(tt, c)
	})
}
