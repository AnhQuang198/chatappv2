package internal

import (
	"chatappv2/config/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func InitRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1", middleware.CORSMiddleware())
	{
		test := v1.Group("/test")
		{
			test.GET("", func(c *gin.Context) {

			})
		}
	}
	return r
}
