package models

import (
	"fmt"
	"time"

	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	uuid "github.com/satori/go.uuid"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

type (
	Filter[T any] struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
		Target   T   `json:"target"`
	}
	Results[T any] struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
		Count    int `json:"count"`
		Entries  []T `json:"entries"`
	}
	HistoricalFilter struct {
		SensorUUID uuid.UUID  `json:"sensorUUID"`
		From       *time.Time `json:"from"`
		To         *time.Time `json:"to"`
	}
)

func pgController() *Controller {
	return NewController(postgres.NewDefault(), cache.RedisDefault(), []byte(random.String()))
}

func sqliteController() *Controller {
	return NewController(sqlite.NewMem(), cache.Bigcache(), []byte(random.String()))
}
