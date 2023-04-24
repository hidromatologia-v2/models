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
