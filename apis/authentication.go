package apis

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"shark-auth/autherrors"
	"shark-auth/token"
	"shark-auth/user"
)

const TOKEN = "token"

type GetTokenRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type GetTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetToken(userRepo user.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getTokenRequest GetTokenRequest
		err := json.NewDecoder(r.Body).Decode(&getTokenRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !userRepo.IsValid(getTokenRequest.UserName, getTokenRequest.Password) {
			logrus.Error("password does not match")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tkn, err := token.CreateJwtToken(getTokenRequest.UserName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := GetTokenResponse{AccessToken: tkn}

		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(TOKEN)
	if err != nil {
		if err == http.ErrNoCookie {
			logrus.Error("cookie not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value

	claims, err := token.ParseJwtToken(tokenString)
	if err != nil {
		if err == autherrors.ErrAuthenticationFailed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logrus.Infof("request received from user: %s", claims.Username)

	w.Write([]byte("Hello world"))
}

// takes the previous token (which is still valid), and returns a new token with a renewed expiry time.
// To minimize misuse of a JWT, the expiry time is usually kept in the order of a few minutes.
// Typically the client application would refresh the token in the background.
func refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(TOKEN)
	if err != nil {
		if err == http.ErrNoCookie {
			logrus.Error("cookie not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := cookie.Value
	claims, err := token.ParseJwtToken(tknStr)
	if err != nil {
		if err == autherrors.ErrAuthenticationFailed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new tkn will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jwtToken, err := token.CreateJwtToken(claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := GetTokenResponse{AccessToken: jwtToken}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
