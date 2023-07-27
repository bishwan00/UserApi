package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/netSpot/goApi/initializers"
	"github.com/netSpot/goApi/models"
)

func RequireAuth(c *gin.Context){
	tokenString,err:=c.Cookie("Authorization")

	if err!=nil{
		c.AbortWithStatus(http.StatusUnauthorized)
	}


token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(os.Getenv("SECERET")), nil
})

if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	if float64(time.Now().Unix())>claims["exp"].(float64){
				c.AbortWithStatus(http.StatusUnauthorized)
	}

	var user models.User
	initializers.DB.First(&user,claims["sub"])

		if user.ID==0{
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		user.Password=""
		c.Set("user",user)
		c.Next()



} else {
		c.AbortWithStatus(http.StatusUnauthorized)
}

}