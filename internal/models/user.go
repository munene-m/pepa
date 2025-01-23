package models

import (
	"time"
)

type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email         string    `gorm:"unique;not null;size:255" json:"email"`
	Name          string    `gorm:"size:255" json:"name"`
	
	GoogleID      string    `gorm:"unique;size:255" json:"google_id"`
	
	ProfilePicture string    `gorm:"type:text" json:"profile_picture"`
	Locale         string    `gorm:"size:10" json:"locale"`
	
	EmailVerified  bool      `json:"email_verified"`
	
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
