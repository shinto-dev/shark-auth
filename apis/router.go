package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"shark-auth/user"
)

func Router() *gin.Engine {
	logrus.Info("starting server")

	var users = map[string]string{
		"user1": "password1",
		"user2": "password2",
	}

	r := gin.Default()
	r.POST("/token", GetToken(user.NewSampleRepository(users)))
	r.GET("/welcome", welcome)
	r.PATCH("/token", Refresh)
	r.DELETE("/token", DeleteToken)
	return r
}
