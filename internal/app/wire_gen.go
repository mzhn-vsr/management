// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"mzhn/management/internal/config"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/services/chatservice"
	"mzhn/management/internal/services/faqservice"
	"mzhn/management/internal/services/feedbackservice"
	"mzhn/management/internal/storage/chatapi"
	"mzhn/management/internal/storage/classifierapi"
	"mzhn/management/internal/storage/pg"
)

import (
	_ "github.com/jackc/pgx/stdlib"
	_ "mzhn/management/docs"
)

// Injectors from wire.go:

func New() (*App, func(), error) {
	configConfig := config.New()
	db, cleanup, err := initPG(configConfig)
	if err != nil {
		return nil, nil, err
	}
	faqStore := pg.NewFaqStore(db)
	faqService := faqservice.New(faqStore)
	chatApi, err := connectToChatService(configConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	chatService := chatapi.New(chatApi)
	classifierApi, err := connectToClassifyService(configConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	classifierapiClassifierApi := classifierapi.New(classifierApi)
	feedbackStore := pg.NewFeedbackStore(db)
	chatserviceChatService := chatservice.New(chatService, classifierapiClassifierApi, feedbackStore)
	feedbackService := feedbackservice.New(feedbackStore, feedbackStore)
	app := newApp(configConfig, faqService, chatserviceChatService, feedbackService)
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

func initPG(cfg *config.Config) (*sqlx.DB, func(), error) {
	host := cfg.Pg.Host
	port := cfg.Pg.Port
	user := cfg.Pg.User
	pass := cfg.Pg.Pass
	name := cfg.Pg.Name

	cs := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=disable`, user, pass, host, port, name)
	slog.Info("connecting to database", slog.String("conn", cs))

	db, err := sqlx.Connect("pgx", cs)
	if err != nil {
		return nil, nil, err
	}
	slog.Info("send ping to database")

	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to database", sl.Err(err), slog.String("conn", cs))
		return nil, func() { db.Close() }, err
	}
	slog.Info("connected to database", slog.String("conn", cs))

	return db, func() { db.Close() }, nil
}

func connectToChatService(cfg *config.Config) (*config.ChatApi, error) {
	return &cfg.ChatService, nil
}

func connectToClassifyService(cfg *config.Config) (*config.ClassifierApi, error) {
	return &cfg.ClassifierApi, nil
}
