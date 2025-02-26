package models

import (
	"time"
)

type User struct {
	ID                uint      `gorm:"primaryKey"`
	Name              string    `json:"name"`
	Email             string    `json:"email" gorm:"unique"`
	Password          string    `json:"password"`
	Salt              string    `json:"salt"`
	Role              string    `json:"role" gorm:"default:'user'"`
	Verified          bool      `json:"verified"`
	VerificationCode  string    `json:"verification_code"` // Баталгаажуулах код
	VerificationToken string    `json:"verification_token"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Address           string    `json:"address"`
	Phone             string    `json:"phone"`
}
