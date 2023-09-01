package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"playhouse-server/env"
	"playhouse-server/repo"
	"playhouse-server/response"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	SessionID int `json:"sessionID"`
}

func CreateSessionToken() string {
	sessionRepo := repo.SessionRepo()
	s := sessionRepo.NewSession()
	sessionID := s.ID
	claims := JWTClaims{
		SessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := env.JWT_SECRET()
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(response.Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
	}

	return tokenStr
}

func IsSessionTokenValid(tokenStr string) (bool, int) {
	if tokenStr == "" {
		panic(response.Error{
			Code:  http.StatusForbidden,
			Cause: errors.New("not a valid token"),
		})
	}

	secret := env.JWT_SECRET()
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		panic(response.Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		sessionRepo := repo.SessionRepo()
		isAvailable := sessionRepo.IsSessionAvailable(claims.SessionID)
		return isAvailable, claims.SessionID
	} else {
		panic(response.Error{
			Code:  http.StatusForbidden,
			Cause: errors.New("token is not available"),
		})
	}

	return false, -1
}
