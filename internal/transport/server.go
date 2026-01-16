package transport

import (
	"TestHitalent/internal/config"
	"TestHitalent/internal/models"
	"TestHitalent/pkg/logger"
	"context"
	"encoding/json"
	"net/http"
)

type HiTalentServiceInterface interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	CreateMessage(chatId string, message *models.Message) (*models.Message, error)
	GetChat(chatId string) (*models.Chat, error)
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
	mux.Handle("/api/v1/chats", CreateChatHandler(s))
	mux.HandleFunc("/api/v1/chats/{id}/messages", CreateMessageHandler(s))
	mux.HandleFunc("/api/v1/chats/{id}", GetChatHandler(s))
	mux.HandleFunc("/api/v1/chats/{id}", DeleteChatHandler(s))
	logger.GetLoggerFromCtx(s.ctx).Info("HTTP server is running")
	addr := s.cfg.Host + ":" + s.cfg.Port
	return http.ListenAndServe(addr, mux)
}

func CreateChatHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Internal server error 1", http.StatusInternalServerError)
			}
		}()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		req := new(models.Chat)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		chat, err := s.service.CreateChat(req)
		if err != nil {
			http.Error(w, "Internal server error 2", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(models.Chat{ID: chat.ID, Title: chat.Title, CreatedAt: chat.CreatedAt})
		if err != nil {
			http.Error(w, "Internal server error 3", http.StatusInternalServerError)
			return
		}
	}
}

func CreateMessageHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Internal server error 1", http.StatusInternalServerError)
				return
			}
		}()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")

		defer r.Body.Close()

		req := new(models.Message)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		msg, err := s.service.CreateMessage(id, req)
		if err != nil {
			http.Error(w, "Internal server error 2", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(models.Message{ID: msg.ID, ChatID: msg.ChatID, Text: msg.Text, CreatedAt: msg.CreatedAt})
		if err != nil {
			http.Error(w, "Internal server error 3", http.StatusInternalServerError)
			return
		}
	}
}

func GetChatHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Internal server error 1", http.StatusInternalServerError)
				return
			}
		}()
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := r.PathValue("id")

		defer r.Body.Close()

		req := new(models.Chat)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		chat, err := s.service.GetChat(id)
		if err != nil {
			http.Error(w, "Internal server error 2", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(chat)
		if err != nil {
			http.Error(w, "Internal server error 3", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteChatHandler(s *HiTalentServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Internal server error 1", http.StatusInternalServerError)
				return
			}
		}()
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := r.PathValue("id")
		defer r.Body.Close()
		err := s.service.DeleteChat(id)
		if err != nil {
			http.Error(w, "Internal server error 2", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(`{"status":"chat deleted"}`)
		if err != nil {
			http.Error(w, "Internal server error 3", http.StatusInternalServerError)
			return
		}
	}
}
