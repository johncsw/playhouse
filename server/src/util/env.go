package util

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
)

type Env struct {
}

var (
	env *Env
)

func NewEnv() *Env {
	if env == nil {
		env = &Env{}
	}
	return env
}

func (Env) Load() {
	// export ENV_PATH manually if you're on local development
	// env_local: ../conf/.env_local
	// env_docker: ./conf/.env_docker
	envPath := os.Getenv("ENV_PATH")
	if err := godotenv.Load(envPath); err != nil {
		panic(err)
	}
}

func (Env) SessionTTLHour() int {
	sessionTTLHour, err := strconv.Atoi(os.Getenv("APP_SESSION_TTL_HOUR"))
	if err != nil {
		panic(ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}
	return sessionTTLHour
}

func (Env) DSN() string {
	return os.Getenv("DB_DSN")
}

func (Env) JWTSecret() string {
	return os.Getenv("APP_JWT_SECRET")
}

func (Env) CORS_ALLOWED_WEBSITE() string {
	return os.Getenv("CORS_ALLOWED_WEBSITE")
}
