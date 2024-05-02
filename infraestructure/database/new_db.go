package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dsn := "alisson:0601@tcp(127.0.0.1:3306)/clube?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Configurar o pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to configure database connection pool")
	}
	sqlDB.SetMaxIdleConns(10)  // Número máximo de conexões inativas no pool
	sqlDB.SetMaxOpenConns(100) // Número máximo de conexões abertas no pool
}

func NewDb() *gorm.DB {
	return db
}
