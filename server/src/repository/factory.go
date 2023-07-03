package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"playhouse-server/util"
)

type Factory struct {
	DB *gorm.DB
}

func NewFactory() *Factory {
	driverConfig := postgres.Config{
		DSN:                  util.EnvGetDSN(),
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(driverConfig), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	f := &Factory{DB: db}
	result := f.DB.Raw("select 1;")
	if result.Error != nil {
		panic(fmt.Sprintf("Fail to connect to database: %v", result.Error))
	} else {
		log.Println("Database connection successful.")
	}

	return f
}

func (f *Factory) NewSessionRepo() *SessionRepo {
	return &SessionRepo{DB: f.DB}
}
