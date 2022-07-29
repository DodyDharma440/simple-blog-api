package config

import (
	"final-project/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	username = "postgres"
	password = "admin"
	dbName   = "sanber_go_final_project"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=localhost user=%v password=%v dbname=%v port=5432 sslmode=disable TimeZone=Asia/Shanghai", username, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database is connected")
	db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Article{},
		&models.ArticleTag{},
		&models.Tag{},
		&models.ArticleCategory{},
		&models.ArticleComment{},
	)
	return db
}
