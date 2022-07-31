package routes

import (
	"final-project/controllers"
	"final-project/middlewares"
	"net/http"

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
	authRoutes := r.Group("/")
	authRoutes.Use(middlewares.JwtAuth())
	authRoutes.PATCH("/change-password", controllers.ChangePassword)
	authRoutes.GET("/my-profile", controllers.MyProfile)

	// users
	userRoutes := r.Group("/users")
	userRoutes.Use(middlewares.JwtAuth(), middlewares.AdminOnly())
	userRoutes.GET("/", controllers.GetUsers)
	userRoutes.GET("/:id", controllers.GetUser)
	userRoutes.POST("/", controllers.CreateUser)
	userRoutes.PUT("/:id", controllers.UpdateUser)
	userRoutes.DELETE("/:id", controllers.DeleteUser)

	// categories
	categoriesRoutes := r.Group("/categories")
	categoriesRoutes.Use(middlewares.JwtAuth(), middlewares.AdminOnly())
	r.GET("/categories", controllers.GetCategories)
	r.GET("/categories/:id", controllers.GetCategory)
	categoriesRoutes.POST("/", controllers.CreateCategory)
	categoriesRoutes.PUT("/:id", controllers.UpdateCategory)
	categoriesRoutes.POST("/:id", controllers.DeleteCategory)

	// tags
	tagsRoutes := r.Group("/tags")
	tagsRoutes.Use(middlewares.JwtAuth(), middlewares.AdminOnly())
	r.GET("/tags", controllers.GetTags)
	r.GET("/tags/:id", controllers.GetTag)
	tagsRoutes.POST("/", controllers.CreateTag)
	tagsRoutes.PUT("/:id", controllers.UpdateTag)
	tagsRoutes.POST("/:id", controllers.DeleteTag)

	// articles
	articleRoutes := r.Group("/articles")
	articleRoutes.Use(middlewares.JwtAuth(), middlewares.AdminOnly())
	r.GET("/articles", controllers.GetArticles)
	r.GET("/articles/:id", controllers.GetArticle)
	r.GET("/articles/slug/:slug", controllers.GetArticleBySlug)
	articleRoutes.POST("/", controllers.CreateArticle)
	articleRoutes.PUT("/:id", controllers.UpdateArticle)
	articleRoutes.DELETE("/:id", controllers.DeleteArticle)
	articleRoutes.PATCH("/publish/:id", controllers.PublishArticle)
	articleRoutes.PATCH("/unpublish/:id", controllers.UnpublishArticle)

	commentRoutes := r.Group("/articles")
	commentRoutes.Use(middlewares.JwtAuth())
	r.GET("/articles/:id/comments", controllers.GetComments)
	commentRoutes.POST("/:id/comments", controllers.CreateComment)
	commentRoutes.DELETE("/comments/:id", controllers.DeleteComment)
	r.GET("/articles/comments/:id/replies", controllers.GetReplyComments)
	commentRoutes.POST("/comments/:id/replies", controllers.CreateReplyComment)
	commentRoutes.DELETE("/comments/replies/:id", controllers.DeleteReplyComment)

	// docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// file
	r.StaticFS("/file", http.Dir("public"))
	return r
}
