package tables

import (
	"fmt"
	"net/mail"

	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
)

const (
	Email = "email"
)

type Message struct {
	Model
	Type      string `json:"type" gorm:"not null;"`
	Recipient string `json:"recipient" gorm:"not null;"`
	Subject   string `json:"subject" gorm:"not null;"`
	Body      string `json:"body" gorm:"not null;"`
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
			Body:      gofakeit.NewCrypto().Word(),
		}
	default:
		return &Message{
			Type:      Email,
			Recipient: gofakeit.NewCrypto().Person().Contact.Email,
			Subject:   gofakeit.NewCrypto().Word(),
			Body:      gofakeit.NewCrypto().Word(),
		}
	}
}
