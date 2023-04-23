package models

import (
	"testing"

	"github.com/wneessen/go-mail"
)

func testSendMail(t *testing.T, c *Controller) {

}

func TestSendMail(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		c.MailFrom = "sulcud@mail.com"
		c.MailHost = "localhost"
		c.MailOptions = []mail.Option{
			mail.WithPort(1025), mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername("sulcud"), mail.WithPassword("sulcud"),
		}
		defer c.Close()
		testSendMail(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		c.MailFrom = "sulcud@mail.com"
		c.MailHost = "localhost"
		c.MailOptions = []mail.Option{
			mail.WithPort(1025), mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername("sulcud"), mail.WithPassword("sulcud"),
		}
		defer c.Close()
		testSendMail(tt, c)
	})
}
