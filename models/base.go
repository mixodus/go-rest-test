package models

import "time"

type Base struct {
	Id        string     `gorm:"column:id;type(uuid);primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"-"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}
