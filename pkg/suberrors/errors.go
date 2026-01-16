package suberrors

import "errors"

var (
	ErrInvalidChatId = errors.New("invalid chat id format")
	ErrNotPositiveChatId = errors.New("chat id must be positive")
)