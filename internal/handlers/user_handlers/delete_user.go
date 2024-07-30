package userhandlers

import (
	"log/slog"
	"message_processing-service/api/response"
	errMsg "message_processing-service/internal/err"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func DeleteUserHandler(log *slog.Logger, userRepo User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.delete.user"
		log := log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("ID parameter is missing in the URL")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Missing user ID"))
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid user ID"))
			return
		}

		err = userRepo.DeleteUserById(r.Context(), id)
		if err != nil {
			log.Error("Failed to delete user", errMsg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to delete user"))
			return
		}
		log.Info("user deleted")
		render.Status(r, http.StatusNoContent)
		render.JSON(w, r, response.OK())
	}
}
