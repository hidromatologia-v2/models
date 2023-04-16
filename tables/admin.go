package tables

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hidromatologia-v2/models/common/random"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	Model
	Username     string `json:"username" gorm:"unique;not null;"`
	Password     string `json:"password" gorm:"-"`
	PasswordHash []byte `json:"-" gorm:"not null;"`
}

func RandomAdmin() *Admin {
	number, _ := rand.Int(rand.Reader, big.NewInt(100))
	person := gofakeit.NewCrypto().Person()
	return &Admin{
		Username: fmt.Sprintf("%s%d", person.FirstName, number),
		Password: random.String()[:72],
	}
}

func (a *Admin) Claims() jwt.MapClaims {
	return jwt.MapClaims{
		"uuid": a.UUID.String(),
		"exp":  time.Now().Add(DefaultExpirationDelta).Unix(),
	}
}

func (a *Admin) FromClaims(m jwt.MapClaims) error {
	userUUID, ok := m["uuid"]
	if !ok {
		return fmt.Errorf("incomplete UUID")
	}
	a.UUID = uuid.FromStringOrNil(userUUID.(string))
	return nil
}

func (a *Admin) Authenticate(password string) bool {
	return a.PasswordHash != nil && bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(password)) == nil
}

func (a *Admin) BeforeSave(tx *gorm.DB) error {
	bErr := a.Model.BeforeSave(tx)
	if bErr != nil {
		return bErr
	}
	if len(a.Password) == 0 {
		return nil
	}
	var err error
	a.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(a.Password), DefaultPasswordCost)
	return err
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	bErr := a.Model.BeforeSave(tx)
	if bErr != nil {
		return bErr
	}
	if len(a.PasswordHash) == 0 {
		return fmt.Errorf("password is empty")
	}
	return nil
}
