package env

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"playhouse-server/response"
	"strconv"
)

// Load function would load .env file from the path specified in ENV_PATH
// export ENV_PATH manually if you're on local development
// env_local: ../conf/.env_local
// env_docker: ./conf/.env_docker
func Load() {
	envPath := os.Getenv("ENV_PATH")
	if err := godotenv.Load(envPath); err != nil {
		panic(err)
	}
}

func SESSION_TTL_HOUR() int {
	sessionTTLHour, err := strconv.Atoi(os.Getenv("APP_SESSION_TTL_HOUR"))
	if err != nil {
		panic(response.Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
	}
	return sessionTTLHour
}

func DSN() string {
	return os.Getenv("DB_DSN")
}

func JWT_SECRET() string {
	return os.Getenv("APP_JWT_SECRET")
}

func CORS_ALLOWED_WEBSITE() string {
	return os.Getenv("CORS_ALLOWED_WEBSITE")
}

func CHUNK_STORAGE_PATH() string {
	return os.Getenv("CHUNK_STORAGE_PATH")
}

func SHELL_PATH() string {
	return os.Getenv("SHELL_PATH")
}

func CLIENT_URL() string {
	return os.Getenv("CLIENT_URL")
}
