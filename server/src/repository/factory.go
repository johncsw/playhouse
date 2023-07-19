package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"playhouse-server/util"
)

var (
	db *gorm.DB
)

type Factory struct {
	db *gorm.DB
}

func NewFactory() *Factory {
	if db == nil {
		env := util.NewEnv()
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

	return &Factory{db: db}
}

func (f *Factory) NewTransaction() transaction {
	return transaction{db: f.db}
}

func (f *Factory) NewSessionRepo() *sessionrepo {
	return &sessionrepo{db: f.db}
}

func (f *Factory) NewVideoRepo() *videorepo {
	return &videorepo{db: f.db}
}

func (f *Factory) NewChunkRepo() *chunkrepo {
	return &chunkrepo{db: f.db}
}
