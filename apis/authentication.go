package apis

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"shark-auth/user"
)

const Token = "token"

var jwtKey = []byte("my_secret_key")

type GetTokenRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type GetTokenResponse struct {
	 AccessToken string `json:"access_token"`
	 RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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

		expireAt := time.Now().Add(5 * time.Minute)
		tokenString, err := createJwtToken(getTokenRequest, expireAt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := GetTokenResponse{AccessToken: tokenString}

		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func createJwtToken(creds GetTokenRequest, expireAt time.Time) (string, error) {
	claims := Claims{
		Username: creds.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
func welcome(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(Token)
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
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logrus.Error("parsing token failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Hello world"))
}

// takes the previous token (which is still valid), and returns a new token with a renewed expiry time.
// To minimize misuse of a JWT, the expiry time is usually kept in the order of a few minutes.
// Typically the client application would refresh the token in the background.
func refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(Token)
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
	var claims Claims
	tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logrus.Error("parsing tkn failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new tkn will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new tkn as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
