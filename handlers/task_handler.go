package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"gotask/db"
	"gotask/models"

	"github.com/gin-gonic/gin"
)

var (
	tasks  []models.Task
	nextID uint = 1
	mu     sync.Mutex
)



func GetTasks(c *gin.Context) {
	var tasks  []models.Task
	query := db.DB.Model(&models.Task{})

	if status := c.Query("status"); status != ""{
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})

}



func CreateTask(c *gin.Context) {
	var input models.CreateTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	newTask := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
	}

	if err := db.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": newTask})
}



func GetTask(c *gin.Context) {
	var task  []models.Task

	
	if err := db.DB.First(&task, c.Param("id")).Error;  err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"data": task})

}


func UpdateTask(c *gin.Context) {
	var task  []models.Task

	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&task).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": task})

}



func DeleteTask(c *gin.Context) {
	var task  []models.Task

	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "task not found"})
		return
	}

	db.DB.Delete(&task)
	c.JSON(http.StatusNoContent, nil)
}


func FilterTasksByStatus(c *gin.Context) {

	status := c.Query("status")

	var tasks []models.Task
	query := db.DB.Model(&models.Task{})

	if status != ""{
		query = query.Where("status = ?", status)

	}

	if err := query.Find(&tasks).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), })
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks,})
	
}



func PatchTask(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var input models.PatchTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, task := range tasks {
		if task.ID == uint(id) {

			if input.Title != nil {
				tasks[i].Title = *input.Title
			}

			if input.Description != nil {
				tasks[i].Description = *input.Description
			}

			if input.Status != nil {
				tasks[i].Status = *input.Status
			}

			c.JSON(http.StatusOK, gin.H{
				"data": tasks[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "task not found",})
}