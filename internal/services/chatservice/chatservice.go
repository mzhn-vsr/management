package chatservice

import (
	"context"
	"log/slog"
	"mzhn/management/internal/entity"
	"mzhn/management/internal/lib/logger/sl"
)

type ChatRepository interface {
	Invoke(ctx context.Context, input string) (*entity.ChatInvokeOutput, error)
}

type ClassifierRepository interface {
	Classify(ctx context.Context, input string) (*entity.ClassifierResponse, error)
}

type ChatService struct {
	chat       ChatRepository
	classifier ClassifierRepository
	logger     *slog.Logger
}

func New(chat ChatRepository, classifier ClassifierRepository) *ChatService {
	return &ChatService{
		chat:       chat,
		classifier: classifier,
		logger:     slog.Default().With(slog.String("struct", "ChatService")),
	}
}

func (svc *ChatService) Invoke(ctx context.Context, input string) (*entity.ChatInvokeAnswer, error) {
	log := svc.logger.With(slog.String("method", "Invoke"))

	log.Debug("invoking chat message", slog.String("input", input))
	answer, err := svc.chat.Invoke(ctx, input)
	if err != nil {
		log.Warn("cannot invoke chat message", sl.Err(err))
		return nil, err
	}

	classify, err := svc.classifier.Classify(ctx, input)
	if err != nil {
		log.Warn("cannot classify", sl.Err(err))
		return nil, err
	}

	return &entity.ChatInvokeAnswer{
		Answer: answer.Answer,
		Class1: classify.Class1,
		Class2: classify.Class2,
	}, nil
}
