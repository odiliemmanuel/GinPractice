package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gotask/db"
	"gotask/models"

	"github.com/gin-gonic/gin"
)





func GetTasks(c *gin.Context) {
	var tasks  []models.Task
	user, _ := c.MustGet("currentUser").(models.User)

	query := db.DB.Where("user_id = ?", user.ID).Find(&tasks)

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

	user, _ := c.MustGet("currentUser").(models.User)

	newTask := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusPending,
		DueDate:     input.DueDate,
		CreatedAt:   time.Now(),
		UserID:      user.ID,
	}

	if err := db.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": newTask})
}



func GetTask(c *gin.Context) {
	var task  models.Task
	user, _ := c.MustGet("currentUser").(models.User)


	if err := db.DB.First(&task, c.Param("id")).Error;  err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}


	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
		return
	}
	

	
	c.JSON(http.StatusOK, gin.H{"data": task})

}

func UpdateTask(c *gin.Context) {
	var task models.Task
	user, _ := c.MustGet("currentUser").(models.User)

	// 1. Find the task
	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// 2. Enforce user ownership
	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
		return
	}

	// 3. Bind incoming JSON
	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Validate Status Enum
	if input.Status != "" {
		isValid := false
		allowedStatuses := []string{"pending", "in_progress", "done"}

		for _, s := range allowedStatuses {
			if string(input.Status) == s {
				isValid = true
				break
			}
		}

		if !isValid {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "invalid status value. Allowed choices are: 'pending', 'in_progress', 'done'",
			})
			return
		}
	}

	// 5. Save changes to DB
	if err := db.DB.Model(&task).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}



func GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	query := db.DB.Model(&models.Task{})

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Pagination defaults
	page := 1
	limit := 10

	// Read query params
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	// Apply pagination
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"page":  page,
		"limit": limit,
	})
}


func DeleteTask(c *gin.Context) {
	var task  models.Task
	user, _ := c.MustGet("currentUser").(models.User)

	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error" : "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
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
	var task models.Task
	user, _ := c.MustGet("currentUser").(models.User)

	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
		return
	}


	var input models.PatchTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	
	if err := db.DB.Model(&task).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}