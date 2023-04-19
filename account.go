package models

import (
	"errors"

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
	if query.RowsAffected == 0 {
		return ErrUnauthorized
	}
	return nil
}
