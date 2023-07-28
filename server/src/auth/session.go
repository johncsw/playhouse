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
type SessionAuthenticator struct {
}

var (
	sessionautheticator *SessionAuthenticator
)

func NewSessionAuthenticator() *SessionAuthenticator {
	if sessionautheticator == nil {
		sessionautheticator = &SessionAuthenticator{}
	}
	return sessionautheticator
}

func (a SessionAuthenticator) InitializeSession() string {
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

func (a SessionAuthenticator) IsJWTValid(tokenStr string) bool {
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
		return sessionRepo.IsSessionAvailable(claims.SessionID)
	} else {
		panic(authError)
	}

	return false
}

func (a SessionAuthenticator) getClaims(tokenStr string) *JWTClaims {
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
		return claims
	} else {
		panic(authError)
	}
}

func (a *SessionAuthenticator) GetSessionId(r *http.Request) int {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		tokenStr = r.URL.Query().Get("token")
	}

	claims := a.getClaims(tokenStr)
	return claims.SessionID
}
