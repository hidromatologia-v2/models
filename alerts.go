package models

import (
	"errors"

	"github.com/hidromatologia-v2/models/tables"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (c *Controller) CreateAlert(session *tables.User, alert *tables.Alert) error {
	return c.DB.Create(&tables.Alert{
		UserUUID:   session.UUID,
		Name:       alert.Name,
		SensorUUID: alert.SensorUUID,
		Condition:  alert.Condition,
		Value:      alert.Value,
		Enabled:    alert.Enabled,
	}).Error
}

func (c *Controller) DeleteAlert(session *tables.User, alert *tables.Alert) error {
	query := c.DB.
		Where("uuid = ?", alert.UUID).
		Where("user_uuid = ?", session.UUID).
		Delete(&tables.Alert{})
	if err := query.Error; err != nil {
		return err
	}
	if query.RowsAffected != 1 {
		return ErrUnauthorized
	}
	return nil
}

func (c *Controller) UpdateAlert(session *tables.User, alert *tables.Alert) error {
	query := c.DB.
		Where("user_uuid = ?", session.UUID).
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
	if query.RowsAffected != 1 {
		return ErrUnauthorized
	}
	return nil
}

func (c *Controller) QueryOneAlert(session *tables.User, alert *tables.Alert) (*tables.Alert, error) {
	result := new(tables.Alert)
	query := c.DB.
		Where("user_uuid = ?", session.UUID).
		Where("uuid = ?", alert.UUID).
		First(result)
	err := query.Error
	if err == nil {
		return result, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUnauthorized
	}
	return nil, err
}

func (c *Controller) QueryManyAlert(session *tables.User, filter *Filter[tables.Alert]) (*Results[tables.Alert], error) {
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	query := c.DB.Where("user_uuid = ?", session.UUID)
	if filter.Target.SensorUUID != uuid.Nil {
		query = query.Where("sensor_uuid = ?", filter.Target.SensorUUID)
	}
	if filter.Target.Name != nil {
		query = query.Where("name LIKE ?", *filter.Target.Name)
	}
	if filter.Target.Condition != nil {
		query = query.Where("condition = ?", *filter.Target.Condition)
	}
	if filter.Target.Value != nil {
		query = query.Where("value = ?", *filter.Target.Value)
	}
	if filter.Target.Enabled != nil {
		query = query.Where("enabled = ?", *filter.Target.Enabled)
	}
	result := &Results[tables.Alert]{
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}
	query = query.
		Order("uuid DESC").
		Offset(filter.PageSize * (filter.Page - 1)).
		Limit(filter.PageSize).
		Find(&result.Entries)
	if err := query.Error; err != nil {
		return nil, err
	}
	result.Count = len(result.Entries)
	return result, nil
}

func (c *Controller) CheckAlert(registry *tables.SensorRegistry) ([]tables.Alert, error) {
	var alerts []tables.Alert
	aErr := c.DB.
		Preload("User").
		Where(`
		(
			enabled = true AND sensor_uuid = ?
		) AND
		(
			(condition = ? AND ? <  value) OR
			(condition = ? AND ? >  value) OR
			(condition = ? AND ? <= value) OR
			(condition = ? AND ? >= value)
		)
		`,
			registry.SensorUUID,
			tables.Lt, registry.Value,
			tables.Gt, registry.Value,
			tables.Le, registry.Value,
			tables.Ge, registry.Value,
		).
		Find(&alerts).
		Update("enabled", false).
		Error
	if aErr != nil {
		return nil, aErr
	}
	return alerts, nil
}
