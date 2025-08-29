package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/aryadisastra/authentication/internal/handlers"
	"github.com/aryadisastra/authentication/internal/httpx"
	"github.com/aryadisastra/authentication/internal/middleware"
)

func New(db *gorm.DB, jwtSecret string, ttl int) *gin.Engine {
	r := gin.Default()
	r.Use(httpx.RecoverJSON(), httpx.NotFoundAsJSON())

	r.GET("/health", func(c *gin.Context) { httpx.OK(c, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", handlers.Register(db))
		api.POST("/auth/login", handlers.Login(db, jwtSecret, ttl))

		auth := api.Group("/auth")
		auth.Use(middleware.AuthRequired(jwtSecret))
		{
			auth.GET("/me", handlers.Me(db))
		}
	}
	return r
}
