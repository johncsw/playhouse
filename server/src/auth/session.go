package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"playhouse-server/repository"
	"playhouse-server/util"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	sessionID int
}
type SessionAuthenticator struct {
	RepoFact *repository.Factory
	Env      *util.Env
}

var (
	sessionautheticator *SessionAuthenticator
)

func NewSessionAuthenticator() *SessionAuthenticator {
	if sessionautheticator == nil {
		sessionautheticator = &SessionAuthenticator{
			RepoFact: repository.NewFactory(),
			Env:      util.NewEnv(),
		}
	}
	return sessionautheticator
}

func (a SessionAuthenticator) InitializeSession() string {
	sessionRepo := a.RepoFact.NewSessionRepo()
	s := sessionRepo.NewSession()
	sessionID := s.ID
	claims := JWTClaims{
		sessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(a.Env.JWTSecret())
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		panic(util.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	return tokenStr
}

func (a SessionAuthenticator) IsJWTValid(tokenStr string) bool {
	authError := util.ResponseErr{
		Code:    http.StatusForbidden,
		ErrBody: errors.New("not a valid tokenStr"),
	}
	if tokenStr == "" {
		panic(authError)
	}

	secret := a.Env.JWTSecret()
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		panic(util.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		sessionRepo := a.RepoFact.NewSessionRepo()
		return sessionRepo.IsSessionAvailable(claims.sessionID)
	} else {
		panic(authError)
	}

	return false
}
