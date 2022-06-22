package models

import "time"

type User struct {
	ID             string `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name           string `json:"name"`
	Occupation     string `json:"occupation"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	AvatarFileName string `json:"avatar"`
	Role           string `json:"role"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

//userDataFormat
type UserFormatter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email" `
	Token      string `json:"token"`
}

//Register
type Register struct {
	ID         string `json:"id"`
	Name       string `json:"name" `
	Occupation string `json:"occupation" `
	Email      string `json:"email" `
	Password   string `json:"password" `
}
