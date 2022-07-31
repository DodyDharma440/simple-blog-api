package config

import (
	"final-project/models"
	"final-project/utils"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	environment := utils.GetEnv("ENVIRONMENT", "development")

	host := "localhost"
	username := "postgres"
	password := "admin"
	dbName := "sanber_go_final_project"
	port := "5432"
	sslMode := "disable"

	if environment == "production" {
		host = os.Getenv("DB_HOST")
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		dbName = os.Getenv("DB_NAME")
		port = os.Getenv("DB_PORT")
		sslMode = "require"
	}

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=Asia/Shanghai",
		host,
		username,
		password,
		dbName,
		port,
		sslMode,
	)
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
		&models.ReplyArticleComment{},
	)
	return db
}
