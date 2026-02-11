package models

import (
	"digital-wallet/pkg/response"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	FullName    string         `json:"full_name" gorm:"uniqueIndex;not null"`
	Email       *string        `json:"email" gorm:"uniqueIndex;null"`
	Password    string         `json:"-" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	PhoneNumber string         `json:"phone_number" gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// HashPassword hashes the user password
func HashAndSalt(pwd []byte) (string, error) {
	if len(pwd) > 72 {
		pwd = pwd[:72]
	}

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", response.Wrap(err, "cannot generate hash")
	}
	return string(hash), nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
