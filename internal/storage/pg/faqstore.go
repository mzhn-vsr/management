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
	"github.com/jmoiron/sqlx"
)

type FaqStore struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewFaqStore(db *sqlx.DB) *FaqStore {
	return &FaqStore{
		db:     db,
		logger: slog.Default().With(slog.String("struct", "pg.FaqStore")),
	}
}

func (s *FaqStore) Create(ctx context.Context, e *dto.FaqEntryCreate) (id string, err error) {
	log := s.logger.With(slog.String("method", "Create"))

	values := make([]any, 0, 4)

	builder := squirrel.Insert(faqTable).
		Columns("question", "answer").
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	values = append(values, e.Question)
	values = append(values, e.Answer)

	if e.Classifier1 != nil {
		log.Debug("appening classifier1", slog.String("classifier1", *e.Classifier1))
		builder = builder.Columns("classifier1")
		values = append(values, e.Classifier1)
	}

	if e.Classifier2 != nil {
		log.Debug("appening classifier2", slog.String("classifier2", *e.Classifier2))
		builder = builder.Columns("classifier2")
		values = append(values, e.Classifier2)
	}

	query, args, err := builder.
		Values(values...).
		ToSql()
	if err != nil {
		log.Warn("error with building query", sl.Err(err))
		return "", err
	}

	qlog := log.With(slog.String("query", query), slog.Any("args", args))

	qlog.Debug("executing query")

	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		qlog.Error("error scaning pg row", sl.Err(err))
		return "", err
	}

	qlog.Debug("faq created", slog.String("id", id))

	return id, nil
}

func (s *FaqStore) Find(ctx context.Context, id string) (*entity.FaqEntry, error) {
	log := s.logger.With(slog.String("method", "Find"))

	query, args, err := squirrel.Select("*").
		From(faqTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Warn("error with building query", sl.Err(err))
		return nil, err
	}

	entry := new(entity.FaqEntry)
	if err := s.db.Get(entry, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("faq entry not found", slog.String("id", id))
			return nil, storage.ErrNoFaqEntry
		}
		log.Warn("cannot query ", sl.Err(err))
		return nil, err
	}

	return entry, nil
}

func (s *FaqStore) List(ctx context.Context, filters dto.FaqEntryList) ([]*entity.FaqEntry, uint64, error) {
	log := s.logger.With(slog.String("method", "List"))

	builder := squirrel.Select().
		From(faqTable).
		PlaceholderFormat(squirrel.Dollar)

	if filters.Limit != nil {
		builder.Limit(uint64(*filters.Limit))
	}

	if filters.Offset != nil {
		builder.Offset(uint64(*filters.Offset))
	}

	query, args, err := builder.
		Column("*").
		ToSql()
	if err != nil {
		log.Warn("error with building query", sl.Err(err))
		return nil, 0, err
	}

	qlog := log.With(slog.String("query", query), slog.Any("args", args))

	entries := make([]*entity.FaqEntry, 0, func() uint {
		if filters.Limit != nil {
			return uint(*filters.Limit)
		}
		return 10
	}())
	if err := s.db.SelectContext(ctx, &entries, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			qlog.Debug("faq entries not found")
			return nil, 0, storage.ErrNoFaqEntry
		}
		qlog.Warn("cannot query", sl.Err(err))
		return nil, 0, err
	}

	query, args, err = builder.
		Column("COUNT(*)").
		ToSql()
	if err != nil {
		log.Warn("error with building query", sl.Err(err))
		return nil, 0, err
	}
	var total uint64
	if err := s.db.QueryRow(query, args...).Scan(&total); err != nil {
		log.Warn("cannot query count sql", sl.Err(err))
		return nil, 0, err
	}

	return entries, total, nil
}

func (s *FaqStore) Update(ctx context.Context, updated *dto.FaqEntryUpdate) error {

	log := s.logger.With(slog.String("method", "Update"))

	builder := squirrel.
		Update(faqTable).
		Where(squirrel.Eq{"id": updated.Id}).
		PlaceholderFormat(squirrel.Dollar)

	if updated.Question != nil {
		builder.Set("question", updated.Question)
	}
	if updated.Answer != nil {
		builder.Set("answer", updated.Answer)
	}
	if updated.Classifier1 != nil {
		builder.Set("classifier1", updated.Classifier1)
	}
	if updated.Classifier2 != nil {
		builder.Set("classifier2", updated.Classifier2)
	}

	query, args, err := builder.
		ToSql()
	if err != nil {
		log.Warn("error with building query", sl.Err(err))
		return err
	}

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("faq entry not found", slog.String("id", updated.Id))
			return storage.ErrNoFaqEntry
		}
		log.Warn("cannot query ", sl.Err(err))
		return err
	}

	return nil
}

func (s *FaqStore) Delete(ctx context.Context, id string) error {
	log := s.logger.With(slog.String("method", "Delete"))

	builder := squirrel.
		Delete(faqTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.
		ToSql()
	if err != nil {
		log.Warn("error with building query", sl.Err(err))
		return err
	}

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("faq entry not found", slog.String("id", id))
			return storage.ErrNoFaqEntry
		}
		log.Warn("cannot query ", sl.Err(err))
		return err
	}

	return nil
}
