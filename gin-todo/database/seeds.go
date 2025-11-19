package database

import (
	"github.com/bradtaggart/gin-todo/models"
	"gorm.io/gorm"
)

func SeedTasks(db *gorm.DB) error {
	tasks := []models.Task{
		{Name: "1984", Description: "Dystopian novel about the future.", Priority: 10},
		{Name: "The Great Gatsby", Description: "Narrative on the American dreams.", Priority: 20},
		{Name: "Get milk", Description: "Remember to buy milk from the store.", Priority: 3},
		{Name: "Walk the dog", Description: "Take Fido for a walk in the park.", Priority: 2},
		{Name: "Complete project", Description: "Finish the Go project by the end of the week.", Priority: 1},
		{Name: "Another task", Description: "Finish another task", Priority: 1},
	}

	for _, task := range tasks {
		err := db.FirstOrCreate(&task, models.Task{Name: task.Name}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
