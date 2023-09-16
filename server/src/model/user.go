package model

import "time"

type User struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt *time.Time
	Email     string
}

func (User) TableName() string {
	return "user"
}
