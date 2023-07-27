package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Full_Name string `json:"full_name" binding:"required"`
	Username string `json:"userName" binding:"required,alpha"`
	Password string	`json:"password,omitempty" binding:"required,min=8,max=15"`
	Status string	`json:"status"  gorm:"type:enum('active','inactive');default:inactive"`
	Role string		`json:"role" gorm:"type:enum('admin','operator');default:'operator'"`
}


