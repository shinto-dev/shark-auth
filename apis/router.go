package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"shark-auth/apis/handlers"
	"shark-auth/pkg/user"
)

func API(db *sqlx.DB) *gin.Engine {
	logrus.Info("starting server")

	r := gin.Default()
	r.POST("/token", handlers.GetToken(user.NewUserRepository(db)))
	r.GET("/welcome", handlers.Welcome)
	r.PATCH("/token", handlers.Refresh)
	r.DELETE("/token", handlers.DeleteToken)
	return r
}
