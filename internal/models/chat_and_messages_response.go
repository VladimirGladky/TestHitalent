package models

type ChatAndMessagesResponse struct {
	*Chat
	Messages []*Message `json:"messages"`
}
