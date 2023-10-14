package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
)

var DB *gorm.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ivanyunus"
	dbname   = "go_test"
)

func ConnectDatabse() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.Exec("SELECT 1").Error
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&Player{},
		&TokenSession{},
		&Bank{},
		&PlayersBank{},
		&Transaction{},
	)

	DB = db

	fmt.Println("Successfully connected!")
}
