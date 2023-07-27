package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/netSpot/goApi/initializers"
	"github.com/netSpot/goApi/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)



func GetUserByIDWithCache(c *gin.Context) {
    userIDStr := c.Param("id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := fetchUserByIDWithCache(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func fetchUserByIDWithCache(userID uint) (models.User, error) {
    var user models.User

    data, err := initializers.RDB.Get(initializers.Ctx, "user:"+strconv.Itoa(int(userID))).Result()
    if err == nil {
        err = json.Unmarshal([]byte(data), &user)
        if err != nil {
            return models.User{}, err
        }
        return user, nil
    }

    if err == redis.Nil {
        result := initializers.DB.First(&user, userID)
        if result.Error != nil {
            return models.User{}, result.Error
        }

        jsonData, err := json.Marshal(user)
        if err != nil {
            return models.User{}, err
        }

        err = initializers.RDB.Set(initializers.Ctx, "user:"+strconv.Itoa(int(userID)), jsonData, 5*time.Minute).Err()
        if err != nil {
            return models.User{}, err
        }

        return user, nil
    }

    return models.User{}, err
}

func GetUsers(c *gin.Context) {
	var users []models.User
	var data []models.User
	result := initializers.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	for _,user :=range users{
		user.Password=""
		data = append(data, user)
	}

	c.JSON(http.StatusOK, gin.H{"users": data})
}



func DeleteUser(c *gin.Context) {
	userID := c.Param("id")


	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	result = initializers.DB.Unscoped().Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")


	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	hash,err:= bcrypt.GenerateFromPassword([]byte(updatedUser.Password),10)
	if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
			"error":"Faild to hash password",
		})
			return
	}

	user.Full_Name = updatedUser.Full_Name
	user.Username = updatedUser.Username
	user.Password=string(hash)
	result = initializers.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Profile(c *gin.Context){
		 user,ok:=c.Get("user")
		if !ok{
		c.JSON(http.StatusUnauthorized, gin.H{"user":"Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user":user})


}