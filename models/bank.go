package models

type Bank struct {
	Base
	Id          string `gorm:"column:id;type(uuid);primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string `gorm:"column:name;type(varchar(50))" json:"name"`
	BankCode    string `gorm:"column:bank_code;type(varchar(5))" json:"bank_code"`
	BankInitial string `gorm:"column:bank_initial;type(varchar(10))" json:"bank_initial"`
}
