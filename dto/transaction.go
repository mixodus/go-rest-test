package dto

import "mime/multipart"

type TopUpRequest struct {
	Amount int64                 `form:"amount" binding:"required,number"`
	File   *multipart.FileHeader `form:"file" binding:"required,file"`
}
