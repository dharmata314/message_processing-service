package kafka

import (
	"context"
	"log/slog"
	"message_processing-service/internal/entities"
	errMsg "message_processing-service/internal/err"
	"message_processing-service/internal/models"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func ProduceMessage(broker, topic string, message entities.Message, log *slog.Logger) error {

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})

	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.Itoa(message.ID)),
		Value: []byte(message.Content),
	})
	if err != nil {
		log.Error("error writing message to kafka")
		errMsg.Err(err)
		return err
	}

	return nil
}

func ConsumeMessage(broker, topic string, log *slog.Logger, messageRepo models.MessageRepository) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Error("error reading message from kafka")
			errMsg.Err(err)
			return
		}

		if len(msg.Key) == 0 {
			log.Error("received message with empty key")
			continue
		}

		err = messageRepo.MarkMessageProcessed(context.Background(), string(msg.Key))
		if err != nil {
			log.Error("error marking message as processed")
			errMsg.Err(err)
			return
		}

		if msg.Value != nil {
			log.Info("Message processed", "value", string(msg.Value))
		} else {
			log.Info("Message processed", "value", "nil")
		}
	}
}
