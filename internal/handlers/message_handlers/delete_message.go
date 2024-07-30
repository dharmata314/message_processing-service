package messagehandlers

import (
	"log/slog"
	"message_processing-service/api/response"
	errMsg "message_processing-service/internal/err"
	"message_processing-service/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func DeleteMessageByID(log *slog.Logger, messageRepository models.MessageRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.delete.message"
		log := log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		log.Info("Extracted ID from URL", slog.String("id", idStr))
		if idStr == "" {
			log.Error("ID parameter is missing in the URL")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Missing message ID"))
			return
		}
		if err != nil {
			log.Error("Invalid message ID", errMsg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid messsage ID"))
			return
		}

		err = messageRepository.DeleteMessageByID(r.Context(), id)
		if err != nil {
			log.Error("failed to delete message", errMsg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to delete message"))
			return
		}
		log.Info("message deleted")
		render.Status(r, http.StatusNoContent)
		render.JSON(w, r, response.OK())

	}

}
