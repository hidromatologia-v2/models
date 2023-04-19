package models

import (
	"errors"

	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/tables"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (c *Controller) CreateStation(session *tables.User, station *tables.Station) (*tables.Station, error) {
	var user tables.User
	qErr := c.DB.Where("uuid = ?", session.UUID).First(&user).Error
	if qErr != nil {
		return nil, qErr
	}
	if !*user.Confirmed {
		return nil, ErrUnauthorized
	}
	s := &tables.Station{
		UserUUID:        user.UUID,
		Name:            station.Name,
		Description:     station.Description,
		CountryName:     station.CountryName,
		SubdivisionName: station.SubdivisionName,
		Country:         station.Country,
		Subdivision:     station.Subdivision,
		Latitude:        station.Latitude,
		Longitude:       station.Longitude,
		APIKey:          random.String(),
	}
	err := c.DB.Create(s).Error
	if err != nil {
		panic(err)
	}
	return s, err
}

func (c *Controller) AddSensors(session *tables.User, station *tables.Station, sensors []tables.Sensor) error {
	var s tables.Station
	fErr := c.DB.
		Where("user_uuid = ?", session.UUID).
		Where("uuid = ?", station.UUID).
		First(&s).Error
	if fErr != nil {
		if errors.Is(fErr, gorm.ErrRecordNotFound) {
			return ErrUnauthorized
		}
		return fErr
	}
	ss := make([]tables.Sensor, 0, len(sensors))
	for _, sensor := range sensors {
		ss = append(ss, tables.Sensor{
			StationUUID: s.UUID,
			Type:        sensor.Type,
		})
	}
	err := c.DB.Create(ss).Error
	return err
}

func (c *Controller) DeleteSensors(session *tables.User, station *tables.Station, sensors []tables.Sensor) error {
	var s tables.Station
	fErr := c.DB.
		Where("user_uuid = ?", session.UUID).
		Where("uuid = ?", station.UUID).
		First(&s).Error
	if fErr != nil {
		if errors.Is(fErr, gorm.ErrRecordNotFound) {
			return ErrUnauthorized
		}
		return fErr
	}
	sensorsUUIDs := make([]uuid.UUID, 0, len(sensors))
	for _, sensor := range sensors {
		sensorsUUIDs = append(sensorsUUIDs, sensor.UUID)
	}
	query := c.DB.
		Where("station_uuid = ?", station.UUID).
		Where("uuid IN (?)", sensorsUUIDs).
		Delete(&tables.Sensor{})
	err := query.Error
	if err != nil {
		return err
	}
	if query.RowsAffected == 0 {
		return ErrUnauthorized
	}
	return nil
}
