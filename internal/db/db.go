package db

import (
	"log"

	"github.com/vivaswanth-kashyap/tchat/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	log.Println("Database connected and migrated successfully")
	return nil
}
