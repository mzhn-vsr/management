package feedbackservice

import (
	"context"
	"errors"
	"log/slog"
	"mzhn/management/internal/dto"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/storage"
)

type FeedbackSender interface {
	Send(ctx context.Context, id string, isUseful bool) error
}

type FeedbackStats interface {
	Stats(ctx context.Context) (*dto.FeedbackStats, error)
}

type FeedbackService struct {
	logger *slog.Logger
	sender FeedbackSender
	stats  FeedbackStats
}

func New(sender FeedbackSender, statsrepo FeedbackStats) *FeedbackService {
	return &FeedbackService{
		sender: sender,
		stats:  statsrepo,
		logger: slog.Default().With(slog.String("struct", "FeedbackService")),
	}
}

func (svc *FeedbackService) Send(ctx context.Context, id string, isUseful bool) error {
	log := svc.logger.With(slog.String("method", "Send"))

	log.Debug("sending feedback", slog.String("id", id), slog.Bool("isUseful", isUseful))
	if err := svc.sender.Send(ctx, id, isUseful); err != nil {

		if errors.Is(err, storage.ErrNoMessage) {
			return ErrNoMessage
		}

		if errors.Is(err, storage.ErrFeedbackAlreadySent) {
			return ErrFeedbackAlreadySent
		}

		log.Error("cannot send feedback", sl.Err(err))
		return err
	}

	return nil
}

func (svc *FeedbackService) Stats(ctx context.Context) (*dto.FeedbackStats, error) {
	return svc.stats.Stats(ctx)
}
