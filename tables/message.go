package tables

import (
	"fmt"
	"net/mail"

	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
)

const (
	Email = "email"
	SMS   = "SMS"
)

type Message struct {
	Model
	Type      string `json:"type" gorm:"not null;"`
	Recipient string `json:"recipient" gorm:"not null;"`
	Subject   string `json:"subject" gorm:"not null;"`
	Message   string `json:"message" gorm:"not null;"`
}

func (m *Message) BeforeSave(tx *gorm.DB) error {
	bErr := m.Model.BeforeSave(tx)
	if bErr != nil {
		return bErr
	}
	switch m.Type {
	case Email:
		_, pErr := mail.ParseAddress(m.Recipient)
		if pErr != nil {
			return pErr
		}
	case SMS:
		if !phoneRegex.MatchString(m.Recipient) {
			return fmt.Errorf("invalid phone")
		}
	default:
		return fmt.Errorf("invalid message type")
	}
	return nil
}

func RandomMessage(t string) *Message {
	switch t {
	case Email:
		return &Message{
			Type:      Email,
			Recipient: gofakeit.NewCrypto().Person().Contact.Email,
			Subject:   gofakeit.NewCrypto().Word(),
			Message:   gofakeit.NewCrypto().Word(),
		}
	case SMS:
		return &Message{
			Type:      SMS,
			Recipient: gofakeit.NewCrypto().Person().Contact.Phone,
			Subject:   gofakeit.NewCrypto().Word(),
			Message:   gofakeit.NewCrypto().Word(),
		}
	default:
		return &Message{
			Type:      Email,
			Recipient: gofakeit.NewCrypto().Person().Contact.Email,
			Subject:   gofakeit.NewCrypto().Word(),
			Message:   gofakeit.NewCrypto().Word(),
		}
	}
}
