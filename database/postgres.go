package database

import (
	"console/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/console",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
	)

	log.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connected")
	return db
}

func Init() *gorm.DB {
	db := Connect()
	db.AutoMigrate(&models.Group{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Challenge{})
	db.AutoMigrate(&models.UsersChallenge{})
	db.AutoMigrate(&models.Pool{})

	SetupFromConfig(db)
	SetupAdmin(db)
	return db
}
