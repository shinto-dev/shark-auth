package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"shark-auth/apis/handlers"
)

func PanicHandlerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				handlers.HandleError(context.Background(), w, getPanicErr(r))
				return
			}
		}()
		h.ServeHTTP(w, r)
	})
}
func getPanicErr(recoverErr interface{}) error {
	err, ok := recoverErr.(error)
	if !ok {
		return errors.New(fmt.Sprintf("%v", recoverErr))
	}
	return err
}
