package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"shark-auth/apis/handlers"
	"shark-auth/apis/middlewares"
	"shark-auth/pkg/user"
)

func API(db *sqlx.DB, redisClient *redis.Client) *gin.Engine {
	logrus.Info("starting server")

	r := gin.Default()
	r.Use(middlewares.PanicHandlerMiddleware())
	r.POST("/signup", handlers.Signup(user.NewUserRepository(db)))
	r.POST("/token", handlers.CreateToken(user.NewUserRepository(db), db))
	r.GET("/welcome", handlers.Welcome)
	r.PATCH("/token", handlers.RefreshToken(db))
	r.DELETE("/token", handlers.DeleteToken(db, redisClient))
	return r
}
