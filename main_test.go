package main

import (
	"gotask/db"
	"gotask/models"
	"gotask/routes"
	"gotask/config"
	

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	
	db.Connect(config.Load())

	r := gin.Default()
	routes.SetupRoutes(r)

	return r
}

func generateTestToken(userID uint) string {
	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().Add(time.Hour).Unix(),
		},
	)

	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString
}


func TestCreateTask(t *testing.T) {
	router := setupRouter()

	user := models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	db.DB.Create(&user)

	token := generateTestToken(user.ID)

	body := `{
		"title": "Learn Github Actions",
		"completed": false
	}`

	request, _ := http.NewRequest(
		"POST",
		"/api/v1/tasks",
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, request)

	if w.Code != http.StatusCreated {
		t.Errorf(
			"Expected %d got %d",
			http.StatusCreated,
			w.Code,
		)
	}
}