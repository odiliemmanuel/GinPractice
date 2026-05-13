package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"gotask/models"

	"github.com/gin-gonic/gin"
)

var (
	tasks  []models.Task
	nextID uint = 1
	mu     sync.Mutex
)


func GetTasks(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}


func CreateTask(c *gin.Context) {
	var input models.CreatTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	newTask := models.Task{
		ID:          nextID,
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
	}

	nextID++
	tasks = append(tasks, newTask)

	c.JSON(http.StatusCreated, gin.H{"data": newTask})
}


func GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != ni
	}

	mu.Lock()
	defer mu.Unlock()

	for _, task := range tasks {
		if task.ID == uint(id) {
			c.JSON(http.StatusOK, gin.H{"data": task})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

