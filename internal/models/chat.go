package models

import "time"

type Chat struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null" validate:"required,min=1,max=200"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
