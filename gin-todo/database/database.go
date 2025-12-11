package database

import (
	"fmt"
	"log"

	//"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	// Sqlite driver based on CGO
	//"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bradtaggart/gin-todo/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	//lsdsn := "host=localhost user=brad password=3582 dbname=todo_list port=5432 sslmode=disable TimeZone=America/Denver"
	DB, err = gorm.Open(sqlite.Open("tasks"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connection established")

	DB.AutoMigrate(&models.Task{})

	// Seed the database
	if err := SeedTasks(DB); err != nil {
		log.Fatal("Failed to seed tasks:", err)
	}
}
