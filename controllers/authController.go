package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/netSpot/goApi/initializers"
	"github.com/netSpot/goApi/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body models.User
	var result models.User
	var err error

	if err =  c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to hash password",
		})
		return
	}
	body.Password=string(hash)
	 err = initializers.DB.Create(&body).Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to create user" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "succese",
	})
}

func Login(c *gin.Context) {

	var body struct {
		UserName string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to read body",
		})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "username = ?", body.UserName)
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username or password" + result.Error.Error(),
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("seceret")))
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"accese_token":tokenString})
}
