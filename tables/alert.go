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
	User        User       `json:"user" gorm:"foreignKey:UserUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserUUID    uuid.UUID  `json:"userUUID" gorm:"uniqueIndex:idx_unique_alarm;not null;"`
	Name        *string    `json:"name" gorm:"uniqueIndex:idx_unique_alarm;not null;"`
	Sensor      Sensor     `json:"sensor" gorm:"foreignKey:SensorUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SensorUUID  uuid.UUID  `json:"sensorUUID" gorm:"not null;"`
	ConditionOP *string    `json:"conditionOP" gorm:"-"`
	Condition   *Condition `json:"-" gorm:"not null;"`
	Value       *float64   `json:"value" gorm:"not null;"`
	Enabled     *bool      `json:"enabled" gorm:"not null;default:FALSE"`
}

func (a *Alert) AfterFind(tx *gorm.DB) error {
	op := new(string)
	*op = Conditions[*a.Condition]
	a.ConditionOP = op
	return nil
}

func PickCondition(condition string) (Condition, error) {
	op, found := ConditionsOPs[condition]
	if !found {
		return 0, fmt.Errorf("invalid condition operator")
	}
	return op, nil
}

func (a *Alert) BeforeUpdate(tx *gorm.DB) error {
	if a.ConditionOP == nil {
		a.Condition = nil
	} else {
		op, opErr := PickCondition(*a.ConditionOP)
		if opErr != nil {
			return opErr
		}
		if a.Condition == nil {
			a.Condition = new(Condition)
		}
		*a.Condition = op
	}
	return nil
}

func (a *Alert) BeforeCreate(tx *gorm.DB) error {
	if a.Name == nil {
		return fmt.Errorf("no name provided")
	}
	if a.ConditionOP == nil {
		return fmt.Errorf("no condition provided")
	}
	if a.Value == nil {
		return fmt.Errorf("not value provided")
	}
	op, opErr := PickCondition(*a.ConditionOP)
	if opErr != nil {
		return opErr
	}
	if a.Condition == nil {
		a.Condition = new(Condition)
	}
	*a.Condition = op
	return nil
}

func RandomAlert(user *User, sensor *Sensor) *Alert {
	conditionIndexBigInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ConditionsOPs))))
	conditionIndex := int(conditionIndexBigInt.Int64())
	value, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	var (
		name      = new(string)
		v         = new(float64)
		condition = new(string)
	)
	*name = fmt.Sprintf("%s %s %s %s", gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word(), gofakeit.NewCrypto().Word())
	*v = float64(value.Int64())
	index := 0
	for c := range ConditionsOPs {
		if index != conditionIndex {
			index++
			continue
		}
		*condition = c
		break
	}
	return &Alert{
		UserUUID:    user.UUID,
		Name:        name,
		SensorUUID:  sensor.UUID,
		ConditionOP: condition,
		Value:       v,
	}
}
