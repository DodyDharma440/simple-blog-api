package routes

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// auth
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	// users
	userRoutes := r.Group("/users")
	userRoutes.Use(middlewares.JwtAuth())
	userRoutes.GET("/", controllers.GetUsers)
	userRoutes.GET("/:id", controllers.GetUser)
	userRoutes.POST("/", controllers.CreateUser)
	userRoutes.PUT("/:id", controllers.UpdateUser)
	userRoutes.DELETE("/:id", controllers.DeleteUser)

	// categories
	categoriesRoutes := r.Group("/categories")
	categoriesRoutes.Use(middlewares.JwtAuth())
	r.GET("/categories", controllers.GetCategories)
	r.GET("/categories/:id", controllers.GetCategory)
	categoriesRoutes.POST("/", controllers.CreateCategory)
	categoriesRoutes.PUT("/:id", controllers.UpdateCategory)
	categoriesRoutes.POST("/:id", controllers.DeleteCategory)

	// tags
	tagsRoutes := r.Group("/tags")
	tagsRoutes.Use(middlewares.JwtAuth())
	r.GET("/tags", controllers.GetTags)
	r.GET("/tags/:id", controllers.GetTag)
	tagsRoutes.POST("/", controllers.CreateTag)
	tagsRoutes.PUT("/:id", controllers.UpdateTag)
	tagsRoutes.POST("/:id", controllers.DeleteTag)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
