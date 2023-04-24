//go:build vonage

package models

import (
	"os"
	"testing"

	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
	"github.com/vonage/vonage-go-sdk"
)

var (
	vonageAPIKey, vonageSecret, smsRecipient string
)

func init() {
	vonageAPIKey = os.Getenv("VONAGE_API_KEY")
	vonageSecret = os.Getenv("VONAGE_SECRET")
	smsRecipient = os.Getenv("SMS_RECIPIENT")
}

func testSendSMS(t *testing.T, c *Controller) {
	t.Run("Valid", func(tt *testing.T) {
		m := tables.RandomMessage(tables.SMS)
		m.Recipient = smsRecipient
		assert.Nil(tt, c.SendMessage(m))
		var message tables.Message
		assert.Nil(tt, c.DB.Where("uuid = ?", m.UUID).First(&message).Error)
		assert.Equal(tt, m.UUID, message.UUID)
	})
}

func TestSendSMS(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		c := sqliteController()
		c.SMSClient = vonage.NewSMSClient(
			vonage.CreateAuthFromKeySecret(vonageAPIKey, vonageSecret),
		)
		defer c.Close()
		testSendSMS(tt, c)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		c := pgController()
		c.SMSClient = vonage.NewSMSClient(
			vonage.CreateAuthFromKeySecret(vonageAPIKey, vonageSecret),
		)
		defer c.Close()
		testSendSMS(tt, c)
	})
}
