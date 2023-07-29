package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"playhouse-server/env"
	"playhouse-server/repo"
	"playhouse-server/responsebody"
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
		panic(responsebody.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	return tokenStr
}

func IsSessionTokenValid(tokenStr string) (bool, int) {
	authError := responsebody.ResponseErr{
		Code:    http.StatusForbidden,
		ErrBody: errors.New("not a valid token"),
	}
	if tokenStr == "" {
		panic(authError)
	}

	secret := env.JWT_SECRET()
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		panic(responsebody.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		sessionRepo := repo.SessionRepo()
		isSessionValid := sessionRepo.IsSessionAvailable(claims.SessionID)
		return isSessionValid, claims.SessionID
	} else {
		panic(authError)
	}

	return false, -1
}
