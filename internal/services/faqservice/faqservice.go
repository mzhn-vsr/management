package faqservice

import (
	"context"
	"errors"
	"mzhn/management/internal/dto"
	"mzhn/management/internal/entity"
	"mzhn/management/internal/storage"
)

type FaqStore interface {
	Create(ctx context.Context, entry *dto.FaqEntryCreate) (string, error)
	Find(ctx context.Context, id string) (*entity.FaqEntry, error)
	List(ctx context.Context, filters dto.FaqEntryList) ([]*entity.FaqEntry, uint64, error)
	Update(ctx context.Context, entry *dto.FaqEntryUpdate) error
	Delete(ctx context.Context, id string) error
}

type FaqService struct {
	FaqStore
}

func New(faqstore FaqStore) *FaqService {
	return &FaqService{
		FaqStore: faqstore,
	}
}

func (s *FaqService) Find(ctx context.Context, id string) (*entity.FaqEntry, error) {
	e, err := s.FaqStore.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNoFaqEntry) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e, nil
}
