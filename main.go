package main

import (
	"final-project/config"
	"final-project/docs"
	"final-project/routes"
	"final-project/utils"
	"log"
	"os"

	_ "final-project/docs"

	"github.com/joho/godotenv"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err.Error())
		// fmt.Println(err.Error())
	}

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Blog API"
	docs.SwaggerInfo.Description = "This API Blog."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = utils.GetEnv("SWAGGER_HOST", "localhost:8080")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db := config.ConnectDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := routes.SetupRouter(db)
	r.Run()
}
