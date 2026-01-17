package service

import (
	"TestHitalent/internal/models"
	"TestHitalent/internal/repository/mocks"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHiTalentService_CreateChatSuccess(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)
	ch := &models.Chat{
		Title: "Test title 1",
	}
	expResp := &models.Chat{
		ID:        1,
		Title:     "Test title 1",
		CreatedAt: time.Now(),
	}
	repo.EXPECT().CreateChat(ch).Return(expResp, nil).Times(1)
	srv := NewHiTalentService(context.Background(), repo)
	chat, err := srv.CreateChat(ch)
	require.NoError(t, err)
	require.Equal(t, expResp, chat)
}

func TestHiTalentService_CreateChatFail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)

	cases := []struct {
		name   string
		chat   *models.Chat
		expErr string
	}{
		{
			name:   "nil value",
			chat:   nil,
			expErr: "chat is nil",
		},
		{
			name: "empty title",
			chat: &models.Chat{
				Title: "",
			},
			expErr: "required",
		},
		{
			name: "whitespace only title",
			chat: &models.Chat{
				Title: "   ",
			},
			expErr: "required",
		},
		{
			name: "title too long",
			chat: &models.Chat{
				Title: strings.Repeat("a", 201),
			},
			expErr: "max",
		},
	}

	srv := NewHiTalentService(context.Background(), repo)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chat, err := srv.CreateChat(tc.chat)
			require.Error(t, err)
			require.Nil(t, chat)
			require.Contains(t, err.Error(), tc.expErr)
		})
	}
}

func TestHiTalentService_GetChatSuccess(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)
	chatID := "1"
	limit := 20
	expResp := &models.ChatAndMessagesResponse{
		Chat: &models.Chat{
			ID:        1,
			Title:     "Test Chat",
			CreatedAt: time.Now(),
		},
		Messages: []*models.Message{
			{ID: 1, ChatID: 1, Text: "Message 1", CreatedAt: time.Now()},
			{ID: 2, ChatID: 1, Text: "Message 2", CreatedAt: time.Now()},
		},
	}

	repo.EXPECT().GetChat(1, limit).Return(expResp, nil).Times(1)
	srv := NewHiTalentService(context.Background(), repo)
	result, err := srv.GetChat(chatID, limit)
	require.NoError(t, err)
	require.Equal(t, expResp, result)
}

func TestHiTalentService_GetChatFail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)

	cases := []struct {
		name   string
		chatID string
		limit  int
		expErr string
	}{
		{
			name:   "invalid chat ID",
			chatID: "invalid",
			limit:  20,
			expErr: "invalid chat id",
		},
		{
			name:   "negative chat ID",
			chatID: "-1",
			limit:  20,
			expErr: "chat id must be positive",
		},
		{
			name:   "zero chat ID",
			chatID: "0",
			limit:  20,
			expErr: "chat id must be positive",
		},
		{
			name:   "chat ID with spaces",
			chatID: " 1 ",
			limit:  20,
			expErr: "invalid chat id",
		},
		{
			name:   "chat ID with letters",
			chatID: "1abc",
			limit:  20,
			expErr: "invalid chat id",
		},
	}

	srv := NewHiTalentService(context.Background(), repo)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := srv.GetChat(tc.chatID, tc.limit)
			require.Error(t, err)
			require.Nil(t, result)
			require.Contains(t, err.Error(), tc.expErr)
		})
	}
}

func TestHiTalentService_CreateMessageSuccess(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)
	chatID := "1"
	msg := &models.Message{
		Text: "Test message",
	}
	expResp := &models.Message{
		ID:        1,
		ChatID:    1,
		Text:      "Test message",
		CreatedAt: time.Now(),
	}

	repo.EXPECT().CreateMessage(1, msg).Return(expResp, nil).Times(1)
	srv := NewHiTalentService(context.Background(), repo)
	result, err := srv.CreateMessage(chatID, msg)
	require.NoError(t, err)
	require.Equal(t, expResp, result)
}

func TestHiTalentService_CreateMessageFail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)

	cases := []struct {
		name    string
		chatID  string
		message *models.Message
		expErr  string
	}{
		{
			name:   "invalid chat ID",
			chatID: "invalid",
			message: &models.Message{
				Text: "Test message",
			},
			expErr: "invalid chat id",
		},
		{
			name:   "negative chat ID",
			chatID: "-1",
			message: &models.Message{
				Text: "Test message",
			},
			expErr: "chat id must be positive",
		},
		{
			name:   "zero chat ID",
			chatID: "0",
			message: &models.Message{
				Text: "Test message",
			},
			expErr: "chat id must be positive",
		},
		{
			name:    "nil message",
			chatID:  "1",
			message: nil,
			expErr:  "message is nil",
		},
		{
			name:   "empty text",
			chatID: "1",
			message: &models.Message{
				Text: "",
			},
			expErr: "required",
		},
		{
			name:   "whitespace only text",
			chatID: "1",
			message: &models.Message{
				Text: "   ",
			},
			expErr: "required",
		},
		{
			name:   "text too long",
			chatID: "1",
			message: &models.Message{
				Text: strings.Repeat("a", 5001),
			},
			expErr: "max",
		},
	}

	srv := NewHiTalentService(context.Background(), repo)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := srv.CreateMessage(tc.chatID, tc.message)
			require.Error(t, err)
			require.Nil(t, result)
			require.Contains(t, err.Error(), tc.expErr)
		})
	}
}

func TestHiTalentService_DeleteChatSuccess(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)
	chatID := "1"

	repo.EXPECT().DeleteChat(1).Return(nil).Times(1)
	srv := NewHiTalentService(context.Background(), repo)
	err := srv.DeleteChat(chatID)
	require.NoError(t, err)
}

func TestHiTalentService_DeleteChatFail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mocks.NewMockHiTalentRepositoryInterface(ctl)

	cases := []struct {
		name   string
		chatID string
		expErr string
	}{
		{
			name:   "invalid chat ID",
			chatID: "invalid",
			expErr: "invalid chat id",
		},
		{
			name:   "negative chat ID",
			chatID: "-1",
			expErr: "chat id must be positive",
		},
		{
			name:   "zero chat ID",
			chatID: "0",
			expErr: "chat id must be positive",
		},
		{
			name:   "chat ID with spaces",
			chatID: " 5 ",
			expErr: "invalid chat id",
		},
		{
			name:   "float chat ID",
			chatID: "1.5",
			expErr: "invalid chat id",
		},
	}

	srv := NewHiTalentService(context.Background(), repo)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := srv.DeleteChat(tc.chatID)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expErr)
		})
	}
}
