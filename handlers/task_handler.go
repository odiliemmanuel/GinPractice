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
	var input models.CreateTaskInput

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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
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


func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}



	var input models.UpdateTaskInput


	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	mu.Lock()
	defer mu.Unlock()

	for count, task := range tasks {
		if task.ID == uint(id){

			tasks[count].Title = input.Title
			tasks[count].Description = input.Description
			tasks[count].Status = input.Status

			c.JSON(http.StatusOK, gin.H{"data": tasks[count]})
			return
		}
	}


	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})

}



func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for count, task := range tasks{
		if task.ID == uint(id){

			tasks = append(tasks[:count], tasks[count +1:]...)

			c.JSON(http.StatusOK, gin.H{"message": "task deleted",})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": " tasks not found",})


}


func FilterTasksByStatus(c *gin.Context) {
	status := c.Query("status")

	mu.Lock()
	defer mu.Unlock()

	if status == ""{
		c.JSON(http.StatusOK, gin.H{"data": tasks, })
		return
	}

	var filteredTasks []models.Task

	for _, task := range tasks {
		if string(task.Status) == status{
			filteredTasks = append(filteredTasks, task)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": filteredTasks})

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

	c.JSON(http.StatusNotFound, gin.H{
		"error": "task not found",})
}