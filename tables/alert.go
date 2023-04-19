package tables

import (
	"fmt"

	"github.com/hidromatologia-v2/models/common/random"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	Lt = "<"
	Gt = ">"
	Le = "<="
	Ge = ">="
)

var Conditions = []string{
	Lt,
	Gt,
	Le,
	Ge,
}

func CheckCondition(condition string) error {
	switch condition {
	case Lt:
		break
	case Gt:
		break
	case Le:
		break
	case Ge:
		break
	default:
		return fmt.Errorf("invalid alert condition")
	}
	return nil
}

type Alert struct {
	Model
	User       User      `json:"user" gorm:"foreignKey:UserUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserUUID   uuid.UUID `json:"userUUID" gorm:"uniqueIndex:idx_unique_alarm;not null;"`
	Name       *string   `json:"name" gorm:"uniqueIndex:idx_unique_alarm;not null;"`
	Sensor     Sensor    `json:"sensor" gorm:"foreignKey:SensorUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SensorUUID uuid.UUID `json:"sensorUUID" gorm:"not null;"`
	Condition  *string   `json:"condition" gorm:"not null;"`
	Value      *float64  `json:"value" gorm:"not null;"`
	Enabled    *bool     `json:"enabled" gorm:"not null;default:FALSE"`
}

func (a *Alert) BeforeCreate(tx *gorm.DB) error {
	if a.Name == nil {
		return fmt.Errorf("no name provided")
	}
	if a.Condition == nil {
		return fmt.Errorf("no condition provided")
	}
	if a.Value == nil {
		return fmt.Errorf("not value provided")
	}
	if cErr := CheckCondition(*a.Condition); cErr != nil {
		return cErr
	}
	return nil
}

func RandomAlert(user *User, sensor *Sensor) *Alert {
	name := random.Name()
	value := random.Float(1000.0)
	condition := Conditions[random.Int(len(Conditions))]
	return &Alert{
		UserUUID:   user.UUID,
		Name:       &name,
		SensorUUID: sensor.UUID,
		Condition:  &condition,
		Value:      &value,
	}
}
