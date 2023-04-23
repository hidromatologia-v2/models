package models

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/vonage/vonage-go-sdk"
)

const (
	SMSCooldown = time.Hour
	SMSFrom     = "Hidromatoligia Service"
)

func (c *Controller) sendSMS(message *tables.Message) error {
	cachingKey := fmt.Sprintf("sms-%s", hex.EncodeToString([]byte(message.Recipient)))
	// Cache recipient
	var exists bool
	gErr := c.Cache.Get(cachingKey, &exists)
	if gErr == nil && exists {
		return nil
	}
	sErr := c.Cache.Set(cachingKey, true, store.WithExpiration(SMSCooldown))
	if sErr != nil {
		return sErr
	}
	// Register message
	iErr := c.DB.Create(message).Error
	if iErr != nil {
		return iErr
	}
	// Prepare sent
	_, _, err := c.SMSClient.Send(
		SMSFrom,
		message.Recipient,
		fmt.Sprintf("%s\n%s", message.Subject, message.Message),
		vonage.SMSOpts{
			// TODO:
			// Callback: fmt.Sprintf("%s%s", c.ServerHost, SMSCallbackRoute),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) sendEmail(message *tables.Message) error {
	panic("Implement me!")
}

func (c *Controller) SendMessage(message *tables.Message) error {
	switch message.Type {
	case tables.SMS:
		return c.sendSMS(message)
	case tables.Email:
		return c.sendEmail(message)
	default:
		return fmt.Errorf("invalid message type")
	}
}
