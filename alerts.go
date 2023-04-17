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
