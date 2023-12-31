package models

type Player struct {
	Base
	FirstName     string      `gorm:"column:first_name;type(varchar(50))" json:"first_name" `
	LastName      string      `gorm:"column:last_name;type(varchar(50))" json:"last_name" `
	Password      string      `gorm:"column:password;type(varchar(50))" json:"password" `
	Email         string      `gorm:"column:email;type(varchar(100));unique;not null" json:"email" `
	Phone         string      `gorm:"column:phone;type(varchar(20))" json:"phone" `
	Balance       int64       `gorm:"column:balance;type(bigint)" json:"balance" `
	PlayersBankId *string     `gorm:"column:players_bank_id;type(uuid)" json:"players_bank_id"`
	PlayersBank   PlayersBank `json:"players_bank"`
}
