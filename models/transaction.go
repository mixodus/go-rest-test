package models

type Status string
type TransactionType string

const (
	PENDING Status = "pending"
	FAILED  Status = "failed"
	SUCCESS Status = "success"
)

const (
	DEBIT  TransactionType = "debit"
	CREDIT TransactionType = "credit"
)

type Transaction struct {
	Base
	PlayerId        string          `gorm:"column:player_id;type:uuid;not null" json:"player_id" `
	PlayersBankId   string          `gorm:"column:players_bank_id;type:uuid;not null" json:"players_bank_id"`
	Amount          int64           `gorm:"column:amount;type:bigint" json:"amount"`
	Status          Status          `gorm:"column:status;default:'pending'" json:"status"`
	TransactionType TransactionType `gorm:"column:transaction_type;not null" json:"transaction_type"`
	FileName        string          `gorm:"column:file_name;type:varchar(255)" json:"file_name"`
	Notes           string          `gorm:"column:notes;type:varchar(255)" json:"notes"`
}
