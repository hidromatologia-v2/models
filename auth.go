package models

import (
	"errors"

	"github.com/hidromatologia-v2/models/tables"
	"gorm.io/gorm"
)

func (c *Controller) Authenticate(u *tables.User) (*tables.User, error) {
	user := new(tables.User)
	err := c.DB.
		Where("username = ?", u.Username).
		First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	if !user.Authenticate(u.Password) {
		return nil, ErrUnauthorized
	}
	return user, nil
}
