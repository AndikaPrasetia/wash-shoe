package dto

import "time"

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type SignupRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type AuthUser struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"` // Audit
	UpdatedAt time.Time  `json:"updated_at"` // Audit
	DeletedAt *time.Time `json:"-"`          // Soft delete
}
