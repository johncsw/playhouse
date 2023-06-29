package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type DB struct {
	Connection *gorm.DB
}

func NewDB() *DB {
	config := postgres.Config{
		DSN:                  os.Getenv("DB_DSN"),
		PreferSimpleProtocol: true,
	}
	connection, err := gorm.Open(postgres.New(config), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	return &DB{
		Connection: connection,
	}
}
