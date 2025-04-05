package middleware

import "time"

type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expire"`
	Code      int       `json:"code"`
}

type AuthUser struct {
	ID       uint   `json:"ID"`
	Username string `json:"Name"`
}
