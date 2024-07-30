package messagesrepo

import (
	"context"
	"fmt"
	"log/slog"
	"message_processing-service/internal/entities"
	errMsg "message_processing-service/internal/err"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewMessageRepository(db *pgxpool.Pool, log *slog.Logger) *MessageRepository {
	return &MessageRepository{db: db, log: log}
}

func (m *MessageRepository) CreateMessage(ctx context.Context, message *entities.Message) error {
	err := m.db.QueryRow(ctx, `INSERT INTO messages (content, status) VALUES ($1, $2) RETURNING id`, message.Content, "pending").Scan(&message.ID)
	if err != nil {
		m.log.Error("failed to create user", errMsg.Err(err))
		return err
	}
	return nil
}

func (m *MessageRepository) FindMessageByID(ctx context.Context, id int) (entities.Message, error) {
	query, err := m.db.Query(ctx, `SELECT content, status FROM messages WHERE id = $1`, id)
	if err != nil {
		m.log.Error("error querying messages", errMsg.Err(err))
		return entities.Message{}, err
	}
	defer query.Close()
	row := entities.Message{}
	if !query.Next() {
		m.log.Error("message not found")
		return entities.Message{}, fmt.Errorf("user not found")
	} else {
		err := query.Scan(&row.ID, &row.Content, &row.Status)
		if err != nil {
			m.log.Error("error scanning messages", errMsg.Err(err))
			return entities.Message{}, err
		}
	}
	return row, nil
}

func (m *MessageRepository) DeleteMessageByID(ctx context.Context, id int) error {
	_, err := m.db.Exec(ctx, `DELETE FROM messages WHERE id = $1`, id)
	if err != nil {
		m.log.Error("failed to delete user", errMsg.Err(err))
		return err
	}
	return nil
}

func (m *MessageRepository) MarkMessageProcessed(ctx context.Context, message_id string) error {

	id, err := strconv.Atoi(message_id)
	if err != nil {
		m.log.Error("error converting message id to int", errMsg.Err(err))
		return err
	}

	_, err = m.db.Exec(ctx, `UPDATE messages SET status = 'processed', processed_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	if err != nil {
		m.log.Error("failed to update message status", errMsg.Err(err))
		return err
	}
	return nil
}

func (m *MessageRepository) GetStatistics(ctx context.Context) (map[string]int, error) {

	var total, processed, pending int
	statMap := make(map[string]int)

	err := m.db.QueryRow(ctx, `SELECT COUNT(*) FROM messages`).Scan(&total)
	if err != nil {
		m.log.Error("failed getting statistics", errMsg.Err(err))
		return nil, err
	}
	statMap["total_messages"] = total

	err = m.db.QueryRow(ctx, `SELECT COUNT(*) FROM messages WHERE status = 'processed'`).Scan(&processed)
	if err != nil {
		m.log.Error("failed getting statistics", errMsg.Err(err))
		return nil, err
	}
	statMap["processed_messages"] = processed

	err = m.db.QueryRow(ctx, `SELECT COUNT(*) FROM messages WHERE status = 'pending'`).Scan(&pending)
	if err != nil {
		m.log.Error("failed getting statistics", errMsg.Err(err))
		return nil, err
	}
	statMap["pending_messages"] = pending
	return statMap, nil

}
