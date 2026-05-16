package handlers

import (
	"net/http"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gotask/db"
	"gotask/models"
)


func RegisterUser(c*gin.Context){
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email: input.Email, 
		Password: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data" : gin.H{"id": user.ID, "username": user.Username, "email": user.Email}})
}


func LoginUser(c *gin.Context){
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signed})


}

var jwtSecret = []byte("your_secret_key")

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func RefreshToken(c *gin.Context) {
	var input RefreshTokenRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the token, ignoring expiration constraints explicitly
	token, err := jwt.Parse(input.Token, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithExpirationRequired())

	// Validate errors safely
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload structure"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to extract claims"})
		return
	}

	// Fetch standard claim numeric dates safely
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid metadata parameters"})
		return
	}
	
	expiredAt := time.Unix(int64(expFloat), 0)

	// Max allowable Grace Period Window: 7 days
	if time.Since(expiredAt) > (7 * 24 * time.Hour) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token grace window has closed, please reauthenticate"})
		return
	}

	// Generate and issue a new 24-Hour Token
	newClaims := jwt.MapClaims{
		"sub": claims["sub"], // User ID
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := newToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate a new token token string"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}