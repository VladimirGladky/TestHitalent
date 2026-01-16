package repository

import (
	"TestHitalent/internal/models"
	"TestHitalent/pkg/suberrors"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type HiTalentRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewHiTalentRepository(db *gorm.DB, ctx context.Context) *HiTalentRepository {
	return &HiTalentRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r *HiTalentRepository) CreateChat(chat *models.Chat) (*models.Chat, error) {
	if err := r.db.WithContext(r.ctx).Create(chat).Error; err != nil {
		return nil, err
	}
	return chat, nil
}

func (r *HiTalentRepository) GetChat(chatId int, limit int) (*models.ChatAndMessagesResponse, error) {
	var chat models.Chat

	if err := r.db.
		WithContext(r.ctx).
		First(&chat, chatId).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, suberrors.ErrChatNotFound
		}
		return nil, err
	}

	var messages []*models.Message

	if err := r.db.
		WithContext(r.ctx).
		Where("chat_id = ?", chatId).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error; err != nil {

		return nil, err
	}

	return &models.ChatAndMessagesResponse{
		Chat:     &chat,
		Messages: messages,
	}, nil
}

func (r *HiTalentRepository) DeleteChat(chatId int) error {
	result := r.db.
		WithContext(r.ctx).
		Delete(&models.Chat{}, chatId)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return suberrors.ErrChatNotFound
	}

	return nil
}

func (r *HiTalentRepository) CreateMessage(chatId int, message *models.Message) (*models.Message, error) {
	message.ChatID = chatId

	err := r.db.
		WithContext(r.ctx).
		Create(message).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, suberrors.ErrChatNotFound
		}
		return nil, err
	}

	return message, nil
}
