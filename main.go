package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/netSpot/goApi/controllers"
	"github.com/netSpot/goApi/initializers"
	"github.com/netSpot/goApi/middleware"
)

func init(){
initializers.LoadEnvVariables()
initializers.ConnectToDB()
initializers.SetupRedis()
}

func main() {
	r := gin.Default()
	 logFile, err := os.Create("server.log")
    if err != nil {
        log.Fatal("Failed to create log file:", err)
    }

	    gin.DefaultWriter = io.MultiWriter(logFile)
		    r.Use(gin.Logger())



		r.POST("/upload",middleware.RequireAuth,controllers.ImageUpload)
		r.POST("/users/signup", controllers.Signup)
		r.POST("/users/login", controllers.Login)
		r.GET("/users",middleware.RequireAuth,controllers.GetUsers)
		r.DELETE("/users/:id",middleware.RequireAuth,middleware.CheckRole,controllers.DeleteUser)
		r.PATCH("/users/:id",middleware.RequireAuth,controllers.UpdateUser)
		r.GET("/users/:id", middleware.RequireAuth,controllers.GetUserByIDWithCache)
		r.GET("/users/profile",middleware.RequireAuth,controllers.Profile)

	r.Run() 
}