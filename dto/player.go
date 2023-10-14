package dto

import "time"

// request
type RegisterRequest struct {
	FirstName string `gorm:"first_name" json:"first_name" binding:"required,min=3,max=50"`
	LastName  string `gorm:"last_name" json:"last_name"`
	Password  string `gorm:"column:password" json:"password" binding:"required"`
	Email     string `gorm:"email" json:"email" binding:"required,email"`
	Phone     string `gorm:"phone" json:"phone" binding:"required,min=9,max=13"`
}

type LoginRequest struct {
	Email    string `gorm:"column:email" json:"email" binding:"required"`
	Password string `gorm:"column:password" json:"password" binding:"required"`
}

// response
type RegisterResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type ProfileResponse struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// userlist response
type Base struct {
	Id        string    `gorm:"column:id;" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type Player struct {
	Base
	FirstName string `gorm:"column:first_name" json:"first_name" binding:"required"`
	LastName  string `gorm:"column:last_name" json:"last_name" binding:"required"`
	Password  string `gorm:"column:password" json:"password" binding:"required"`
	Email     string `gorm:"column:email" json:"email" binding:"required"`
	Phone     string `gorm:"column:phone" json:"phone" binding:"required"`
}
