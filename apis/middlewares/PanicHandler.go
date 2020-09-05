package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"shark-auth/apis/handlers"
	"shark-auth/pkg/errorcode"
)

func PanicHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("stacktrace", fmt.Sprintf("%+v", getPanicErr(r))).
					Error("Unexpected error")
				c.JSON(http.StatusBadRequest,
					handlers.NewErrorResponse(errorcode.ERROR_INTERNAL_ERROR, "internal-error"))
				return
			}
		}()
		c.Next()
	}
}
func getPanicErr(recoverErr interface{}) error {
	err, ok := recoverErr.(error)
	if !ok {
		return errors.New(fmt.Sprintf("%v", recoverErr))
	}
	return err
}
