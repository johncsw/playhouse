package util

import (
	"errors"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"playhouse-server/middleware"
	"strconv"
)

func LoadEnv() {
	if err := godotenv.Load("../conf/.env"); err != nil {
		panic(errors.New("error loading .env file"))
	}
}

func EnvGetSessionTTLHour() int {
	sessionTTLHour, err := strconv.Atoi(os.Getenv("APP_SESSION_TTL_HOUR"))
	if err != nil {
		panic(middleware.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}
	return sessionTTLHour
}

func EnvGetDSN() string {
	return os.Getenv("DB_DSN")
}

func EnvGetJWTSecret() string {
	return os.Getenv("APP_JWT_SECRET")
}
