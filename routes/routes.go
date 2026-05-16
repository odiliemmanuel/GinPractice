package routes

import (
	"gotask/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	
	r.GET("/tasks", handlers.GetTasks)

	r.GET("/tasks/:id", handlers.GetTask)

	r.POST("/tasks", handlers.CreateTask)

	r.PUT("/tasks/:id", handlers.UpdateTask)

	r.PATCH("/tasks/:id", handlers.PatchTask)

	r.DELETE("/tasks/:id", handlers.DeleteTask)

}
