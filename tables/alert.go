package tables

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Condition int

const (
	Lt Condition = iota
	Gt
	Le
	Ge
)

var ConditionsOPs = map[string]Condition{
	"<":  Lt,
	">":  Gt,
	"<=": Le,
	">=": Ge,
}

type Alert struct {
	Model
	User        User      `json:"user" gorm:"foreignKey:UserUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserUUID    uuid.UUID `json:"userUUID" gorm:"uniqueIndex:idx_unique_alarm;not null;"`
	Name        string    `json:"name" gorm:"uniqueIndex:idx_unique_alarm;not null;"`
	Sensor      Sensor    `json:"sensor" gorm:"foreignKey:SensorUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SensorUUID  uuid.UUID `json:"sensorUUID" gorm:"not null;"`
	ConditionOP string    `json:"conditionOP" gorm:"-"`
	Condition   Condition `json:"-" gorm:"not null;"`
	Value       int       `json:"value" gorm:"not null;"`
}

func (a *Alert) BeforeSave(tx *gorm.DB) error {
	bErr := a.Model.BeforeSave(tx)
	if bErr != nil {
		return bErr
	}
	op, found := ConditionsOPs[a.ConditionOP]
	if !found {
		return fmt.Errorf("invalid condition operator")
	}
	a.Condition = op
	return nil
}
