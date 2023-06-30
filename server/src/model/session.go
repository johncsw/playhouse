package model

import "time"

type Session struct {
	ID          int `gorm:"primaryKey"`
	IsAvailable bool
	DueAt       *time.Time
	CreatedAt   *time.Time
}

func (Session) TableName() string {
	return "session"
}
