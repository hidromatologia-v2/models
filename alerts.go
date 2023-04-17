package models

import "github.com/hidromatologia-v2/models/tables"

func (c *Controller) CreateAlert(user *tables.User, alert *tables.Alert) error {
	return c.DB.Create(&tables.Alert{
		UserUUID:    user.UUID,
		Name:        alert.Name,
		SensorUUID:  alert.SensorUUID,
		ConditionOP: alert.ConditionOP,
		Value:       alert.Value,
		Enabled:     alert.Enabled,
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
