package transport

import (
	"TestHitalent/internal/config"
	"TestHitalent/internal/models"
	"TestHitalent/internal/service/mocks"
	"TestHitalent/pkg/suberrors"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateChatHandler_Success(t *testing.T) {
	ctx := context.Background()
	ctl := gomock.NewController(t)
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}
	defer ctl.Finish()

	srv := mocks.NewMockHiTalentServiceInterface(ctl)

	inputChat := &models.Chat{
		Title: "Test Chat",
	}
	expectedChat := &models.Chat{
		ID:        1,
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}

	srv.EXPECT().CreateChat(gomock.Any()).Return(expectedChat, nil).Times(1)

	server := NewHiTalentServer(cfg, srv, ctx)

	body, _ := json.Marshal(inputChat)
	req := httptest.NewRequest("POST", "/api/v1/chats", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	CreateChatHandler(server)(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response models.Chat
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, expectedChat.ID, response.ID)
	require.Equal(t, expectedChat.Title, response.Title)
}

func TestCreateChatHandler_Fail(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}

	cases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "empty title",
			requestBody:    `{"title": }`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			srv := mocks.NewMockHiTalentServiceInterface(ctl)
			server := NewHiTalentServer(cfg, srv, ctx)

			req := httptest.NewRequest("POST", "/api/v1/chats", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			CreateChatHandler(server)(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			require.Contains(t, w.Body.String(), tc.expectedError)
		})
	}
}

func TestGetChatHandler_Success(t *testing.T) {
	ctx := context.Background()
	ctl := gomock.NewController(t)
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}
	defer ctl.Finish()

	srv := mocks.NewMockHiTalentServiceInterface(ctl)

	expectedResponse := &models.ChatAndMessagesResponse{
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

	srv.EXPECT().GetChat("1", 20).Return(expectedResponse, nil).Times(1)

	server := NewHiTalentServer(cfg, srv, ctx)

	req := httptest.NewRequest("GET", "/api/v1/chats/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	GetChatHandler(server)(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response models.ChatAndMessagesResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, expectedResponse.Chat.ID, response.Chat.ID)
	require.Equal(t, expectedResponse.Chat.Title, response.Chat.Title)
	require.Len(t, response.Messages, 2)
}

func TestGetChatHandler_Fail(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}

	cases := []struct {
		name           string
		chatID         string
		limit          string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid limit parameter",
			chatID:         "1",
			limit:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid limit parameter",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			srv := mocks.NewMockHiTalentServiceInterface(ctl)
			server := NewHiTalentServer(cfg, srv, ctx)

			url := "/api/v1/chats/" + tc.chatID
			if tc.limit != "" {
				url += "?limit=" + tc.limit
			}
			req := httptest.NewRequest("GET", url, nil)
			req.SetPathValue("id", tc.chatID)

			w := httptest.NewRecorder()

			GetChatHandler(server)(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			require.Contains(t, w.Body.String(), tc.expectedError)
		})
	}
}

func TestCreateMessageHandler_Success(t *testing.T) {
	ctx := context.Background()
	ctl := gomock.NewController(t)
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}
	defer ctl.Finish()

	srv := mocks.NewMockHiTalentServiceInterface(ctl)

	inputMessage := &models.Message{
		Text: "Test message",
	}
	expectedMessage := &models.Message{
		ID:        1,
		ChatID:    1,
		Text:      "Test message",
		CreatedAt: time.Now(),
	}

	srv.EXPECT().CreateMessage("1", gomock.Any()).Return(expectedMessage, nil).Times(1)

	server := NewHiTalentServer(cfg, srv, ctx)

	body, _ := json.Marshal(inputMessage)
	req := httptest.NewRequest("POST", "/api/v1/chats/1/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	CreateMessageHandler(server)(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response models.Message
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, expectedMessage.ID, response.ID)
	require.Equal(t, expectedMessage.Text, response.Text)
}

func TestCreateMessageHandler_Fail(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}

	cases := []struct {
		name           string
		chatID         string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid JSON",
			chatID:         "1",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "empty body",
			chatID:         "1",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "malformed JSON",
			chatID:         "1",
			requestBody:    `{"text": }`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			srv := mocks.NewMockHiTalentServiceInterface(ctl)
			server := NewHiTalentServer(cfg, srv, ctx)

			req := httptest.NewRequest("POST", "/api/v1/chats/"+tc.chatID+"/messages", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tc.chatID)

			w := httptest.NewRecorder()

			CreateMessageHandler(server)(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			require.Contains(t, w.Body.String(), tc.expectedError)
		})
	}
}

func TestDeleteChatHandler_Success(t *testing.T) {
	ctx := context.Background()
	ctl := gomock.NewController(t)
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}
	defer ctl.Finish()

	srv := mocks.NewMockHiTalentServiceInterface(ctl)

	srv.EXPECT().DeleteChat("1").Return(nil).Times(1)

	server := NewHiTalentServer(cfg, srv, ctx)

	req := httptest.NewRequest("DELETE", "/api/v1/chats/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	DeleteChatHandler(server)(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Empty(t, w.Body.String())
}

func TestDeleteChatHandler_Fail(t *testing.T) {
	ctx := context.Background()
	ctl := gomock.NewController(t)
	cfg := &config.Config{
		Host: "localhost",
		Port: "4047",
	}
	defer ctl.Finish()

	srv := mocks.NewMockHiTalentServiceInterface(ctl)

	srv.EXPECT().DeleteChat("999").Return(suberrors.ErrChatNotFound).Times(1)

	server := NewHiTalentServer(cfg, srv, ctx)

	req := httptest.NewRequest("DELETE", "/api/v1/chats/999", nil)
	req.SetPathValue("id", "999")

	w := httptest.NewRecorder()

	DeleteChatHandler(server)(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Chat not found")
}
