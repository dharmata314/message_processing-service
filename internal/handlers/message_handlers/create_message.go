package messagehandlers

import (
	"log/slog"
	"message_processing-service/api/response"
	"message_processing-service/internal/entities"
	errMsg "message_processing-service/internal/err"
	"message_processing-service/internal/kafka"
	"message_processing-service/internal/models"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type RequestMessage struct {
	Content string `json:"content" validate:"required"`
}

type ResponseMessage struct {
	response.Response
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func NewMessage(log *slog.Logger, MessageRepository models.MessageRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.createMessage.New"

		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestMessage
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", errMsg.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body request", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid requets", errMsg.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		message := entities.Message{Content: req.Content}
		err = MessageRepository.CreateMessage((r.Context()), &message)
		if err != nil {
			log.Error("failed to create message", errMsg.Err((err)))
			render.JSON(w, r, response.Error("failed to create message"))
			return
		}
		log.Info("message added to postgres")
		err = kafka.ProduceMessage("localhost:9092", "test", message, log)
		if err != nil {
			log.Error("failed to send message to kafka", errMsg.Err(err))
			render.JSON(w, r, response.Error("failed to send message to kafka"))
			return
		}
		responseOK(w, r, "pending", message.ID)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, status string, messageID int) {
	render.JSON(w, r, ResponseMessage{
		response.OK(),
		messageID,
		status,
	})
}
