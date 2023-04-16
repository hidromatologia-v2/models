package models

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
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

func (c *Controller) Register(u *tables.User) error {
	return c.DB.Create(u).Error
}

func (c *Controller) Authorize(jwtS string) (*tables.User, error) {
	token, pErr := jwt.Parse(jwtS, c.JWT.KeyFunc)
	if pErr != nil {
		return nil, pErr
	}
	var fUser tables.User
	fErr := fUser.FromClaims(token.Claims.(jwt.MapClaims))
	if fErr != nil {
		return nil, fErr
	}
	user := new(tables.User)
	err := c.DB.
		Where("uuid = ?", fUser.UUID).
		First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	return user, nil
}
