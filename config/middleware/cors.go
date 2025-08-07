package middleware

import (
	"chatappv2/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	cfg, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	origins := strings.Split(cfg.Server.AllowOrigin, ",")
	corsConfig := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(corsConfig)
}
