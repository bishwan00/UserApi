package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netSpot/goApi/models"
)

func CheckRole(c *gin.Context) {

	user,_ := c.Get("user")
	if user.(models.User).Role !="admin"{
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized,gin.H{"message":"Unauthorized"})
	}
	c.Next()
}