package models

import (
	"errors"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/tables"
	"gorm.io/gorm"
)

func (c *Controller) QueryAccount(session *tables.User) (*tables.User, error) {
	account := new(tables.User)
	err := c.DB.
		Where("uuid = ?", session.UUID).
		First(&account).Error
	if err == nil {
		return account, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUnauthorized
	}
	return nil, err
}

/*
Update account user should have the Password set, the details to update
*/
func (c *Controller) UpdateAccount(session, details *tables.User, password string) error {
	var account tables.User
	fErr := c.DB.Where("uuid = ?", session.UUID).First(&account).Error
	if fErr != nil {
		if errors.Is(fErr, gorm.ErrRecordNotFound) {
			return ErrUnauthorized
		}
		return fErr
	}
	if !account.Authenticate(password) {
		return ErrUnauthorized
	}
	query := c.DB.
		Where("uuid = ?", session.UUID).
		Limit(1).
		Updates(
			&tables.User{
				Model:    tables.Model{UUID: session.UUID},
				Name:     details.Name,
				Phone:    details.Phone,
				Email:    details.Email,
				Password: details.Password,
			},
		)
	if err := query.Error; err != nil {
		return err
	}
	if query.RowsAffected != 1 {
		return ErrUnauthorized
	}
	return nil
}

func (c *Controller) RequestConfirmation(session *tables.User) (string, string, error) {
	var user tables.User
	qErr := c.DB.Where("uuid = ?", session.UUID).Where("NOT confirmed").First(&user).Error
	if qErr != nil {
		if errors.Is(qErr, gorm.ErrRecordNotFound) {
			return "", "", ErrUnauthorized
		}
		return "", "", qErr
	}
	emailCode := random.String()[:5]
	smsCode := random.String()[:5]
	sErr := c.Cache.Set(emailCode, user, store.WithExpiration(time.Hour))
	if sErr != nil {
		return "", "", sErr
	}
	sErr = c.Cache.Set(smsCode, user, store.WithExpiration(time.Hour))
	if sErr != nil {
		c.Cache.Del(emailCode)
		return "", "", sErr
	}
	return emailCode, smsCode, nil
}

func (c *Controller) ConfirmAccount(emailCode, smsCode string) error {
	var (
		emailUser, smsUser tables.User
	)
	eErr := c.Cache.Get(emailCode, &emailUser)
	if eErr != nil {
		return eErr
	}
	eErr = c.Cache.Get(smsCode, &smsUser)
	if eErr != nil {
		return eErr
	}
	if emailUser.UUID != smsUser.UUID {
		return ErrUnauthorized
	}
	go c.Cache.Del(emailCode)
	go c.Cache.Del(smsCode)
	var user tables.User
	qErr := c.DB.Where("uuid = ?", emailUser.UUID).First(&user).Error
	if qErr != nil {
		if errors.Is(qErr, gorm.ErrRecordNotFound) {
			return ErrUnauthorized
		}
		return qErr
	}
	query := c.DB.
		Where("uuid = ?", user.UUID).
		Limit(1).
		Updates(&tables.User{
			Model:     user.Model,
			Confirmed: &tables.True,
		})
	if err := query.Error; err != nil {
		return err
	}
	if query.RowsAffected != 1 {
		return ErrUnauthorized
	}
	return nil
}
