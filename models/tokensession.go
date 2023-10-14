package models

// !!UNUSED since we use redis for session management!!
type TokenSession struct {
	Base
	Id       string `gorm:"column:id;type(uuid);primaryKey;default:gen_random_uuid()" json:"id"`
	PlayerId string `gorm:"column:player_id;type(uuid);not null" json:"player_id" `
	Token    string `gorm:"column:token;type(varchar(100));unique;not null" json:"token" `
}
