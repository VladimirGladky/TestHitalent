package transport

import (
	"TestHitalent/internal/config"
	"TestHitalent/internal/models"
	"TestHitalent/pkg/logger"
	"TestHitalent/pkg/suberrors"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

//go:generate mockgen -source=server.go -destination=../service/mocks/mock_service.go -package=mocks HiTalentServiceInterface

type HiTalentServiceInterface interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	CreateMessage(chatId string, message *models.Message) (*models.Message, error)
	GetChat(chatId string, limit int) (*models.ChatAndMessagesResponse, error)
	DeleteChat(chatId string) error
}

type HiTalentServer struct {
	cfg     *config.Config
	service HiTalentServiceInterface
	ctx     context.Context
}

func NewHiTalentServer(cfg *config.Config, service HiTalentServiceInterface, ctx context.Context) *HiTalentServer {
	return &HiTalentServer{
		cfg:     cfg,
		service: service,
		ctx:     ctx,
	}
}

func (s *HiTalentServer) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/chats", CreateChatHandler(s))
	mux.HandleFunc("POST /api/v1/chats/{id}/messages", CreateMessageHandler(s))
	mux.HandleFunc("GET /api/v1/chats/{id}", GetChatHandler(s))
	mux.HandleFunc("DELETE /api/v1/chats/{id}", DeleteChatHandler(s))
	logger.GetLoggerFromCtx(s.ctx).Info("HTTP server is running")
	addr := s.cfg.Host + ":" + s.cfg.Port
	return http.ListenAndServe(addr, mux)
}

func CreateChatHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error 1", "description": "` + fmt.Sprint(rec) + `"}`))
			}
		}()

		defer r.Body.Close()
		req := new(models.Chat)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "Invalid request body", "description": "` + err.Error() + `"}`))
			return
		}
		chat, err := s.service.CreateChat(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error 2", "description": "` + err.Error() + `"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(models.Chat{ID: chat.ID, Title: chat.Title, CreatedAt: chat.CreatedAt})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error 3", "description": "` + err.Error() + `"}`))
			return
		}
	}
}

func CreateMessageHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error 1", "description": "` + fmt.Sprint(rec) + `"}`))
				return
			}
		}()

		id := r.PathValue("id")

		defer r.Body.Close()

		req := new(models.Message)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "Invalid request body", "description": "` + err.Error() + `"}`))
			return
		}

		msg, err := s.service.CreateMessage(id, req)
		if err != nil {
			if errors.Is(err, suberrors.ErrChatNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error": "Chat not found"}`))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error 2", "description": "` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(models.Message{ID: msg.ID, ChatID: msg.ChatID, Text: msg.Text, CreatedAt: msg.CreatedAt})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error 3", "description": "` + err.Error() + `"}`))
			return
		}
	}
}

func GetChatHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error 1", "description": "` + fmt.Sprint(rec) + `"}`))
				return
			}
		}()
		id := r.PathValue("id")

		limitStr := r.URL.Query().Get("limit")
		limit := 20

		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error": "Invalid limit parameter", "description": "` + err.Error() + `"}`))
				return
			}
			limit = parsedLimit
			if limit > 100 {
				limit = 100
			}
			if limit < 1 {
				limit = 1
			}
		}

		defer r.Body.Close()
		chatAndMessage, err := s.service.GetChat(id, limit)
		if err != nil {
			if errors.Is(err, suberrors.ErrChatNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error": "Chat not found"}`))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error 2", "description": "` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(chatAndMessage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error 3", "description": "` + err.Error() + `"}`))
			return
		}
	}
}

func DeleteChatHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error 1", "description": "` + fmt.Sprint(rec) + `"}`))
				return
			}
		}()
		id := r.PathValue("id")
		defer r.Body.Close()
		err := s.service.DeleteChat(id)
		if err != nil {
			if errors.Is(err, suberrors.ErrChatNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error": "Chat not found"}`))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error 2", "description": "` + err.Error() + `"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
