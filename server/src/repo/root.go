package repo

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"playhouse-server/env"
)

var (
	db *gorm.DB
)

func Init() {
	driverConfig := postgres.Config{
		DSN:                  env.DSN(),
		PreferSimpleProtocol: true,
	}

	connection, err := gorm.Open(postgres.New(driverConfig), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	result := connection.Raw("select 1;")
	if result.Error != nil {
		panic(fmt.Sprintf("Fail to connect to database: %v", result.Error))
	} else {
		log.Println("Database connection successful.")
	}

	db = connection
}

func NewTransaction(statements func(*gorm.DB) error) error {
	return db.Transaction(statements)
}

var (
	sessionRepo *sessionrepo
	videoRepo   *videorepo
	chunkRepo   *chunkrepo
	userRepo    *userrepo
)

func SessionRepo() *sessionrepo {
	if sessionRepo == nil {
		sessionRepo = &sessionrepo{}
	}
	return sessionRepo
}

func VideoRepo() *videorepo {
	if videoRepo == nil {
		videoRepo = &videorepo{}
	}
	return videoRepo
}

func ChunkRepo() *chunkrepo {
	if chunkRepo == nil {
		chunkRepo = &chunkrepo{}
	}
	return chunkRepo
}

func UserRepo() *userrepo {
	if userRepo == nil {
		userRepo = &userrepo{}
	}
	return userRepo
}
