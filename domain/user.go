package domain

import (
	"time"
)

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Name       string    `json:"name"`
	Password   string    `json:"-" gorm:"not null"` // "-" means it won't appear in JSON
	Provider   string    `json:"provider"`          // "google", "email", etc.
	ProviderId string    `json:"provider_id"`       // OAuth provider's user ID
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
