package main

import (
	"github.com/bradtaggart/gin-todo/controllers"
	"github.com/bradtaggart/gin-todo/database"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	database.Connect()

	router.GET("/tasks", controllers.GetTasks)
	router.GET("/tasks/:id", controllers.GetTask)
	router.POST("/tasks", controllers.CreateTask)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)

	router.Run(":8080")
}
