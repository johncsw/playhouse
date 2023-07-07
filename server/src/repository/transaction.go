package repository

import (
	"gorm.io/gorm"
)

type transaction struct {
	db *gorm.DB
}

func (tx *transaction) Execute(statements func(*gorm.DB) error) error {
	err := tx.db.Transaction(statements)
	return err
}
