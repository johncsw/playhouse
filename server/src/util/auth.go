package util

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"playhouse-server/middleware"
	"time"
)

func GenJWT(sessionID int, due *time.Time) string {
	payload := map[string]interface{}{
		"sessionID": sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
		"exp":     due.Unix(),
	})

	secret := []byte(EnvGetJWTSecret())
	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(middleware.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	return tokenString
}
