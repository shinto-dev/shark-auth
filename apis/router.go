package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"shark-auth/apis/handlers"
	"shark-auth/apis/middlewares"
)

func API(db *sqlx.DB, redisClient *redis.Client) *gin.Engine {
	logrus.Info("starting server")

	r := gin.Default()
	r.Use(middlewares.PanicHandlerMiddleware())
	r.POST("/signup", handlers.HandleUserSignup(db))
	r.POST("/token", handlers.HandleTokenCreate(db))
	r.GET("/welcome", handlers.HandleWelcome(redisClient))
	r.PATCH("/token", handlers.HandleTokenRefresh(db))
	r.DELETE("/token", handlers.HandleTokenDelete(db, redisClient))
	return r
}
