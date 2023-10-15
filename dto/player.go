package dto

import (
	"time"

	"github.com/mixodus/go-rest-test/models"
)

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
	Balance   int64  `json:"balance"`
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
	FirstName     string             `gorm:"column:first_name;type(varchar(50))" json:"first_name" `
	LastName      string             `gorm:"column:last_name;type(varchar(50))" json:"last_name" `
	Password      string             `gorm:"column:password;type(varchar(50))" json:"-" `
	Email         string             `gorm:"column:email;type(varchar(100));unique;not null" json:"email" `
	Phone         string             `gorm:"column:phone;type(varchar(20))" json:"phone" `
	Balance       int64              `gorm:"column:balance;type(bigint)" json:"balance" `
	PlayersBankId *string            `gorm:"column:players_bank_id;type(uuid)" json:"players_bank_id"`
	PlayersBank   models.PlayersBank `json:"players_bank"`
}
