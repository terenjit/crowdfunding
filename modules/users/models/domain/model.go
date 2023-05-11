package models

import (
	"crowdfunding/pkg/token"
	"time"
)

type User struct {
	ID         string    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Occupation string    `json:"occupation"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Avatar     string    `json:"avatar"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// userDataFormat
type UserFormatter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Occupation string `json:"occupation"`
	Email      string `json:"email" `
	Token      string `json:"token"`
}

// Register
type Register struct {
	ID         string `json:"id"`
	Name       string `json:"name" `
	Username   string `json:"username"`
	Occupation string `json:"occupation" `
	Email      string `json:"email" `
	Password   string `json:"password" `
}

// LoginRequest
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersGetList struct {
	Status   string      `json:"status"`
	Search   interface{} `json:"search"`
	Page     int         `json:"page"`
	Quantity int         `json:"quantity"`
	Query    string      `query:"query"`
	SortBy   string      `json:"sort_by"`
	UserID   string      `json:"user_id"`
	Opts     token.Claim `json:"opts"`
}

type UpdatedUser struct {
	ID         string    `json:"id" param:"id" validate:"required"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Occupation string    `json:"occupation"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Avatar     string    `json:"avatar"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"update_by"`
}
