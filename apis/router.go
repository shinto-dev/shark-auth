package apis

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"shark-auth/user"
)

func Router() *mux.Router {
	logrus.Info("starting server")

	var users = map[string]string{
		"user1": "password1",
		"user2": "password2",
	}

	r := mux.NewRouter()
	r.HandleFunc("/token", GetToken(user.NewSampleRepository(users))).
		Methods("POST")
	r.HandleFunc("/welcome", welcome)
	r.HandleFunc("/refresh", refresh)
	r.HandleFunc("/token", refresh).Methods("DELETE")

	return r
}
