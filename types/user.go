package types

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/mail"
	"regexp"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hidromatologia-v2/models/common/random"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	DefaultPasswordCost    = 12
	DefaultExpirationDelta = 30 * 24 * time.Hour
)

type User struct {
	Model
	Username     string `json:"username" gorm:"unique;not null;"`
	Password     string `json:"password" gorm:"-"`
	PasswordHash []byte `json:"-" gorm:"not null;"`
	Name         string `json:"name" gorm:"not null"`
	Phone        string `json:"phone" gorm:"unique;not null"`
	Email        string `json:"email" gorm:"unique;not null"`
	Confirmed    bool   `json:"confimed" gorm:"not null;default:FALSE;" `
}

func RandomUser() *User {
	number, _ := rand.Int(rand.Reader, big.NewInt(100))
	person := gofakeit.Person()
	return &User{
		Username: fmt.Sprintf("%s%d", person.FirstName, number),
		Email:    person.Contact.Email,
		Phone:    person.Contact.Phone,
		Password: random.String()[:72],
	}
}

func (u *User) Claims() jwt.MapClaims {
	return jwt.MapClaims{
		"uuid": u.UUID.String(),
		"exp":  time.Now().Add(DefaultExpirationDelta).Unix(),
	}
}

func (u *User) FromClaims(m jwt.MapClaims) error {
	userUUID, ok := m["uuid"]
	if !ok {
		return fmt.Errorf("incomplete UUID")
	}
	u.UUID = uuid.FromStringOrNil(userUUID.(string))
	return nil
}

func (u *User) Authenticate(password string) bool {
	return u.PasswordHash != nil && bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)) == nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	bErr := u.Model.BeforeSave(tx)
	if bErr != nil {
		return bErr
	}
	if len(u.Password) == 0 {
		return nil
	}
	var err error
	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), DefaultPasswordCost)
	return err
}

var phoneRegex = regexp.MustCompile(`(?m)^\d{2,18}$`)

func (u *User) BeforeCreate(tx *gorm.DB) error {
	bErr := u.Model.BeforeSave(tx)
	if bErr != nil {
		return bErr
	}
	if len(u.PasswordHash) == 0 {
		return fmt.Errorf("password is empty")
	}
	if !phoneRegex.MatchString(u.Phone) {
		println(u.Phone)
		return fmt.Errorf("invalid phone")
	}
	if _, pErr := mail.ParseAddress(u.Email); pErr != nil {
		return pErr
	}

	return nil
}
