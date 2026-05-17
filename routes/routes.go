package routes

import (
	"gotask/handlers"
	"gotask/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("api/v1")


	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/auth/refresh", handlers.RefreshToken)
		auth.POST("/login", handlers.LoginUser)
	}

	tasks := api.Group("/tasks")
	tasks.Use(middleware.RequireAuth)

	{
		tasks.GET("", handlers.GetTasks)

		tasks.GET("/tasks/:id", handlers.GetTask)

		tasks.POST("", handlers.CreateTask)

		tasks.PUT("/tasks/:id", handlers.UpdateTask)

		tasks.PATCH("/tasks/:id", handlers.PatchTask)

		tasks.DELETE("/tasks/:id", handlers.DeleteTask)
	}
	
	

}

