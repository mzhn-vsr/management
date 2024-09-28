package pg

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mzhn/management/internal/dto"
	"mzhn/management/internal/entity"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/storage"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type FeedbackStore struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewFeedbackStore(db *sqlx.DB) *FeedbackStore {
	return &FeedbackStore{
		db:     db,
		logger: slog.Default().With(slog.String("struct", "pg.FeedbackStore")),
	}
}

func (store *FeedbackStore) Save(ctx context.Context, question string, answer string) (id string, err error) {
	log := store.logger.With(slog.String("method", "Save"))

	query, args, err := squirrel.
		Insert(feedbackTable).
		Columns("question", "answer").
		Values(question, answer).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Warn("cannot build query", sl.Err(err))
		return "", nil
	}

	log = log.With(slog.String("query", query), slog.Any("args", args))

	log.Debug("executing")
	if err := store.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		log.Warn("cannot execute", sl.Err(err))
		return "", err
	}

	return id, nil
}

func (store *FeedbackStore) Send(ctx context.Context, id string, isUseful bool) error {
	log := store.logger.With(slog.String("method", "Send"))

	feedback, err := store.Find(ctx, id)
	if err != nil {
		return err
	}

	if feedback.IsUseful != nil {
		return storage.ErrFeedbackAlreadySent
	}

	query, args, err := squirrel.
		Update(feedbackTable).
		Set("is_useful", isUseful).
		Where(squirrel.Eq{"id": feedback.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Warn("cannot build query", sl.Err(err))
		return nil
	}

	log = log.With(slog.String("query", query), slog.Any("args", args))

	log.Debug("executing")
	if _, err := store.db.ExecContext(ctx, query, args...); err != nil {
		if e, ok := err.(pgx.PgError); ok {
			log.Warn("pg error", sl.PgError(e))
			switch e.Code {
			case "22P02":
				return storage.ErrNoMessage
			default:
				return err
			}
		}

		log.Warn("cannot execute", sl.Err(err))
		return err
	}

	return nil
}

func (store *FeedbackStore) Find(ctx context.Context, id string) (*entity.Feedback, error) {

	log := store.logger.With(slog.String("method", "Find"))

	query, args, err := squirrel.
		Select("*").
		From(feedbackTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Warn("cannot build query", sl.Err(err))
		return nil, err
	}

	log = log.With(slog.String("query", query), slog.Any("args", args))

	feedback := new(entity.Feedback)
	log.Debug("executing")
	if err := store.db.GetContext(ctx, feedback, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoMessage
		}

		log.Warn("cannot execute", sl.Err(err))
		return nil, err
	}

	return feedback, nil
}

func (store *FeedbackStore) Stats(ctx context.Context) (*dto.FeedbackStats, error) {

	log := store.logger.With(slog.String("method", "Find"))

	query, args, err := squirrel.
		Select("COUNT(case is_useful when true then true else null end), COUNT(case is_useful when false then false else null end), COUNT(id)").
		From(feedbackTable).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Warn("cannot build query", sl.Err(err))
		return nil, err
	}

	log = log.With(slog.String("query", query), slog.Any("args", args))

	stats := new(dto.FeedbackStats)
	log.Debug("executing")
	if err := store.db.QueryRowContext(ctx, query, args...).Scan(&stats.Positive, &stats.Negative, &stats.Total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoMessage
		}

		log.Warn("cannot execute", sl.Err(err))
		return nil, err
	}

	return stats, nil
}
