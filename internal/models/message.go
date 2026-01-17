package models

import "time"

type Message struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	ChatID    int       `json:"chat_id" gorm:"not null;constraint:OnDelete:CASCADE"`
	Text      string    `json:"text" gorm:"type:text;not null" validate:"required,min=1,max=5000"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
