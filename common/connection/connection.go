package connection

import (
	"fmt"
	"testing"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/memphisdev/memphis.go"
	redis_v9 "github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wneessen/go-mail"
)

func PostgresController() *models.Controller {
	opts := models.Options{
		Database:      postgres.NewDefault(),
		EmailCache:    cache.Redis(&redis_v9.Options{Addr: "127.0.0.1:6379", DB: 1}),
		PasswordCache: cache.Redis(&redis_v9.Options{Addr: "127.0.0.1:6379", DB: 2}),
		JWTSecret:     []byte(random.String()),
		Mail: &models.MailOptions{
			From: "sulcud@mail.com",
			Host: "127.0.0.1",
			Options: []mail.Option{
				mail.WithPort(1025), mail.WithSMTPAuth(mail.SMTPAuthPlain),
				mail.WithUsername(""), mail.WithPassword(""),
				mail.WithTLSPolicy(mail.NoTLS),
			},
		},
	}
	return models.NewController(&opts)
}

const (
	testingStation = "testing"
)

func DefaultConsumer(t *testing.T) *memphis.Consumer {
	return NewConsumer(t, testingStation, random.String())
}

func DefaultProducer(t *testing.T) *memphis.Producer {
	return NewProducer(t, testingStation, random.String())
}

func NewConsumer(t *testing.T, station, name string) *memphis.Consumer {
	conn, cErr := memphis.Connect(
		"127.0.0.1",
		"root",
		memphis.Password("memphis"),
		// memphis.ConnectionToken("memphis"),
	)
	assert.Nil(t, cErr)
	name = fmt.Sprintf("%s-%s", name, uuid.NewV4().String())
	c, err := conn.CreateConsumer(station, name)
	assert.Nil(t, err)
	return c
}

func NewProducer(t *testing.T, station, name string) *memphis.Producer {
	conn, cErr := memphis.Connect(
		"127.0.0.1",
		"root",
		memphis.Password("memphis"),
		// memphis.ConnectionToken("memphis"),
	)
	assert.Nil(t, cErr)
	name = fmt.Sprintf("%s-%s", name, uuid.NewV4().String())
	prod, err := conn.CreateProducer(station, name)
	assert.Nil(t, err)
	return prod
}
