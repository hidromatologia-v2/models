package models

import (
	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/session"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/wneessen/go-mail"
	"gorm.io/gorm"
)

type Controller struct {
	DB          *gorm.DB
	Cache       cache.Cache
	JWT         *session.JWT
	MailHost    string
	MailFrom    string
	MailOptions []mail.Option
}

func (c *Controller) Close() error {
	conn, err := c.DB.DB()
	if err != nil {
		return err
	}
	return conn.Close()
}

type (
	MailOptions struct {
		Host    string
		From    string
		Options []mail.Option
	}
	Options struct {
		Database  *gorm.DB
		Cache     cache.Cache
		JWTSecret []byte
		Mail      *MailOptions
	}
)

func NewController(options *Options) *Controller {
	c := &Controller{}
	if options.Database != nil {
		c.DB = options.Database
	}
	if options.Cache != nil {
		c.Cache = options.Cache
	}
	if options.JWTSecret != nil {
		c.JWT = session.New(options.JWTSecret)
	}
	if options.Mail != nil {
		c.MailHost = options.Mail.Host
		c.MailFrom = options.Mail.From
		c.MailOptions = options.Mail.Options
	}
	c.DB.AutoMigrate(
		&tables.User{}, &tables.Station{}, &tables.Sensor{},
		&tables.SensorRegistry{}, &tables.Alert{}, &tables.Admin{},
		&tables.Message{},
	)
	return c
}
