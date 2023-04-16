package tables

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/brianvoe/gofakeit/v6"
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

var Conditions = map[Condition]string{
	Lt: "<",
	Gt: ">",
	Le: "<=",
	Ge: ">=",
}

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
	Value       float64   `json:"value" gorm:"not null;"`
}

func (a *Alert) AfterFind(tx *gorm.DB) error {
	a.ConditionOP = Conditions[a.Condition]
	return nil
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

func RandomAlert(user *User, sensor *Sensor) *Alert {
	conditionIndexBigInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ConditionsOPs))))
	conditionIndex := int(conditionIndexBigInt.Int64())
	var condition string
	index := 0
	for c := range ConditionsOPs {
		if index != conditionIndex {
			index++
			continue
		}
		condition = c
		break
	}
	value, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	return &Alert{
		UserUUID:    user.UUID,
		Name:        fmt.Sprintf("%s %s %s %s", gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word()),
		SensorUUID:  sensor.UUID,
		ConditionOP: condition,
		Value:       float64(value.Int64()),
	}
}
