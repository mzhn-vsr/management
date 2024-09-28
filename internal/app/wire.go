//go:build wireinject
// +build wireinject

package app

import (
	"fmt"
	"log/slog"

	"mzhn/management/internal/config"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/services/chatservice"
	"mzhn/management/internal/services/faqservice"
	"mzhn/management/internal/services/feedbackservice"
	"mzhn/management/internal/storage/chatapi"
	"mzhn/management/internal/storage/classifierapi"
	"mzhn/management/internal/storage/pg"

	"github.com/google/wire"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		newApp,

		faqservice.New,
		chatservice.New,
		feedbackservice.New,

		pg.NewFeedbackStore,
		pg.NewFaqStore,
		chatapi.New,
		classifierapi.New,

		initPG,
		connectToChatService,
		connectToClassifyService,
		config.New,

		wire.Bind(new(faqservice.FaqStore), new(*pg.FaqStore)),
		wire.Bind(new(chatservice.ChatRepository), new(*chatapi.ChatService)),
		wire.Bind(new(chatservice.ClassifierRepository), new(*classifierapi.ClassifierApi)),

		wire.Bind(new(chatservice.MessageSaver), new(*pg.FeedbackStore)),
		wire.Bind(new(feedbackservice.FeedbackSender), new(*pg.FeedbackStore)),
		wire.Bind(new(feedbackservice.FeedbackStats), new(*pg.FeedbackStore)),
	))
}

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
