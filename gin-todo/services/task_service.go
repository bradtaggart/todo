package services

import (
	"fmt"

	"github.com/bradtaggart/gin-todo/database"
	"github.com/bradtaggart/gin-todo/models"
)

func GetAllTasks(tasks *[]models.Task) error {
	if err := database.DB.Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

func GetTaskByID(task *models.Task, id string) error {
	if err := database.DB.Where("id = ?", id).First(task).Error; err != nil {
		return err
	}
	return nil
}

func CreateTask(task *models.Task) error {
	fmt.Println(task)
	if err := database.DB.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTask(task *models.Task, id string) error {
	if err := database.DB.Model(task).Where("id = ?", id).Updates(task).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTask(task *models.Task, id string) error {
	if err := database.DB.Where("id = ?", id).Delete(task).Error; err != nil {
		return err
	}
	return nil
}
