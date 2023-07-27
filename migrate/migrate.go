package main

import (
	"github.com/netSpot/goApi/initializers"
	"github.com/netSpot/goApi/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main(){
	initializers.DB.AutoMigrate(&models.User{})
}