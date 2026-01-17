package service

import (
	"TestHitalent/internal/models"
	"TestHitalent/pkg/suberrors"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type HiTalentRepositoryInterface interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	GetChat(chatId int, limit int) (*models.ChatAndMessagesResponse, error)
	CreateMessage(chatId int, message *models.Message) (*models.Message, error)
	DeleteChat(chatId int) error
}

type HiTalentService struct {
	repo     HiTalentRepositoryInterface
	ctx      context.Context
	validate *validator.Validate
}

func NewHiTalentService(ctx context.Context, repo HiTalentRepositoryInterface) *HiTalentService {
	return &HiTalentService{
		repo:     repo,
		ctx:      ctx,
		validate: validator.New(),
	}
}

func (s *HiTalentService) CreateChat(chat *models.Chat) (*models.Chat, error) {
	if chat == nil {
		return nil, errors.New("chat is nil")
	}

	chat.Title = strings.TrimSpace(chat.Title)

	if err := s.validate.Struct(chat); err != nil {
		return nil, err
	}

	return s.repo.CreateChat(chat)
}

func (s *HiTalentService) GetChat(chatId string, limit int) (*models.ChatAndMessagesResponse, error) {
	chatID, err := strconv.Atoi(chatId)
	if err != nil {
		return nil, suberrors.ErrInvalidChatId
	}
	if chatID <= 0 {
		return nil, suberrors.ErrNotPositiveChatId
	}
	return s.repo.GetChat(chatID, limit)
}

func (s *HiTalentService) CreateMessage(chatId string, message *models.Message) (*models.Message, error) {
	chatID, err := strconv.Atoi(chatId)
	if err != nil {
		return nil, suberrors.ErrInvalidChatId
	}
	if chatID <= 0 {
		return nil, suberrors.ErrNotPositiveChatId
	}

	if message == nil {
		return nil, errors.New("message is nil")
	}

	message.Text = strings.TrimSpace(message.Text)

	if err = s.validate.Struct(message); err != nil {
		return nil, err
	}

	return s.repo.CreateMessage(chatID, message)
}

func (s *HiTalentService) DeleteChat(chatId string) error {
	chatID, err := strconv.Atoi(chatId)
	if err != nil {
		return suberrors.ErrInvalidChatId
	}
	if chatID <= 0 {
		return suberrors.ErrNotPositiveChatId
	}
	return s.repo.DeleteChat(chatID)
}
