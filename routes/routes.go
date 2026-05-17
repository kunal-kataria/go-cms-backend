package routes

import (
	"github.com/kunal-kataria/go-cms-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitializeRoutes sets up all API routes
func InitializeRoutes(router *gin.Engine, db *gorm.DB) {
	router.Use(func(c *gin.Context) { // Middleware: add DB connection to the gin.context
		c.Set("db", db)
		c.Next() //pases connection to next middlewares
	})

	v1 := router.Group("api/v1")
	{
		// Page Routes
		pages := v1.Group("/pages")
		{
			pages.GET("/", controllers.GetPages)
			pages.GET("/:id", controllers.GetPage)
			pages.POST("/", controllers.CreatePage)
			pages.PUT("/:id", controllers.UpdatePage)
			pages.DELETE("/:id", controllers.DeletePage)
		}

		// Post Routes
		posts := v1.Group("/posts")
		{
			posts.GET("/", controllers.GetPosts)
			posts.GET("/:id", controllers.GetPost)
			posts.POST("/", controllers.CreatePost)
			posts.PUT("/:id", controllers.UpdatePost)
			posts.DELETE("/:id", controllers.DeletePost)
		}

		// Media Routes
		media := v1.Group("/media")
		{
			media.GET("/", controllers.GetMedia)
			media.GET("/:id", controllers.GetMediaByID)
			media.POST("/", controllers.CreateMedia)
			media.DELETE("/:id", controllers.DeleteMedia)
		}
	}
}
