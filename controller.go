package models

import (
	"github.com/hidromatologia-v2/models/session"
	"github.com/hidromatologia-v2/models/tables"
	"gorm.io/gorm"
)

type Controller struct {
	DB  *gorm.DB
	JWT *session.JWT
}

func (c *Controller) Close() error {
	conn, err := c.DB.DB()
	if err != nil {
		return err
	}
	return conn.Close()
}

func NewController(db *gorm.DB, secret []byte) *Controller {
	c := &Controller{
		DB:  db,
		JWT: session.New(secret),
	}
	c.DB.AutoMigrate(
		&tables.User{}, &tables.Station{}, &tables.Sensor{},
		&tables.SensorRegistry{}, &tables.Alert{},
	)
	return c
}
