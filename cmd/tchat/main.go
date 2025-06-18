package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vivaswanth-kashyap/tchat/internal/cli"
	"github.com/vivaswanth-kashyap/tchat/internal/db"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	db.InitDB(dbPath)
	cli.Execute()
}
