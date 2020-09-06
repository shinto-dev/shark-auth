package apis

import (
	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"shark-auth/apis/handlers"
	"shark-auth/apis/middlewares"
)

func API(db *sqlx.DB, redisClient *redis.Client) http.Handler {
	logrus.Info("starting server")
	r := mux.NewRouter()
	r.Use(middlewares.PanicHandlerMiddleware)
	r.HandleFunc("/signup", handlers.HandleUserSignup(db)).Methods(http.MethodPost)
	r.HandleFunc("/token", handlers.HandleTokenCreate(db)).Methods(http.MethodPost)
	r.HandleFunc("/welcome", handlers.HandleWelcome(redisClient)).Methods(http.MethodGet)
	r.HandleFunc("/token", handlers.HandleTokenRefresh(db)).Methods(http.MethodPatch)
	r.HandleFunc("/token", handlers.HandleTokenDelete(db, redisClient)).Methods(http.MethodDelete)

	return r
}
