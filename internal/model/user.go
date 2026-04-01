package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	UserStatusNormal  = 1
	UserStatusBanned  = 2
	UserRoleUser      = 1
	UserRoleAdmin     = 2
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"size:32;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:128;uniqueIndex" json:"email"`
	Phone        string         `gorm:"size:16" json:"phone"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Status       int            `gorm:"default:1;not null" json:"status"`
	Role         int            `gorm:"default:1;not null" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *User) IsNormal() bool {
	return u.Status == UserStatusNormal
}

type RefreshToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	TokenHash string    `gorm:"size:255;not null" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
