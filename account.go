package models

import (
	"errors"

	"github.com/hidromatologia-v2/models/tables"
	"gorm.io/gorm"
)

func (c *Controller) QueryAccount(user *tables.User) (*tables.User, error) {
	account := new(tables.User)
	err := c.DB.
		Where("uuid = ?", user.UUID).
		First(&account).Error
	if err == nil {
		return account, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUnauthorized
	}
	return nil, err
}
