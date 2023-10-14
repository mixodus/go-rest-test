package models

type PlayersBank struct {
	Base
	Id                string `gorm:"column:id;type(uuid);primaryKey;default:gen_random_uuid()" json:"id"`
	PlayerId          string `gorm:"column:player_id;type(uuid);not null" json:"player_id" `
	BankId            string `gorm:"column:bank_id;type(uuid);not null" json:"bank_id"`
	BankAccountNumber string `gorm:"column:bank_account_number;type(varchar(50))" json:"bank_account_number"`
	BankAccountName   string `gorm:"column:bank_account_name;type(varchar(50))" json:"bank_account_name"`
	Bank              Bank   `gorm:"foreignKey:BankId;references:Id" json:"bank"`
}
