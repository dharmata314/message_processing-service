package messagehandlers

import (
	"log/slog"
	"message_processing-service/api/response"
	errMsg "message_processing-service/internal/err"
	"message_processing-service/internal/models"
	"net/http"

	"github.com/go-chi/render"
)

type ResponseStatistics struct {
	response.Response
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Processed int `json:"processed"`
}

func GetStatistics(log *slog.Logger, messageRepository models.MessageRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.get.statistics"

		log := log.With(
			slog.String("options", loggerOptions))

		statMap, err := messageRepository.GetStatistics(r.Context())
		if err != nil {
			log.Error("failed to acquire statistics", errMsg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to acquire statistics"))
			return
		}

		log.Info("statistics received")
		responseOkStatistics(w, r, statMap)

	}
}

func responseOkStatistics(w http.ResponseWriter, r *http.Request, statMap map[string]int) {
	render.JSON(w, r, ResponseStatistics{
		response.OK(),
		statMap["total_messages"],
		statMap["pending_messages"],
		statMap["processed_messages"],
	})
}
