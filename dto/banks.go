package dto

type AddBankRequest struct {
	BankCode      string `json:"bank_code" binding:"required,min=3"`
	AccountNumber string `json:"account_number" binding:"required,number,min=5"`
	AccountName   string `json:"account_name" binding:"required"`
}
