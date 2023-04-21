package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/wiliamvj/go-open-ai/internal/domain/entity"
	"github.com/wiliamvj/go-open-ai/internal/infra/db"
)

type ChatRepositoryMySQL struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewChatRepositoryMySQL(dbt *sql.DB) *ChatRepositoryMySQL {
	return &ChatRepositoryMySQL{
		DB:      dbt,
		Queries: db.New(dbt),
	}
}

func (r *ChatRepositoryMySQL) CreateChat(ctx context.Context, chat *entity.Chat) error {
	err :=
		r.Queries.CreateChat(
			ctx,
			db.CreateChatParams{
				ID:               chat.ID,
				UserID:           chat.UserID,
				InitialMessageID: chat.InitialSystemMessage.Content,
				Status:           chat.Status,
				TokenUsage:       int32(chat.TokensUsage),
				Model:            chat.Config.Model.Name,
				ModelMaxTokens:   int32(chat.Config.Model.MaxTokens),
				Temperature:      float64(chat.Config.Temperature),
				TopP:             float64(chat.Config.TopP),
				N:                int32(chat.Config.N),
				Stop:             chat.Config.Stop[0],
				MaxTokens:        int32(chat.Config.MaxTokens),
				PresencePenalty:  float64(chat.Config.PresencePenalty),
				FrequencyPenalty: float64(chat.Config.FrequencyPenalty),
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		)
	if err != nil {
		return err
	}

	err = r.Queries.AddMessage(
		ctx,
		db.AddMessageParams{
			ID:        chat.InitialSystemMessage.ID,
			ChatID:    chat.ID,
			Content:   chat.InitialSystemMessage.Content,
			Role:      chat.InitialSystemMessage.Role,
			Tokens:    int32(chat.InitialSystemMessage.Tokens),
			CreatedAt: chat.InitialSystemMessage.CreatedAt,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepositoryMySQL) FindChatByID(ctx context.Context, chatID string) (*entity.Chat, error) {
	chat := &entity.Chat{}
	res, err := r.Queries.FindChatByID(ctx, chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}
	chat.ID = res.ID
	chat.UserID = res.UserID
	chat.Status = res.Status
	chat.TokensUsage = int(res.TokenUsage)
	chat.Config = &entity.ChatConfig{
		Model: &entity.Model{
			Name:      res.Model,
			MaxTokens: int(res.ModelMaxTokens),
		},
		Temperature:      float32(res.Temperature),
		TopP:             float32(res.TopP),
		N:                int(res.N),
		Stop:             []string{res.Stop},
		MaxTokens:        int(res.MaxTokens),
		PresencePenalty:  float32(res.PresencePenalty),
		FrequencyPenalty: float32(res.FrequencyPenalty),
	}

	messages, err := r.Queries.FindMessagesByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		chat.Messages = append(chat.Messages, &entity.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      message.Role,
			Tokens:    int(message.Tokens),
			Model:     &entity.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}

	erasedMessages, err := r.Queries.FindErasedMessagesByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	for _, message := range erasedMessages {
		chat.ErasedMessages = append(chat.ErasedMessages, &entity.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      message.Role,
			Tokens:    int(message.Tokens),
			Model:     &entity.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}
	return chat, nil
}
