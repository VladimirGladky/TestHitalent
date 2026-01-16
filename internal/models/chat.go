package models

import "time"

type Chat struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" validate:"required,min=1,max=200"`
	CreatedAt time.Time `json:"created_at"`
}
