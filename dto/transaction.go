package dto

import "mime/multipart"

type TopUpRequest struct {
	Amount int64                 `form:"amount" binding:"required,number"`
	File   *multipart.FileHeader `form:"file" binding:"required,file"`
}

type SpentRequest struct {
	Amount int64  `form:"amount" binding:"required,number"`
	Notes  string `form:"notes" binding:"required"`
}
