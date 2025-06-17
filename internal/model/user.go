package model

import "time"

type User struct {
	ID          string    `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Provider    string    `json:"provider"`
	ProviderID  string    `json:"provider_id"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
