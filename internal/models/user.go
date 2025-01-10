package models

import (
	"time"
)

type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email         string    `gorm:"unique;not null;size:255" json:"email"`
	Name          string    `gorm:"size:255" json:"name"`
	FirstName     string    `gorm:"size:100" json:"first_name"`
	LastName      string    `gorm:"size:100" json:"last_name"`
	
	GoogleID      string    `gorm:"unique;size:255" json:"google_id"`
	AccessToken   string    `gorm:"type:text" json:"-"` // Omitted from JSON for security
	RefreshToken  string    `gorm:"type:text" json:"-"` // Omitted from JSON for security
	
	ProfilePicture string    `gorm:"type:text" json:"profile_picture"`
	Locale         string    `gorm:"size:10" json:"locale"`
	
	EmailVerified  bool      `json:"email_verified"`
	
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
