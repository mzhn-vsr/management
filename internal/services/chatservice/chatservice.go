package chatservice

import (
	"context"
	"log/slog"
	"mzhn/management/internal/entity"
	"mzhn/management/internal/lib/logger/sl"
)

type MessageSaver interface {
	Save(ctx context.Context, q, a string) (id string, err error)
}

type ChatRepository interface {
	Invoke(ctx context.Context, input string) (*entity.ChatInvokeOutput, error)
}

type ClassifierRepository interface {
	Classify(ctx context.Context, input string) (*entity.ClassifierResponse, error)
}

type ChatService struct {
	chat       ChatRepository
	classifier ClassifierRepository
	saver      MessageSaver
	logger     *slog.Logger
}

func New(chat ChatRepository, classifier ClassifierRepository, saver MessageSaver) *ChatService {
	return &ChatService{
		chat:       chat,
		classifier: classifier,
		saver:      saver,
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

	log.Debug("saving message")
	id, err := svc.saver.Save(ctx, input, answer.Answer)
	if err != nil {
		log.Error("cannot save message", sl.Err(err))
		return nil, err
	}

	return &entity.ChatInvokeAnswer{
		Id:     id,
		Answer: answer.Answer,
		Class1: classify.Output.Class1,
		Class2: classify.Output.Class2,
	}, nil
}
