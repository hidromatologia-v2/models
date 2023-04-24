package models

import (
	"fmt"

	"github.com/hidromatologia-v2/models/tables"
	"github.com/wneessen/go-mail"
)

func (c *Controller) sendEmail(message *tables.Message) error {
	// Register message
	iErr := c.DB.Create(message).Error
	if iErr != nil {
		return iErr
	}
	// -- Send the message
	m := mail.NewMsg()
	if err := m.From(c.MailFrom); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}
	if err := m.To(message.Recipient); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}
	m.Subject(message.Subject)
	m.SetBodyString(mail.TypeTextPlain, message.Body)
	client, err := mail.NewClient(
		c.MailHost,
		c.MailOptions...,
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}
	return nil
}

func (c *Controller) SendMessage(message *tables.Message) error {
	switch message.Type {
	case tables.Email:
		return c.sendEmail(message)
	default:
		return fmt.Errorf("invalid message type")
	}
}
