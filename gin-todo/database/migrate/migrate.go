package main

import (
	"github.com/bradtaggart/gin-todo/database"
	"github.com/bradtaggart/gin-todo/models"
)

func init() {
	database.Connect()
}

func main() {
	database.DB.AutoMigrate(&models.Task{})
}
