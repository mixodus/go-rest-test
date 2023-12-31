package models

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
)

var DB *gorm.DB

func ConnectDatabse() {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)
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

	fmt.Println("Successfully connected to " + dbname + "!")
}
