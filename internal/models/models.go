package models

import (
	"context"
	"message_processing-service/internal/entities"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *entities.Message) error
	FindMessageByID(ctx context.Context, id int) (entities.Message, error)
	DeleteMessageByID(ctx context.Context, id int) error
	MarkMessageProcessed(ctx context.Context, message_id string) error
	GetStatistics(ctx context.Context) (map[string]int, error)
}
