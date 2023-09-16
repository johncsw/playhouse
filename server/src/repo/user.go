package repo

import (
	"errors"
	"gorm.io/gorm"
	"playhouse-server/model"
	"playhouse-server/request"
	"time"
)

type userrepo struct {
}

func (r *userrepo) NewUser(b *request.AuthRegistrationBody, tx *gorm.DB) (*model.User, error) {
	var executor *gorm.DB
	if tx == nil {
		executor = db
	} else {
		executor = tx
	}
	now := time.Now().UTC()
	usr := &model.User{
		Email:     b.Email,
		CreatedAt: &now,
	}

	err := executor.Create(usr).Error

	return usr, err
}

func (r *userrepo) GetUserByEmail(email string, tx *gorm.DB) (*model.User, error) {
	var executor *gorm.DB
	if tx == nil {
		executor = db
	} else {
		executor = tx
	}

	var u model.User
	err := executor.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userrepo) DeleteAll(tx *gorm.DB) error {
	var executor *gorm.DB
	if tx == nil {
		executor = db
	} else {
		executor = tx
	}

	return executor.Where("id > 0").Delete(&model.User{}).Error
}
