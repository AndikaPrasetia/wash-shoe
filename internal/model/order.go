package model

import "time"

type Order struct {
	ID          int32     `json:"id"`
	UserID      string    `json:"user_id"`
	ServiceType string    `json:"service_type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
