package repo

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"playhouse-server/env"
	"strings"
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
		panic(fmt.Sprintf("fail to connect to database: %v", result.Error))
	} else {
		log.Println("database connection successful.")
	}

	db = connection
}

//go:embed schema/*.sql
var embeddedSQLs embed.FS

func SetUpSchema() {
	var db *sql.DB
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	dsn := strings.Replace(env.DSN(), "dbname=playhouse", "dbname=defaultdb", 1)
	conn, openErr := sql.Open("pgx", dsn)
	if openErr != nil {
		panic(openErr)
	}
	db = conn

	goose.SetBaseFS(embeddedSQLs)
	if dialectErr := goose.SetDialect("pgx"); dialectErr != nil {
		panic(dialectErr)
	}

	if migrateErr := goose.Up(db, "schema", goose.WithNoVersioning()); migrateErr != nil {
		panic(migrateErr)
	}
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
