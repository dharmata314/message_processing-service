package main

import (
	"context"
	"fmt"
	"log/slog"
	"message_processing-service/internal/config"
	"message_processing-service/internal/database"
	messagesrepo "message_processing-service/internal/database/messages_repo"
	usersrepo "message_processing-service/internal/database/users_repo"
	errMsg "message_processing-service/internal/err"
	messagehandlers "message_processing-service/internal/handlers/message_handlers"
	userhandlers "message_processing-service/internal/handlers/user_handlers"
	"message_processing-service/internal/jwt"
	"message_processing-service/internal/kafka"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()
	log.Debug("debug messages are active")

	pg, err := connectToPostgres(cfg, log)
	if err != nil {
		log.Error("failed to create postgres db", errMsg.Err(err))
		os.Exit(1)
	}

	log.Info("connecting to postgres")

	defer pg.Close()
	if pg == nil {
		log.Error("failed to connect to postgres")
		os.Exit(1)
	}

	if err := pg.Ping(context.Background()); err != nil {
		log.Error("failed to ping postgres db", errMsg.Err(err))
		os.Exit(1)
	}

	log.Info("postgres db connected successfully")

	log.Info("application started")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	messageRepository := messagesrepo.NewMessageRepository(pg.Db, log)
	userRepository := usersrepo.NewUserRepository(pg.Db, log)
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, log)

	router.Post("/users/new", userhandlers.NewUser(log, userRepository))
	router.Post("/login", userhandlers.LoginFunc(log, userRepository, jwtManager))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Delete("/users/{id}", userhandlers.DeleteUserHandler(log, userRepository)) // не парсит id запроса

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Patch("/users/{id}", userhandlers.NewUpdateUserHandler(userRepository, log))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Post("/message", messagehandlers.NewMessage(log, messageRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Get("/message/statistics", messagehandlers.GetStatistics(log, messageRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Delete("/message/{id}", messagehandlers.DeleteMessageByID(log, messageRepository))

	log.Info("starting server", slog.String("addr", cfg.HTTPServer.Addr))

	server := &http.Server{
		Addr:              cfg.HTTPServer.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server", errMsg.Err(err))
		}
	}()

	go kafka.ConsumeMessage("kafka:9092", "test", log, messageRepository)

	select {}

}

func setupLogger() *slog.Logger {
	var log *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}

func connectToPostgres(cfg *config.Config, log *slog.Logger) (*database.Postgres, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	pg, err := database.NewPG(context.Background(), connString, log, cfg)
	return pg, err
}
