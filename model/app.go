package model

import (
	"time"
)

type App struct {
	ID          int       `json:"id", gorm:"primary_key"`
	Name        string    `json:"name" gorm:"not null`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`
}
