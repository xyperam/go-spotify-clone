package models

import "time"

type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Username         string    `json:"username" gorm:"unique;not null"`
	Password         string    `json:"password" gorm:"not null"`
	Email            string    `json:"email" gorm:"unique;not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsPremium        bool      `json:"is_premium" gorm:"default:false"`
	PremiumExpiresAt time.Time `json:"premium_expires_at"`
}
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
