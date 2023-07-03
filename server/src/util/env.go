package util

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"playhouse-server/middleware"
	"strconv"
)

func LoadEnv() {
	// export ENV_PATH manually if you're on local development
	// env_local: ../conf/.env_local
	// env_docker: ./conf/.env_docker
	envPath := os.Getenv("ENV_PATH")
	if err := godotenv.Load(envPath); err != nil {
		panic(err)
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
