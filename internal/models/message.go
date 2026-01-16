package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	Text      string    `json:"text" validate:"required,min=1,max=5000"`
	CreatedAt time.Time `json:"created_at"`
}
