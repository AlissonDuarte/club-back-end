package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  )
  
func NewDb() *gorm.DB {
	dsn := "alisson:0601@tcp(127.0.0.1:3306)/clube?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return db
  }