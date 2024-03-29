package models

import (
	"errors"
	"fmt"

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
		UserUUID:    user.UUID,
		Name:        station.Name,
		Description: station.Description,
		Country:     station.Country,
		Subdivision: station.Subdivision,
		Latitude:    station.Latitude,
		Longitude:   station.Longitude,
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

func (c *Controller) DeleteStation(session *tables.User, station *tables.Station) error {
	query := c.DB.
		Where("user_uuid = ?", session.UUID).
		Where("uuid = ?", station.UUID).
		Delete(&tables.Station{})
	if dErr := query.Error; dErr != nil {
		return dErr
	}
	if query.RowsAffected != 1 {
		return ErrUnauthorized
	}
	return nil
}

func (c *Controller) UpdateStation(session *tables.User, station *tables.Station) error {
	query := c.DB.
		Where("user_uuid = ?", session.UUID).
		Where("uuid = ?", station.UUID).
		Limit(1).
		Updates(
			&tables.Station{
				Model:       tables.Model{UUID: station.UUID},
				Name:        station.Name,
				Description: station.Description,
				Country:     station.Country,
				Subdivision: station.Subdivision,
				Latitude:    station.Latitude,
				Longitude:   station.Longitude,
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

func (c *Controller) QueryStationNoAPIKey(station *tables.Station) (*tables.Station, error) {
	s := new(tables.Station)
	qErr := c.DB.
		Preload("Sensors").
		Where("uuid = ?", station.UUID).
		First(s).Error
	return s, qErr
}

func (c *Controller) QueryStationAPIKey(session *tables.User, station *tables.Station) (*tables.Station, error) {
	s := new(tables.Station)
	qErr := c.DB.
		Preload("Sensors").
		Where("uuid = ?", station.UUID).
		Where("user_uuid = ?", session.UUID).
		First(s).Error
	if qErr != nil && errors.Is(qErr, gorm.ErrRecordNotFound) {
		return nil, ErrUnauthorized
	}
	s.APIKeyJSON = &s.APIKey
	return s, qErr
}

func (c *Controller) QueryManyStation(filter *Filter[tables.Station]) (*Results[tables.Station], error) {
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	query := c.DB
	if filter.Target.Name != nil {
		query = query.Where("name LIKE ?", *filter.Target.Name)
	}
	if filter.Target.Description != nil {
		query = query.Where("description LIKE ?", *filter.Target.Description)
	}
	if filter.Target.Country != nil {
		query = query.Where("country = ?", *filter.Target.Country)
	}
	if filter.Target.Subdivision != nil {
		query = query.Where("subdivision = ?", *filter.Target.Subdivision)
	}
	if filter.Target.Latitude != nil {
		query = query.Where("latitude = ?", *filter.Target.Latitude)
	}
	if filter.Target.Longitude != nil {
		query = query.Where("longitude = ?", *filter.Target.Longitude)
	}
	result := &Results[tables.Station]{
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}
	query = query.
		Preload("Sensors").
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

func (c *Controller) Historical(filter *HistoricalFilter) ([]tables.SensorRegistry, error) {
	query := c.DB.
		Where("sensor_uuid = ?", filter.SensorUUID)
	if filter.From != nil {
		query = query.
			Where("? <= created_at", filter.From)
	}
	if filter.To != nil {
		query = query.
			Where("created_at <= ?", filter.To)
	}
	var results []tables.SensorRegistry
	qErr := query.Find(&results).Error
	return results, qErr
}

func (c *Controller) PushRegistry(session *tables.Station, registries []tables.SensorRegistry) (err error) {
	var (
		uuidsMaps       = make(map[uuid.UUID]struct{}, 10)
		uuids           = make([]uuid.UUID, 0, 10)
		uuidRegistries  = make([]*tables.SensorRegistry, 0, 10)
		typesMaps       = make(map[string]struct{}, 10)
		types           = make([]string, 0, 10)
		typesRegistries = make([]*tables.SensorRegistry, 0, 10)
	)
	for _, registry := range registries {
		switch {
		case registry.SensorUUID != uuid.Nil:
			if _, ok := uuidsMaps[registry.SensorUUID]; !ok {
				uuidsMaps[registry.SensorUUID] = struct{}{}
				uuids = append(uuids, registry.SensorUUID)
			}
			uuidRegistries = append(uuidRegistries, &registry)
		case registry.SensorType != nil && *registry.SensorType != "":
			if _, ok := typesMaps[*registry.SensorType]; !ok {
				typesMaps[*registry.SensorType] = struct{}{}
				types = append(types, *registry.SensorType)
			}
			typesRegistries = append(typesRegistries, &registry)
		default:
			err = fmt.Errorf("invalid sensor registry")
			return err
		}
	}

	var (
		rcCount int
		rc      = make([]tables.SensorRegistry, 0, len(registries))
	)

	// By UUID
	if len(uuidsMaps) > 0 {
		var (
			sensorsByUUID    []tables.Sensor
			sensorsByUUIDMap = make(map[uuid.UUID]struct{}, len(uuidsMaps))
		)
		err = c.DB.
			Model(&tables.Sensor{}).
			Where("station_uuid = ?", session.UUID).
			Where(`uuid IN (?)`, uuids).
			Select("uuid").
			Find(&sensorsByUUID).Error
		if err != nil {
			err = fmt.Errorf("sensors found by uuid")
			return err
		}
		for _, sensor := range sensorsByUUID {
			sensorsByUUIDMap[sensor.UUID] = struct{}{}
		}
		for _, registry := range uuidRegistries {
			_, ok := sensorsByUUIDMap[registry.SensorUUID]
			if !ok {
				continue
			}
			rcCount++
			rc = append(rc, tables.SensorRegistry{
				SensorUUID: registry.SensorUUID,
				Value:      registry.Value,
			})
		}
	}

	// By Type
	if len(typesMaps) > 0 {
		var (
			sensorsByType    []tables.Sensor
			sensorsByTypeMap = make(map[string]uuid.UUID, len(uuidsMaps))
		)
		err = c.DB.
			Model(&tables.Sensor{}).
			Where("station_uuid = ?", session.UUID).
			Where(`"type" IN (?)`, types).
			Select("type", "uuid").
			Find(&sensorsByType).Error
		if err != nil {
			err = fmt.Errorf("sensors found by uuid")
			return err
		}
		for _, sensor := range sensorsByType {
			sensorsByTypeMap[sensor.Type] = sensor.UUID
		}
		for _, registry := range typesRegistries {
			sensorUUID, ok := sensorsByTypeMap[*registry.SensorType]
			if !ok {
				continue
			}
			rcCount++
			rc = append(rc, tables.SensorRegistry{
				SensorUUID: sensorUUID,
				Value:      registry.Value,
			})
		}
	}

	if rcCount != len(registries) {
		err = ErrUnauthorized
		return err
	}

	// Create registries
	if len(rc) > 0 {
		err = c.DB.Create(rc).Error
	}
	return err
}
