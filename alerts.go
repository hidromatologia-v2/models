package models

import (
	"github.com/hidromatologia-v2/models/tables"
)

func (c *Controller) CreateAlert(user *tables.User, alert *tables.Alert) error {
	return c.DB.Create(&tables.Alert{
		UserUUID:   user.UUID,
		Name:       alert.Name,
		SensorUUID: alert.SensorUUID,
		Condition:  alert.Condition,
		Value:      alert.Value,
		Enabled:    alert.Enabled,
	}).Error
}

func (c *Controller) DeleteAlert(user *tables.User, alert *tables.Alert) error {
	query := c.DB.
		Where("uuid = ?", alert.UUID).
		Where("user_uuid = ?", user.UUID).
		Delete(&tables.Alert{})
	if err := query.Error; err != nil {
		return err
	}
	if query.RowsAffected == 0 {
		return ErrUnauthorized
	}
	return nil
}

func (c *Controller) UpdateAlert(user *tables.User, alert *tables.Alert) error {
	query := c.DB.
		Where("user_uuid = ?", user.UUID).
		Where("uuid = ?", alert.UUID).
		Limit(1).
		Updates(
			&tables.Alert{
				Model:     tables.Model{UUID: alert.UUID},
				Name:      alert.Name,
				Condition: alert.Condition,
				Value:     alert.Value,
				Enabled:   alert.Enabled,
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

func (c *Controller) QueryOneAlert(user *tables.User, alert *tables.Alert) (*tables.Alert, error) {
	result := new(tables.Alert)
	query := c.DB.
		Where("user_uuid = ?", user.UUID).
		Where("uuid = ?", alert.UUID).
		First(result)
	if err := query.Error; err != nil {
		return nil, err
	}
	if query.RowsAffected == 0 {
		return nil, ErrUnauthorized
	}
	return result, nil
}

// func (c *Controller) QueryManyAlert(user *tables.User, filter *Filter[tables.Alert]) (*Results[tables.Alert], error) {
// 	query := c.DB.
// 		Where("user_uuid = ?", user.UUID).
// 		Where("uuid = ?", filter.Target.UUID)
// 	if filter.Target.Name != nil {
// 		query = query.Where("name LIKE ?", *filter.Target.Name)
// 	}
// 	if filter.Target.Value != nil {
// 		query = query.Where("name LIKE ?", *filter.Target.Name)
// 	}
// 	if err := query.Error; err != nil {
// 		return nil, err
// 	}
// 	if query.RowsAffected == 0 {
// 		return nil, ErrUnauthorized
// 	}
// 	return result, nil
// }
